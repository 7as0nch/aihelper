package base

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strings"

	"github.com/example/aichat/backend/models"
	"github.com/example/aichat/backend/models/generator/model"
	"github.com/example/aichat/backend/pkg/auth"
	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"golang.org/x/crypto/bcrypt"
)

type SysUserRepo interface {
	GetByAccount(ctx context.Context, account string) (*model.SysUser, error)
	GetById(ctx context.Context, id int64) (*model.SysUser, error)
	Create(ctx context.Context, user *model.SysUser) error
	GetUserAuthByUserIDAndType(ctx context.Context, userID int64, authType model.AuthType) (*model.SysUserAuth, error)
	GetUserAuthByTypeAndIdentifier(ctx context.Context, authType model.AuthType, identifier string) (*model.SysUserAuth, error)
	CreateUserAuth(ctx context.Context, userAuth *model.SysUserAuth) error
	UpdateUserAuth(ctx context.Context, userAuth *model.SysUserAuth) error
}

type SysUserUseCase struct {
	user SysUserRepo
	auth auth.AuthRepo
}

func NewSysUserUseCase(user SysUserRepo, auth auth.AuthRepo) *SysUserUseCase {
	return &SysUserUseCase{
		user: user,
		auth: auth,
	}
}

// Login validates user credentials and returns a JWT token.
func (s *SysUserUseCase) Login(ctx context.Context, account, password string) (string, error) {
	account = strings.TrimSpace(account)
	password = strings.TrimSpace(password)
	if account == "" || password == "" {
		return "", kerrors.BadRequest("INVALID_ARGUMENT", "username or password is empty")
	}

	user, err := s.user.GetByAccount(ctx, account)
	if err != nil {
		return "", err
	}
	if user == nil || user.ID == 0 {
		return "", kerrors.Unauthorized(auth.Reason, "invalid credentials")
	}

	matched := false
	needUpgrade := false

	authRecord, err := s.user.GetUserAuthByUserIDAndType(ctx, user.ID, model.AuthTypePassword)
	if err != nil {
		return "", err
	}
	if authRecord != nil {
		matched, needUpgrade = verifyPassword(password, authRecord.Secret)
		if matched && needUpgrade {
			if upgradeErr := s.upsertPasswordAuth(ctx, user, account, password); upgradeErr != nil {
				log.Warnf("upgrade userauth password hash failed for user=%d: %v", user.ID, upgradeErr)
			}
		}
	} else {
		matched, needUpgrade = verifyPassword(password, user.Password)
		if matched {
			if upgradeErr := s.upsertPasswordAuth(ctx, user, account, password); upgradeErr != nil {
				log.Warnf("migrate legacy password to userauth failed for user=%d: %v", user.ID, upgradeErr)
			}
			if needUpgrade {
				log.Infof("legacy password upgraded for user=%d", user.ID)
			}
		}
	}

	if !matched {
		return "", kerrors.Unauthorized(auth.Reason, "invalid credentials")
	}

	return s.issueToken(ctx, user)
}

// Register creates a new local account and password credential.
func (s *SysUserUseCase) Register(ctx context.Context, account, password, email, phone, sex string) (string, error) {
	account = strings.TrimSpace(account)
	password = strings.TrimSpace(password)
	if account == "" || password == "" {
		return "", kerrors.BadRequest("INVALID_ARGUMENT", "account or password is empty")
	}

	existing, err := s.user.GetByAccount(ctx, account)
	if err != nil {
		return "", err
	}
	if existing != nil {
		return "", kerrors.BadRequest("ACCOUNT_EXISTS", "account already exists")
	}

	newUser := &model.SysUser{
		Type:        model.SysUserType_Guest,
		Account:     account,
		Name:        account,
		Email:       strings.TrimSpace(email),
		Phonenumber: strings.TrimSpace(phone),
		Sex:         strings.TrimSpace(sex),
		Status:      models.Status_Enabled,
	}
	newUser.New()
	if err = s.user.Create(ctx, newUser); err != nil {
		return "", err
	}

	if err = s.upsertPasswordAuth(ctx, newUser, account, password); err != nil {
		return "", err
	}

	return s.issueToken(ctx, newUser)
}

// LoginByQQ logs in or auto-registers a user by QQ openid.
func (s *SysUserUseCase) LoginByQQ(ctx context.Context, openID, nickname, avatar string) (string, *model.SysUser, error) {
	openID = strings.TrimSpace(openID)
	if openID == "" {
		return "", nil, kerrors.BadRequest("INVALID_ARGUMENT", "qq openid is empty")
	}

	record, err := s.user.GetUserAuthByTypeAndIdentifier(ctx, model.AuthTypeQQ, openID)
	if err != nil {
		return "", nil, err
	}

	var user *model.SysUser
	if record != nil {
		user, err = s.user.GetById(ctx, record.UserID)
		if err != nil {
			return "", nil, err
		}
	}

	if user == nil {
		account, accErr := s.buildQQAccount(ctx, openID)
		if accErr != nil {
			return "", nil, accErr
		}

		name := strings.TrimSpace(nickname)
		if name == "" {
			name = account
		}

		newUser := &model.SysUser{
			Type:    model.SysUserType_Guest,
			Account: account,
			Name:    name,
			Avatar:  strings.TrimSpace(avatar),
			Status:  models.Status_Enabled,
		}
		newUser.New()
		if err = s.user.Create(ctx, newUser); err != nil {
			return "", nil, err
		}

		authRecord := &model.SysUserAuth{
			UserID:     newUser.ID,
			AuthType:   model.AuthTypeQQ,
			Identifier: openID,
			Secret:     "",
		}
		authRecord.New()
		if err = s.user.CreateUserAuth(ctx, authRecord); err != nil {
			return "", nil, err
		}
		user = newUser
	}

	token, err := s.issueToken(ctx, user)
	if err != nil {
		return "", nil, err
	}
	return token, user, nil
}

// GetInfo returns profile from token context.
func (s *SysUserUseCase) GetInfo(ctx context.Context) (*model.SysUser, error) {
	uid, _ := ctx.Value(auth.UserId).(int64)
	user, err := s.user.GetById(ctx, uid)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Logout keeps current stateless behavior.
func (s *SysUserUseCase) Logout(ctx context.Context) error {
	return nil
}

func (s *SysUserUseCase) issueToken(ctx context.Context, user *model.SysUser) (string, error) {
	token, err := s.auth.NewToken(ctx, user.ID, user.Account, user.Phonenumber)
	if err != nil {
		log.Errorf("create token failed: %v", err)
		return "", kerrors.InternalServer("TOKEN_CREATE_FAILED", "create token failed")
	}
	return token, nil
}

func (s *SysUserUseCase) upsertPasswordAuth(ctx context.Context, user *model.SysUser, identifier, password string) error {
	hash, err := hashPassword(password)
	if err != nil {
		return err
	}

	identifier = strings.TrimSpace(identifier)
	if identifier == "" {
		identifier = user.Account
	}

	authRecord, err := s.user.GetUserAuthByUserIDAndType(ctx, user.ID, model.AuthTypePassword)
	if err != nil {
		return err
	}

	if authRecord == nil {
		authRecord = &model.SysUserAuth{
			UserID:     user.ID,
			AuthType:   model.AuthTypePassword,
			Identifier: identifier,
			Secret:     hash,
		}
		authRecord.New()
		return s.user.CreateUserAuth(ctx, authRecord)
	}

	authRecord.Identifier = identifier
	authRecord.Secret = hash
	return s.user.UpdateUserAuth(ctx, authRecord)
}

func (s *SysUserUseCase) buildQQAccount(ctx context.Context, openID string) (string, error) {
	suffix := openID
	if len(suffix) > 8 {
		suffix = suffix[:8]
	}
	base := "qq_" + suffix
	account := base

	for i := 0; i < 5; i++ {
		existing, err := s.user.GetByAccount(ctx, account)
		if err != nil {
			return "", err
		}
		if existing == nil {
			return account, nil
		}
		randSuffix, err := randomHex(3)
		if err != nil {
			return "", err
		}
		account = base + "_" + randSuffix
	}
	return "", kerrors.InternalServer("QQ_ACCOUNT_CREATE_FAILED", "failed to allocate qq account")
}

func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func verifyPassword(input, stored string) (matched bool, needUpgrade bool) {
	stored = strings.TrimSpace(stored)
	if stored == "" {
		return false, false
	}
	if err := bcrypt.CompareHashAndPassword([]byte(stored), []byte(input)); err == nil {
		return true, false
	}
	if input == stored {
		return true, true
	}
	return false, false
}

func randomHex(size int) (string, error) {
	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}


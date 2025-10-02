package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	// ErrUserNotFound is user not found.
	ErrUserNotFound = errors.NotFound("USER_NOT_FOUND", "user not found")
)

// User is a User model.
type User struct {
	ID           uint64
	Username     string
	Email        string
	PasswordHash string
	Role         string
	Status       int
	Avatar       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	LastLoginAt  *time.Time
}

// UserRepo is a User repo.
type UserRepo interface {
	Save(context.Context, *User) (*User, error)
	Update(context.Context, *User) (*User, error)
	FindByID(context.Context, uint64) (*User, error)
	FindByUsername(context.Context, string) (*User, error)
	UpdateLastLoginAt(context.Context, uint64, *time.Time) error
}

// UserUsecase is a User usecase.
type UserUsecase struct {
	repo UserRepo
	log  *log.Helper
}

// NewUserUsecase new a User usecase.
func NewUserUsecase(repo UserRepo, logger log.Logger) *UserUsecase {
	return &UserUsecase{repo: repo, log: log.NewHelper(logger)}
}

// CreateUser creates a User, and returns the new User.
func (uc *UserUsecase) CreateUser(ctx context.Context, u *User) (*User, error) {
	uc.log.WithContext(ctx).Infof("CreateUser: %v", u.Username)
	return uc.repo.Save(ctx, u)
}

// UpdateUser updates a User, and returns the updated User.
func (uc *UserUsecase) UpdateUser(ctx context.Context, u *User) (*User, error) {
	uc.log.WithContext(ctx).Infof("UpdateUser: %v", u.Username)
	return uc.repo.Update(ctx, u)
}

// GetUserByID gets a User by ID.
func (uc *UserUsecase) GetUserByID(ctx context.Context, id uint64) (*User, error) {
	uc.log.WithContext(ctx).Infof("GetUserByID: %v", id)
	return uc.repo.FindByID(ctx, id)
}

// GetUserByUsername gets a User by username.
func (uc *UserUsecase) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	uc.log.WithContext(ctx).Infof("GetUserByUsername: %v", username)
	return uc.repo.FindByUsername(ctx, username)
}

// UpdateLastLoginAt updates the last login time of a User.
func (uc *UserUsecase) UpdateLastLoginAt(ctx context.Context, id uint64, lastLoginAt *time.Time) error {
	uc.log.WithContext(ctx).Infof("UpdateLastLoginAt: %v", id)
	return uc.repo.UpdateLastLoginAt(ctx, id, lastLoginAt)
}
package data

import (
	"context"
	"time"

	"github.com/aichat/backend/internal/biz"
	"gorm.io/gorm"
)

// User model
type User struct {
	ID           uint64 `gorm:"primaryKey"`
	Username     string `gorm:"uniqueIndex;size:50"`
	Email        string `gorm:"uniqueIndex;size:100"`
	PasswordHash string `gorm:"size:255"`
	Role         string `gorm:"size:20"`
	Status       int    `gorm:"default:1"`
	Avatar       string `gorm:"size:255"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	LastLoginAt  *time.Time
}

// TableName overrides the table name
func (User) TableName() string {
	return "users"
}

// userRepo implements the biz.UserRepo interface
type userRepo struct {
	data *Data
}

// NewUserRepo creates a new user repository
func NewUserRepo(data *Data) biz.UserRepo {
	return &userRepo{
		data: data,
	}
}

// Save saves a user
func (r *userRepo) Save(ctx context.Context, user *biz.User) (*biz.User, error) {
	u := &User{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Role:         user.Role,
		Status:       user.Status,
		Avatar:       user.Avatar,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		LastLoginAt:  user.LastLoginAt,
	}
	
	if err := r.data.DB().Create(u).Error; err != nil {
		return nil, err
	}
	
	user.ID = u.ID
	return user, nil
}

// Update updates a user
func (r *userRepo) Update(ctx context.Context, user *biz.User) (*biz.User, error) {
	u := &User{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Role:         user.Role,
		Status:       user.Status,
		Avatar:       user.Avatar,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		LastLoginAt:  user.LastLoginAt,
	}
	
	if err := r.data.DB().Save(u).Error; err != nil {
		return nil, err
	}
	
	return &biz.User{
		ID:           u.ID,
		Username:     u.Username,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Role:         u.Role,
		Status:       u.Status,
		Avatar:       u.Avatar,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
		LastLoginAt:  u.LastLoginAt,
	}, nil
}

// FindByID finds a user by ID
func (r *userRepo) FindByID(ctx context.Context, id uint64) (*biz.User, error) {
	var u User
	if err := r.data.DB().Where("id = ?", id).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, biz.ErrUserNotFound
		}
		return nil, err
	}
	
	return &biz.User{
		ID:           u.ID,
		Username:     u.Username,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Role:         u.Role,
		Status:       u.Status,
		Avatar:       u.Avatar,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
		LastLoginAt:  u.LastLoginAt,
	}, nil
}

// FindByUsername finds a user by username
func (r *userRepo) FindByUsername(ctx context.Context, username string) (*biz.User, error) {
	var u User
	if err := r.data.DB().Where("username = ?", username).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, biz.ErrUserNotFound
		}
		return nil, err
	}
	
	return &biz.User{
		ID:           u.ID,
		Username:     u.Username,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Role:         u.Role,
		Status:       u.Status,
		Avatar:       u.Avatar,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
		LastLoginAt:  u.LastLoginAt,
	}, nil
}

// UpdateLastLoginAt updates the last login time of a user
func (r *userRepo) UpdateLastLoginAt(ctx context.Context, id uint64, lastLoginAt *time.Time) error {
	return r.data.DB().Model(&User{}).Where("id = ?", id).Update("last_login_at", lastLoginAt).Error
}
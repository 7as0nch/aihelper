package service

import (
	"context"
	"time"

	pb "github.com/aichat/backend/api/aichat/v1"
	"github.com/aichat/backend/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	
	uc  *biz.UserUsecase
	log *log.Helper
}

func NewUserService(uc *biz.UserUsecase, logger log.Logger) *UserService {
	return &UserService{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

// Login implements aichat.v1.UserServiceServer.
func (s *UserService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	s.log.WithContext(ctx).Infof("Login: %v", req.Username)
	
	// 根据用户名查找用户
	user, err := s.uc.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	
	// TODO: 验证密码
	
	// 更新最后登录时间
	now := time.Now()
	if err := s.uc.UpdateLastLoginAt(ctx, user.ID, &now); err != nil {
		s.log.WithContext(ctx).Errorf("Failed to update last login time: %v", err)
		// 不返回错误，因为登录本身是成功的
	}
	
	// TODO: 生成JWT token
	
	return &pb.LoginReply{
		Success: true,
		Message: "Login successful",
		Token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE3MTQ3MzQyMjF9.5cB6",
		Data: &pb.UserInfo{
			UserId:       user.ID,
			Username:     user.Username,
			Email:        user.Email,
			Role:         user.Role,
			CreatedAt:    user.CreatedAt.Format(time.RFC3339),
			LastLoginAt:  func() string {
				if user.LastLoginAt != nil {
					return user.LastLoginAt.Format(time.RFC3339)
				}
				return ""
			}(),
		},
	}, nil
}

// Register implements aichat.v1.UserServiceServer.
func (s *UserService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterReply, error) {
	s.log.WithContext(ctx).Infof("Register: %v", req.Username)
	
	// 创建用户
	user := &biz.User{
		Username:     req.Username,
		Email:        req.Email,
		// TODO: 密码加密
		PasswordHash: req.Password,
		Role:         "user",
		Status:       1,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	
	createdUser, err := s.uc.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	
	return &pb.RegisterReply{
		Success: true,
		Message: "User registered successfully",
		User: &pb.UserInfo{
			UserId:    createdUser.ID,
			Username:  createdUser.Username,
			Email:     createdUser.Email,
			Role:      createdUser.Role,
			CreatedAt: createdUser.CreatedAt.Format(time.RFC3339),
		},
	}, nil
}

// GetUserInfo implements aichat.v1.UserServiceServer.
func (s *UserService) GetUserInfo(ctx context.Context, req *pb.GetUserInfoRequest) (*pb.GetUserInfoReply, error) {
	s.log.WithContext(ctx).Infof("GetUserInfo: %v", req.UserId)
	
	// 根据用户ID获取用户信息
	user, err := s.uc.GetUserByID(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	
	return &pb.GetUserInfoReply{
		Success: true,
		Data: &pb.UserInfo{
			UserId:       user.ID,
			Username:     user.Username,
			Email:        user.Email,
			Role:         user.Role,
			CreatedAt:    user.CreatedAt.Format(time.RFC3339),
			LastLoginAt:  func() string {
				if user.LastLoginAt != nil {
					return user.LastLoginAt.Format(time.RFC3339)
				}
				return ""
			}(),
		},
	}, nil
}
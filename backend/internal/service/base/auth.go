package base

import (
	"context"
	"sync"

	pb "github.com/example/aichat/backend/api/base"
	"github.com/example/aichat/backend/internal/biz/base"
	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthService struct {
	pb.UnimplementedAuthServer
	user *base.SysUserUseCase

	qqStateMu sync.Mutex
	qqState   map[string]qqStateItem
}

func NewAuthService(user *base.SysUserUseCase) *AuthService {
	return &AuthService{
		user:    user,
		qqState: make(map[string]qqStateItem),
	}
}

func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	log.Infof("login request for user: %s", req.Username)
	token, err := s.user.Login(ctx, req.Username, req.Password)
	if err != nil {
		return nil, err
	}
	return &pb.LoginReply{
		AccessToken: token,
	}, nil
}

// getInfo
func (s *AuthService) GetInfo(ctx context.Context, _ *emptypb.Empty) (*pb.GetInfoReply, error) {
	user, err := s.user.GetInfo(ctx)
	if err != nil {
		return nil, err
	}
	if user == nil || user.ID == 0 {
		return nil, kerrors.New(401, "USER_NOT_FOUND", "user not found")
	}
	return &pb.GetInfoReply{
		User: &pb.User{
			UserId:   user.ID,
			NickName: user.Account,
			UserName: user.Name,
			Avatar:   user.Avatar,
		},
		Roles: []string{"superadmin"},
		Permissions: []string{"*:*:*"},
	}, nil
}

// logout
func (s *AuthService) Logout(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, s.user.Logout(ctx)
}

// Register
func (s *AuthService) Register(ctx context.Context, req *pb.RegisterRequest) (*emptypb.Empty, error) {
	_, err := s.user.Register(ctx, req.Account, req.Password, req.Email, req.Phonenumber, req.Sex)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// UpdateProfile
func (s *AuthService) UpdateProfile(ctx context.Context, req *pb.User) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

// GetProfile
func (s *AuthService) GetProfile(ctx context.Context, req *emptypb.Empty) (*pb.GetInfoReply, error) {
	user, err := s.user.GetInfo(ctx)
	if err != nil {
		return nil, err
	}
	if user == nil || user.ID == 0 {
		return nil, kerrors.New(401, "USER_NOT_FOUND", "user not found")
	}
	return &pb.GetInfoReply{
		User: &pb.User{
			UserId:      user.ID,
			NickName:    user.Account,
			UserName:    user.Name,
			Avatar:      user.Avatar,
			Email:       user.Email,
			Phonenumber: user.Phonenumber,
		},
	}, nil
}

// UpdatePwd
func (s *AuthService) UpdatePwd(ctx context.Context, req *pb.UpdatePwdRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}


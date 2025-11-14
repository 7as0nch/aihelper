package base

import (
	"context"

	pb "github.com/example/aichat/backend/api/base"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthService struct {
	pb.UnimplementedAuthServer
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
func (s *AuthService) Logout(ctx context.Context, req *pb.LogoutRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

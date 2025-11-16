package base

import (
	"context"
	kerrors "github.com/go-kratos/kratos/v2/errors"
	pb "github.com/example/aichat/backend/api/base"
	"github.com/example/aichat/backend/internal/biz/base"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthService struct {
	pb.UnimplementedAuthServer
	user *base.SysUserUseCase
}

func NewAuthService(user *base.SysUserUseCase) *AuthService {
	return &AuthService{
		user: user,
	}
}

func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	log.Infof("用户登录 %v", req.Username)
	token, err := s.user.Login(ctx, req.Username)
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
		return nil, kerrors.New(401, "用户不存在", "用户不存在")
	}
	log.Infof("获取用户信息 %#v", user)
	return &pb.GetInfoReply{
		User: &pb.User{
			UserName: user.Name,
			NickName: user.Account,
		},
		Roles: []string{
			"superadmin",
		},
		Permissions: []string{"*:*:*"},
	}, nil
}

// logout
// NOTE: 这里笔记备注：虽然请求参数为空，但仍然需要定义一个空的请求类型，以符合 gRPC 的方法签名要求。前端调用时，传递一个空的请求对象即可。
func (s *AuthService) Logout(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, s.user.Logout(ctx)
}

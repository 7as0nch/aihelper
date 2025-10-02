package server

import (
	v1 "github.com/aichat/backend/api/aichat/v1"
	"github.com/aichat/backend/internal/service"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(opts ...grpc.ServerOption) *grpc.Server {
	// Define server options
	var serverOpts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
		),
	}
	
	// Append additional options
	serverOpts = append(serverOpts, opts...)
	
	// Create gRPC server
	srv := grpc.NewServer(serverOpts...)
	
	// Register services
	return srv
}

// RegisterGRPCServices registers all services to the gRPC server.
func RegisterGRPCServices(srv *grpc.Server, user *service.UserService) {
	// 注册Greeter服务
	// v1.RegisterGreeterServer(srv, greeter)
	
	// 注册UserService
	v1.RegisterUserServiceServer(srv, user)
}

package server

import (
	v1 "github.com/aichat/backend/api/aichat/v1"
	"github.com/aichat/backend/internal/service"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(opts ...http.ServerOption) *http.Server {
	// Define server options
	var serverOpts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
	}
	
	// Append additional options
	serverOpts = append(serverOpts, opts...)
	
	// Create HTTP server
	srv := http.NewServer(serverOpts...)
	
	// Register services
	return srv
}

// RegisterHTTPServices registers all services to the HTTP server.
func RegisterHTTPServices(srv *http.Server, user *service.UserService) {
	// 注册Greeter服务
	// v1.RegisterGreeterHTTPServer(srv, greeter)
	
	// 注册UserService
	v1.RegisterUserServiceHTTPServer(srv, user)
}

package server

import (
	basepb "github.com/example/aichat/backend/api/base"
	chatv1 "github.com/example/aichat/backend/api/chat/v1"
	userfeedbackv1 "github.com/example/aichat/backend/api/userfeedback/v1"
	"github.com/example/aichat/backend/internal/conf"
	"github.com/example/aichat/backend/internal/service"
	"github.com/example/aichat/backend/internal/service/base"
	"github.com/example/aichat/backend/pkg/auth"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	kratoshttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/swagger-api/openapiv2"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.24.0"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server,
	userFeedback *service.UserFeedbackService,
	chat *service.ChatService,
	authServ *base.AuthService,
	authRepo auth.AuthRepo,
	system *base.SystemService,
	logg log.Logger) *kratoshttp.Server {

	// 初始化 tracer provider（开发环境使用采样率100%，生产环境可调整）
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(1.0))),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("aichat-backend-http"),
			semconv.DeploymentEnvironmentKey.String("development"),
		)),
	)

	// 设置全局 tracer provider
	otel.SetTracerProvider(tp)
	var opts = []kratoshttp.ServerOption{
		kratoshttp.Middleware(
			recovery.Recovery(),
			logging.Server(logg),
			tracing.Server(), // 启用分布式追踪中间件
			selector.Server(
				auth.NewHeaderServer(),
				authRepo.Server()).
				Match(auth.NewWhiteListMatcher(map[string]bool{
					basepb.OperationAuthLogin:   true,
				})).Build(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, kratoshttp.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, kratoshttp.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, kratoshttp.Timeout(c.Http.Timeout.AsDuration()))
	}
	opts = append(opts,
		kratoshttp.ResponseEncoder(auth.DefaultResponseEncoder),
		kratoshttp.ErrorEncoder(auth.DefaultErrorEncoder))
	srv := kratoshttp.NewServer(opts...)
	userfeedbackv1.RegisterUserFeedbackHTTPServer(srv, userFeedback)
	chatv1.RegisterChatHTTPServer(srv, chat)
	basepb.RegisterAuthHTTPServer(srv, authServ)
	basepb.RegisterSystemHTTPServer(srv, system)
	srv.HandleFunc("/chat/send", chat.SSEHandler)
	srv.HandlePrefix("/q/", openapiv2.NewHandler())
	return srv
}

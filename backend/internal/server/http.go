package server

import (
	"context"
	"net/http"

	basepb "github.com/example/aichat/backend/api/base"
	aipb "github.com/example/aichat/backend/api/ai"
	chatv1 "github.com/example/aichat/backend/api/chat/v1"
	"github.com/example/aichat/backend/internal/conf"
	"github.com/example/aichat/backend/internal/service"
	"github.com/example/aichat/backend/internal/service/ai"
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
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server,
	chat *service.ChatService,
	authServ *base.AuthService,
	authRepo auth.AuthRepo,
	system *base.SystemService,
	tracker *base.TrackerService,
	aiServ *ai.AIService,
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

	// 设置全局 propagator,使 tracing.Server() 能够从 HTTP headers 提取 trace context
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	var opts = []kratoshttp.ServerOption{
		kratoshttp.Middleware(
			recovery.Recovery(),
			// 先提取 ClientID,这样后续的 logging.Server 可以记录它
			selector.Server(
				auth.NewHeaderServer()).
				Match(func(ctx context.Context, operation string) bool {
					// 对所有请求都执行 NewHeaderServer,提取 ClientID
					return true
				}).Build(),
			logging.Server(logg),  // 此时 ClientID 已经在 context 中
			tracing.Server(),      // 启用分布式追踪中间件,通过 propagator 提取 trace context
			auth.MiddlewareCors(), // 跨域中间件,只对特定接口开放
			// JWT 验证保持在最后,使用原有的白名单逻辑
			selector.Server(
				authRepo.Server()).
				Match(auth.NewWhiteListMatcher(map[string]bool{
					basepb.OperationAuthLogin:    true,
					basepb.OperationTrackerBatch: true,
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
		kratoshttp.ErrorEncoder(auth.DefaultErrorEncoder),
		kratoshttp.RequestDecoder(func(r *http.Request, v interface{}) error {
			// 处理text/plain类型请求
			if r.Header.Get("Content-Type") == "text/plain; charset=utf-8" {
				// 将Content-Type设置为application/json，以便使用默认的JSON解码器
				r.Header.Set("Content-Type", "application/json")
			}
			// 使用默认的请求解码器
			return kratoshttp.DefaultRequestDecoder(r, v)
		}))
	srv := kratoshttp.NewServer(opts...)
	chatv1.RegisterChatHTTPServer(srv, chat)
	basepb.RegisterAuthHTTPServer(srv, authServ)
	basepb.RegisterSystemHTTPServer(srv, system)
	basepb.RegisterTrackerHTTPServer(srv, tracker)
	aipb.RegisterAIHTTPServer(srv, aiServ)
	srv.HandleFunc("/chat/send", chat.SSEHandler)
	srv.HandlePrefix("/q/", openapiv2.NewHandler())
	return srv
}

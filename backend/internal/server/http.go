package server

import (
	"context"
	"net/http"
	_ "net/http/pprof"

	aipb "github.com/example/aichat/backend/api/ai"
	basepb "github.com/example/aichat/backend/api/base"
	chatv1 "github.com/example/aichat/backend/api/chat/v1"
	"github.com/example/aichat/backend/internal/conf"
	"github.com/example/aichat/backend/internal/service"
	"github.com/example/aichat/backend/internal/service/ai"
	"github.com/example/aichat/backend/internal/service/base"
	"github.com/example/aichat/backend/pkg/auth"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	kratoshttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/swagger-api/openapiv2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

func NewHTTPServer(c *conf.Server,
	chat *service.ChatService,
	authServ *base.AuthService,
	authRepo auth.AuthRepo,
	system *base.SystemService,
	tracker *base.TrackerService,
	betaApplication *base.BetaApplicationService,
	aiServ *ai.AIService,
	logg log.Logger) *kratoshttp.Server {
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(1.0))),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("aichat-backend-http"),
			semconv.DeploymentEnvironmentKey.String("development"),
		)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	var opts = []kratoshttp.ServerOption{
		kratoshttp.Middleware(
			metrics.Server(),
			recovery.Recovery(),
			selector.Server(auth.NewHeaderServer()).Match(func(ctx context.Context, operation string) bool {
				return true
			}).Build(),
			logging.Server(logg),
			tracing.Server(),
			auth.MiddlewareCors(),
			selector.Server(authRepo.Server()).Match(auth.NewWhiteListMatcher(map[string]bool{
				basepb.OperationAuthLogin:             true,
				basepb.OperationTrackerBatch:          true,
				basepb.OperationBetaApplicationCreateApplication: true,
				"/auth/qq/login":                    true,
				"/auth/qq/callback":                 true,
				"GET /auth/qq/login":                true,
				"GET /auth/qq/callback":             true,
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
			if r.Header.Get("Content-Type") == "text/plain; charset=utf-8" {
				r.Header.Set("Content-Type", "application/json")
			}
			return kratoshttp.DefaultRequestDecoder(r, v)
		}))
	srv := kratoshttp.NewServer(opts...)
	chatv1.RegisterChatHTTPServer(srv, chat)
	basepb.RegisterAuthHTTPServer(srv, authServ)
	srv.HandleFunc("/auth/qq/login", authServ.HandleQQLogin)
	srv.HandleFunc("/auth/qq/callback", authServ.HandleQQCallback)
	basepb.RegisterSystemHTTPServer(srv, system)
	basepb.RegisterTrackerHTTPServer(srv, tracker)
	basepb.RegisterBetaApplicationHTTPServer(srv, betaApplication)
	aipb.RegisterAIHTTPServer(srv, aiServ)

	srv.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	srv.HandleFunc("/chat/send", func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token != "" {
			if claims, err := authRepo.CheckToken(r.Context(), token); err == nil && claims != nil {
				ctx := context.WithValue(r.Context(), auth.UserId, int64(claims.UserId))
				r = r.WithContext(ctx)
			}
		}
		chat.SSEHandler(w, r)
	})
	srv.Handle("/metrics", promhttp.Handler())
	srv.HandlePrefix("/debug/pprof/", http.DefaultServeMux)
	srv.HandlePrefix("/q/", openapiv2.NewHandler())
	return srv
}

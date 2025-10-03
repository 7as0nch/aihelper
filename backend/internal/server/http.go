package server

import (
	chatv1 "github.com/example/aichat/backend/api/chat/v1"
	reportv1 "github.com/example/aichat/backend/api/report/v1"
	userfeedbackv1 "github.com/example/aichat/backend/api/userfeedback/v1"
	workflowv1 "github.com/example/aichat/backend/api/workflow/v1"
	v1 "github.com/example/aichat/backend/api/helloworld/v1"
	"github.com/example/aichat/backend/internal/conf"
	"github.com/example/aichat/backend/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, 
	greeter *service.GreeterService, 
	userFeedback *service.UserFeedbackService,
	chat *service.ChatService,
	workflow *service.WorkflowAPIService,
	report *service.ReportService,
	logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	v1.RegisterGreeterHTTPServer(srv, greeter)
	userfeedbackv1.RegisterUserFeedbackHTTPServer(srv, userFeedback)
	chatv1.RegisterChatHTTPServer(srv, chat)
	workflowv1.RegisterWorkflowHTTPServer(srv, workflow)
	reportv1.RegisterReportHTTPServer(srv, report)
	return srv
}

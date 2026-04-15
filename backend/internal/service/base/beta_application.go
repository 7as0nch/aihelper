package base

import (
	"context"
	"fmt"
	"net"
	"net/smtp"
	"regexp"
	"strconv"
	"strings"

	pb "github.com/example/aichat/backend/api/base"
	bizbase "github.com/example/aichat/backend/internal/biz/base"
	"github.com/example/aichat/backend/internal/conf"
	"github.com/example/aichat/backend/models/generator/model"
	lib "github.com/example/aichat/backend/pkg/lib"
	kerrors "github.com/go-kratos/kratos/v2/errors"
	"go.uber.org/zap"
)

type BetaApplicationService struct {
	pb.UnimplementedBetaApplicationServer
	beta *bizbase.BetaApplicationUseCase
}

type betaMailer struct {
	mail *conf.Beta_Mail
	log  *zap.Logger
}

func NewBetaApplicationService(beta *bizbase.BetaApplicationUseCase) *BetaApplicationService {
	return &BetaApplicationService{beta: beta}
}

func NewBetaApplicationNotifier(bootstrap *conf.Bootstrap, log *zap.Logger) bizbase.BetaApplicationNotifier {
	var mailConf *conf.Beta_Mail
	if bootstrap != nil && bootstrap.Beta != nil {
		mailConf = bootstrap.Beta.Mail
	}
	return &betaMailer{mail: mailConf, log: log}
}

func (s *BetaApplicationService) CreateApplication(ctx context.Context, req *pb.CreateBetaApplicationRequest) (*pb.CreateBetaApplicationReply, error) {
	if err := validateCreateBetaApplicationRequest(req); err != nil {
		return nil, err
	}

	application, err := s.beta.Submit(
		ctx,
		req.ProductInterest,
		req.ContactType,
		req.ContactValue,
		req.UseCase,
		req.Note,
		req.SourcePage,
		"",
		"",
	)
	if err != nil {
		return nil, err
	}

	return &pb.CreateBetaApplicationReply{
		Id:         application.ID,
		Status:     application.Status,
		MailStatus: application.MailStatus,
	}, nil
}

func validateCreateBetaApplicationRequest(req *pb.CreateBetaApplicationRequest) error {
	productInterest := strings.TrimSpace(req.ProductInterest)
	contactType := strings.TrimSpace(req.ContactType)
	contactValue := strings.TrimSpace(req.ContactValue)
	useCase := strings.TrimSpace(req.UseCase)

	supportedProducts := map[string]bool{
		"litechat": true,
		"aicook": true,
		"tech-sandbox": true,
	}
	if !supportedProducts[productInterest] {
		return kerrors.BadRequest("INVALID_PRODUCT", "请选择有效的产品方向")
	}

	supportedContactTypes := map[string]bool{
		"email": true,
		"qq":    true,
	}
	if !supportedContactTypes[contactType] {
		return kerrors.BadRequest("INVALID_CONTACT_TYPE", "请选择有效的联系方式")
	}
	if contactValue == "" {
		return kerrors.BadRequest("EMPTY_CONTACT_VALUE", "请填写联系方式")
	}
	if len([]rune(contactValue)) > 128 {
		return kerrors.BadRequest("CONTACT_VALUE_TOO_LONG", "联系方式过长")
	}
	if useCase == "" {
		return kerrors.BadRequest("EMPTY_USE_CASE", "请填写使用场景")
	}
	if len([]rune(useCase)) > 1000 {
		return kerrors.BadRequest("USE_CASE_TOO_LONG", "使用场景内容过长")
	}
	if len([]rune(req.Note)) > 1200 {
		return kerrors.BadRequest("NOTE_TOO_LONG", "补充说明内容过长")
	}

	if contactType == "email" && !lib.VerifyEmailFormat(contactValue) {
		return kerrors.BadRequest("INVALID_EMAIL", "邮箱格式不正确")
	}
	if contactType == "qq" {
		matched, _ := regexp.MatchString(`^[1-9][0-9]{4,11}$`, contactValue)
		if !matched {
			return kerrors.BadRequest("INVALID_QQ", "QQ 号格式不正确")
		}
	}

	return nil
}

func (m *betaMailer) Notify(_ context.Context, application *model.BetaApplication) (string, error) {
	if m.mail == nil {
		m.log.Warn("beta application mail skipped because yaml mail config is missing")
		return model.BetaApplicationMailSkipped, nil
	}

	host := strings.TrimSpace(m.mail.Host)
	username := strings.TrimSpace(m.mail.Username)
	password := strings.TrimSpace(m.mail.Password)
	to := strings.TrimSpace(m.mail.To)
	from := strings.TrimSpace(m.mail.From)
	port := int(m.mail.Port)
	if port == 0 {
		port = 587
	}
	if host == "" || username == "" || password == "" {
		m.log.Warn("beta application mail skipped because yaml smtp config is incomplete")
		return model.BetaApplicationMailSkipped, nil
	}
	if from == "" {
		from = username
	}
	if to == "" {
		to = from
	}

	subject := fmt.Sprintf("[LiteChat Beta] %s", application.ProductInterest)
	body := strings.Join([]string{
		fmt.Sprintf("申请ID: %d", application.ID),
		fmt.Sprintf("产品方向: %s", application.ProductInterest),
		fmt.Sprintf("联系方式: %s / %s", application.ContactType, application.ContactValue),
		fmt.Sprintf("使用场景: %s", application.UseCase),
		fmt.Sprintf("补充说明: %s", firstNonEmpty(application.Note, "-")),
		fmt.Sprintf("来源页面: %s", firstNonEmpty(application.SourcePage, "-")),
		fmt.Sprintf("远端地址: %s", firstNonEmpty(application.RemoteAddr, "-")),
		fmt.Sprintf("User-Agent: %s", firstNonEmpty(application.UserAgent, "-")),
	}, "\r\n")
	message := []byte("To: " + to + "\r\n" +
		"From: " + from + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n\r\n" + body)

	addr := net.JoinHostPort(host, strconv.Itoa(port))
	authInfo := smtp.PlainAuth("", username, password, host)
	if err := smtp.SendMail(addr, authInfo, from, []string{to}, message); err != nil {
		return model.BetaApplicationMailFailed, err
	}
	return model.BetaApplicationMailSent, nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			return trimmed
		}
	}
	return ""
}

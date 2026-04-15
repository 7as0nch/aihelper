package base

import (
	"context"
	"strings"

	"github.com/example/aichat/backend/models/generator/model"
	"go.uber.org/zap"
)

type BetaApplicationRepo interface {
	Create(ctx context.Context, application *model.BetaApplication) error
	UpdateMailResult(ctx context.Context, id int64, status string, mailError string) error
}

type BetaApplicationNotifier interface {
	Notify(ctx context.Context, application *model.BetaApplication) (string, error)
}

type BetaApplicationUseCase struct {
	repo     BetaApplicationRepo
	notifier BetaApplicationNotifier
	log      *zap.Logger
}

func NewBetaApplicationUseCase(repo BetaApplicationRepo, notifier BetaApplicationNotifier, log *zap.Logger) *BetaApplicationUseCase {
	return &BetaApplicationUseCase{repo: repo, notifier: notifier, log: log}
}

func (uc *BetaApplicationUseCase) Submit(
	ctx context.Context,
	productInterest string,
	contactType string,
	contactValue string,
	useCase string,
	note string,
	sourcePage string,
	userAgent string,
	remoteAddr string,
) (*model.BetaApplication, error) {
	application := &model.BetaApplication{
		ProductInterest: strings.TrimSpace(productInterest),
		ContactType:     strings.TrimSpace(contactType),
		ContactValue:    strings.TrimSpace(contactValue),
		UseCase:         strings.TrimSpace(useCase),
		Note:            strings.TrimSpace(note),
		SourcePage:      strings.TrimSpace(sourcePage),
		UserAgent:       strings.TrimSpace(userAgent),
		RemoteAddr:      strings.TrimSpace(remoteAddr),
		Status:          model.BetaApplicationStatusSubmitted,
		MailStatus:      model.BetaApplicationMailPending,
	}
	application.New()

	if err := uc.repo.Create(ctx, application); err != nil {
		return nil, err
	}

	mailStatus, err := uc.notifier.Notify(ctx, application)
	if mailStatus == "" {
		mailStatus = model.BetaApplicationMailSkipped
	}
	mailError := ""
	if err != nil {
		mailError = err.Error()
		uc.log.Error("beta application notify failed", zap.Int64("id", application.ID), zap.Error(err))
	}
	if updateErr := uc.repo.UpdateMailResult(ctx, application.ID, mailStatus, mailError); updateErr != nil {
		uc.log.Error("beta application update mail status failed", zap.Int64("id", application.ID), zap.Error(updateErr))
	} else {
		application.MailStatus = mailStatus
		application.MailError = mailError
	}

	return application, nil
}

package data

import (
	"context"

	bizbase "github.com/example/aichat/backend/internal/biz/base"
	"github.com/example/aichat/backend/internal/db"
	"github.com/example/aichat/backend/models/generator/model"
	"go.uber.org/zap"
)

type betaApplicationRepo struct {
	db  db.DataRepo
	log *zap.Logger
}

func NewBetaApplicationRepo(dbRepo db.DataRepo, log *zap.Logger) bizbase.BetaApplicationRepo {
	// if err := dbRepo.GetDB().AutoMigrate(&model.BetaApplication{}); err != nil {
	// 	log.Panic("auto migrate beta application failed", zap.Error(err))
	// }
	return &betaApplicationRepo{db: dbRepo, log: log}
}

func (r *betaApplicationRepo) Create(ctx context.Context, application *model.BetaApplication) error {
	return r.db.DB(ctx).Create(application).Error
}

func (r *betaApplicationRepo) UpdateMailResult(ctx context.Context, id int64, status string, mailError string) error {
	return r.db.DB(ctx).Model(&model.BetaApplication{}).Where("id = ?", id).Updates(map[string]any{
		"mail_status": status,
		"mail_error":  mailError,
	}).Error
}

package biz

import (
	"context"
	"time"

	"go.uber.org/zap"
)

type UserFeedback struct {
	ID        uint64    `gorm:"primaryKey"`
	UserID    string    `gorm:"type:varchar(255);not null"`
	Content   string    `gorm:"type:text;not null"`
	Type      string    `gorm:"type:varchar(50);not null"`
	Status    string    `gorm:"type:varchar(50);default:'pending'"`
	Rating    int32     `gorm:"type:int;default:0"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}


type ChannelRoiParam struct {
	ChannelID int64 `json:"channel_id" jsonschema:"description=渠道id"`
	StartTime string `json:"start_time" jsonschema:"description=开始时间"`
	EndTime   string `json:"end_time" jsonschema:"description=结束时间"`
}
type ChannelRoi struct {
	// 渠道id
	ChannelID int64 `gorm:"type:bigint;not null" jsonschema:"description=渠道id"`
	// 半流程分发收益
	HalfProfit float64 `gorm:"type:decimal(10,2);not null" jsonschema:"description=半流程分发收益"`
	// 半流程分发净收益
	HalfNetProfit float64 `gorm:"type:decimal(10,2);not null" jsonschema:"description=半流程分发净收益"`
}

type UserFeedbackRepo interface {
	Create(ctx context.Context, feedback *UserFeedback) error
	Update(ctx context.Context, feedback *UserFeedback) error
	Delete(ctx context.Context, id uint64) error
	Get(ctx context.Context, id uint64) (*UserFeedback, error)
	List(ctx context.Context, userID string, page, pageSize int32) ([]*UserFeedback, int64, error)
	ListByStatus(ctx context.Context, status string, page, pageSize int32) ([]*UserFeedback, int64, error)
	GetChannelRoi(ctx context.Context, param *ChannelRoiParam) ([]*ChannelRoi, error)
}

type UserFeedbackUsecase struct {
	repo UserFeedbackRepo
	log  *zap.Logger
}

func NewUserFeedbackUsecase(repo UserFeedbackRepo, logger *zap.Logger) *UserFeedbackUsecase {
	return &UserFeedbackUsecase{
		repo: repo,
		log:  logger.Named("NewUserFeedbackUsecase"),
	}
}

func (uc *UserFeedbackUsecase) CreateFeedback(ctx context.Context, feedback *UserFeedback) (*UserFeedback, error) {
	uc.log.Info("Creating feedback for user", zap.String("userID", feedback.UserID))
	return feedback, uc.repo.Create(ctx, feedback)
}

func (uc *UserFeedbackUsecase) GetFeedback(ctx context.Context, id uint64) (*UserFeedback, error) {
	uc.log.Info("Getting feedback", zap.Uint64("id", id))
	return uc.repo.Get(ctx, id)
}

func (uc *UserFeedbackUsecase) ListFeedbacks(ctx context.Context, userID string, page, pageSize int32) ([]*UserFeedback, int64, error) {
	uc.log.Info("Listing feedbacks for user", zap.String("userID", userID), zap.Int32("page", page), zap.Int32("pageSize", pageSize))
	return uc.repo.List(ctx, userID, page, pageSize)
}

func (uc *UserFeedbackUsecase) UpdateFeedback(ctx context.Context, feedback *UserFeedback) (*UserFeedback, error) {
	uc.log.Info("Updating feedback", zap.Uint64("id", feedback.ID))
	return feedback, uc.repo.Update(ctx, feedback)
}

func (uc *UserFeedbackUsecase) DeleteFeedback(ctx context.Context, id uint64) error {
	uc.log.Info("Deleting feedback", zap.Uint64("id", id))
	return uc.repo.Delete(ctx, id)
}

func (uc *UserFeedbackUsecase) ListFeedbacksByStatus(ctx context.Context, status string, page, pageSize int32) ([]*UserFeedback, int64, error) {
	uc.log.Info("Listing feedbacks by status", zap.String("status", status), zap.Int32("page", page), zap.Int32("pageSize", pageSize))
	return uc.repo.ListByStatus(ctx, status, page, pageSize)
}

func (uc *UserFeedbackUsecase) GetChannelRoi(ctx context.Context, param *ChannelRoiParam) ([]*ChannelRoi, error) {
	uc.log.Info("Getting channel ROI", zap.Int64("channelID", param.ChannelID), zap.String("startTime", param.StartTime), zap.String("endTime", param.EndTime))
	return uc.repo.GetChannelRoi(ctx, param)
}
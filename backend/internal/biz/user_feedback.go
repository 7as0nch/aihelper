package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
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
	log  *log.Helper
}

func NewUserFeedbackUsecase(repo UserFeedbackRepo, logger log.Logger) *UserFeedbackUsecase {
	return &UserFeedbackUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (uc *UserFeedbackUsecase) CreateFeedback(ctx context.Context, feedback *UserFeedback) (*UserFeedback, error) {
	uc.log.Infof("Creating feedback for user: %s", feedback.UserID)
	return feedback, uc.repo.Create(ctx, feedback)
}

func (uc *UserFeedbackUsecase) GetFeedback(ctx context.Context, id uint64) (*UserFeedback, error) {
	uc.log.Infof("Getting feedback: %d", id)
	return uc.repo.Get(ctx, id)
}

func (uc *UserFeedbackUsecase) ListFeedbacks(ctx context.Context, userID string, page, pageSize int32) ([]*UserFeedback, int64, error) {
	uc.log.Infof("Listing feedbacks for user: %s, page: %d, pageSize: %d", userID, page, pageSize)
	return uc.repo.List(ctx, userID, page, pageSize)
}

func (uc *UserFeedbackUsecase) UpdateFeedback(ctx context.Context, feedback *UserFeedback) (*UserFeedback, error) {
	uc.log.Infof("Updating feedback: %d", feedback.ID)
	return feedback, uc.repo.Update(ctx, feedback)
}

func (uc *UserFeedbackUsecase) DeleteFeedback(ctx context.Context, id uint64) error {
	uc.log.Infof("Deleting feedback: %d", id)
	return uc.repo.Delete(ctx, id)
}

func (uc *UserFeedbackUsecase) ListFeedbacksByStatus(ctx context.Context, status string, page, pageSize int32) ([]*UserFeedback, int64, error) {
	uc.log.Infof("Listing feedbacks by status: %s, page: %d, pageSize: %d", status, page, pageSize)
	return uc.repo.ListByStatus(ctx, status, page, pageSize)
}

func (uc *UserFeedbackUsecase) GetChannelRoi(ctx context.Context, param *ChannelRoiParam) ([]*ChannelRoi, error) {
	uc.log.Infof("Getting channel ROI for channel: %d, startTime: %s, endTime: %s", param.ChannelID, param.StartTime, param.EndTime)
	return uc.repo.GetChannelRoi(ctx, param)
}
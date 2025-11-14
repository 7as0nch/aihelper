package data

import (
	"context"

	"github.com/example/aichat/backend/internal/biz"
	"go.uber.org/zap"

)

// UserFeedback 用户反馈数据模型
// type UserFeedback struct {
// 	ID        uint           `gorm:"primarykey" json:"id"`
// 	UserID    string         `gorm:"type:varchar(100);not null;index" json:"user_id"`
// 	Content   string         `gorm:"type:text;not null" json:"content"`
// 	Type      string         `gorm:"type:varchar(50);default:'general'" json:"type"`
// 	Status    string         `gorm:"type:varchar(50);default:'pending'" json:"status"`
// 	Rating    int            `gorm:"type:int;default:0" json:"rating"`
// 	CreatedAt time.Time      `json:"created_at"`
// 	UpdatedAt time.Time      `json:"updated_at"`
// 	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
// }

// UserFeedbackRepo 用户反馈数据仓库接口
// type UserFeedbackRepo interface {
// 	Create(ctx context.Context, feedback *UserFeedback) error
// 	Update(ctx context.Context, feedback *UserFeedback) error
// 	Delete(ctx context.Context, id uint) error
// 	Get(ctx context.Context, id uint) (*UserFeedback, error)
// 	List(ctx context.Context, userID string, page, pageSize int) ([]*UserFeedback, int64, error)
// 	ListByStatus(ctx context.Context, status string, page, pageSize int) ([]*UserFeedback, int64, error)
// }

// userFeedbackRepo 用户反馈数据仓库实现
type userFeedbackRepo struct {
	data *Data
	log  *zap.Logger
}

// NewUserFeedbackRepo 创建用户反馈仓库实例
func NewUserFeedbackRepo(data *Data, logger *zap.Logger) biz.UserFeedbackRepo {
	return &userFeedbackRepo{
		data: data,
		log:  logger,
	}
}

// Create 创建用户反馈
func (r *userFeedbackRepo) Create(ctx context.Context, feedback *biz.UserFeedback) error {
	return r.data.db.WithContext(ctx).Create(feedback).Error
}

// Update 更新用户反馈
func (r *userFeedbackRepo) Update(ctx context.Context, feedback *biz.UserFeedback) error {
	return r.data.db.WithContext(ctx).Save(feedback).Error
}

// Delete 删除用户反馈
func (r *userFeedbackRepo) Delete(ctx context.Context, id uint64) error {
	return r.data.db.WithContext(ctx).Delete(&biz.UserFeedback{}, id).Error
}

// Get 获取用户反馈
func (r *userFeedbackRepo) Get(ctx context.Context, id uint64) (*biz.UserFeedback, error) {
	var feedback biz.UserFeedback
	err := r.data.db.WithContext(ctx).First(&feedback, id).Error
	if err != nil {
		return nil, err
	}
	return &feedback, nil
}

// List 获取用户反馈列表
func (r *userFeedbackRepo) List(ctx context.Context, userID string, page, pageSize int32) ([]*biz.UserFeedback, int64, error) {
	var feedbacks []*biz.UserFeedback
	var total int64

	db := r.data.db.WithContext(ctx).Model(&biz.UserFeedback{})
	if userID != "" {
		db = db.Where("user_id = ?", userID)
	}

	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = db.Offset((int(page - 1) * int(pageSize))).Limit(int(pageSize)).
		Order("created_at DESC").Find(&feedbacks).Error
	if err != nil {
		return nil, 0, err
	}

	return feedbacks, total, nil
}

// ListByStatus 根据状态获取用户反馈列表
func (r *userFeedbackRepo) ListByStatus(ctx context.Context, status string, page, pageSize int32) ([]*biz.UserFeedback, int64, error) {
	var feedbacks []*biz.UserFeedback
	var total int64

	db := r.data.db.WithContext(ctx).Model(&biz.UserFeedback{}).Where("status = ?", status)

	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = db.Offset((int(page - 1) * int(pageSize))).Limit(int(pageSize)).
		Order("created_at DESC").Find(&feedbacks).Error
	if err != nil {
		return nil, 0, err
	}

	return feedbacks, total, nil
}

// -- 执行下面sql根据渠道id查询roi。
// SELECT
// 	aho.sub_channel_id AS '渠道ID',
// 	SUM(price) AS '半流程api分发收益',
// 	SUM(income) AS '半流程api分发净收益'
// FROM
// 	api_halfpro_order AS aho
// 	LEFT JOIN api_halfpro_order_line AS ahol ON aho.id = ahol.order_id
// WHERE
// 	aho.channel_id = 2729
// 	AND aho.created_at BETWEEN '2025-09-17 00:00:00'
// 	AND '2025-09-18 00:00:00'
// 	AND ahol.check_status = 1
// 	AND ahol.push_status = 1 
// GROUP BY
// 	aho.sub_channel_id;
func (r *userFeedbackRepo) GetChannelRoi(ctx context.Context, param *biz.ChannelRoiParam) ([]*biz.ChannelRoi, error) {
	var channelRois []*biz.ChannelRoi
	err := r.data.db.WithContext(ctx).Raw(`
		SELECT
			aho.sub_channel_id AS '渠道ID',
			SUM(price) AS '半流程api分发收益',
			SUM(income) AS '半流程api分发净收益'
		FROM
			api_halfpro_order AS aho
			LEFT JOIN api_halfpro_order_line AS ahol ON aho.id = ahol.order_id
		WHERE
			aho.channel_id = ?
			AND aho.created_at BETWEEN ?
			AND ?
			AND ahol.check_status = 1
			AND ahol.push_status = 1 
		GROUP BY
			aho.sub_channel_id;
	`, param.ChannelID, param.StartTime, param.EndTime).Scan(&channelRois).Error
	if err != nil {
		return nil, err
	}
	return channelRois, nil
}

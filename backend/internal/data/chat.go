package data

import (
	"context"

	"github.com/example/aichat/backend/internal/biz"
	"github.com/example/aichat/backend/models/generator/model"
	"github.com/example/aichat/backend/models/generator/query"
	"github.com/go-kratos/kratos/v2/log"
)

type chatRepo struct {
	data DataRepo
	log  *log.Helper
}

// NewChatRepo .
func NewChatRepo(data DataRepo, logger log.Logger) biz.ChatRepo {
	return &chatRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// CreateSession implements biz.ChatRepo.
func (r *chatRepo) CreateSession(ctx context.Context, session *model.AIChat) (*model.AIChat, error) {
	q := query.Use(r.data.GetDB())
	err := q.AIChat.WithContext(ctx).Create(session)
	if err != nil {
		return nil, err
	}
	return session, nil
}

// DeleteSession implements biz.ChatRepo.
func (r *chatRepo) DeleteSession(ctx context.Context, id int64) error {
	q := query.Use(r.data.GetDB())
	_, err := q.AIChat.WithContext(ctx).Where(q.AIChat.ID.Eq(id)).Delete()
	return err
}

// GetSession implements biz.ChatRepo.
func (r *chatRepo) GetSession(ctx context.Context, id int64) (*model.AIChat, error) {
	q := query.Use(r.data.GetDB())
	return q.AIChat.WithContext(ctx).Where(q.AIChat.ID.Eq(id)).First()
}

// ListSessions implements biz.ChatRepo.
func (r *chatRepo) ListSessions(ctx context.Context, userId int64, page int, pageSize int) ([]*model.AIChat, int64, error) {
	q := query.Use(r.data.GetDB())
	offset := (page - 1) * pageSize
	var chats []*model.AIChat
	var count int64
	var err error
	if page > 0 && pageSize > 0 {
		chats, count, err = q.AIChat.WithContext(ctx).Where(q.AIChat.CreatedBy.Eq(userId)).Order(q.AIChat.UpdatedAt.Desc()).FindByPage(offset, pageSize)
		if err != nil {
			return nil, 0, err
		}
	}else {
		chats, err = q.AIChat.WithContext(ctx).Where(q.AIChat.CreatedBy.Eq(userId)).Order(q.AIChat.UpdatedAt.Desc()).Find()
		if err != nil {
			return nil, 0, err
		}
		count = int64(len(chats))
	}
	return chats, count, nil
}

// UpdateSession implements biz.ChatRepo.
func (r *chatRepo) UpdateSession(ctx context.Context, session *model.AIChat) error {
	q := query.Use(r.data.GetDB())
	return q.AIChat.WithContext(ctx).Save(session)
}

// CreateMessage implements biz.ChatRepo.
func (r *chatRepo) CreateMessage(ctx context.Context, message *model.AIChatMessage) (*model.AIChatMessage, error) {
	q := query.Use(r.data.GetDB())
	err := q.AIChatMessage.WithContext(ctx).Create(message)
	if err != nil {
		return nil, err
	}
	return message, nil
}

// DeleteMessagesBySessionID implements biz.ChatRepo.
func (r *chatRepo) DeleteMessagesBySessionID(ctx context.Context, sessionId int64) error {
	q := query.Use(r.data.GetDB())
	_, err := q.AIChatMessage.WithContext(ctx).Where(q.AIChatMessage.SessionID.Eq(sessionId)).Delete()
	return err
}

// ListMessages implements biz.ChatRepo.
func (r *chatRepo) ListMessages(ctx context.Context, sessionId int64) ([]*model.AIChatMessage, error) {
	q := query.Use(r.data.GetDB())
	return q.AIChatMessage.WithContext(ctx).Where(q.AIChatMessage.SessionID.Eq(sessionId)).Order(q.AIChatMessage.CreatedAt.Asc()).Find()
}

package biz

import (
	"context"

	"github.com/example/aichat/backend/models/generator/model"
	"github.com/go-kratos/kratos/v2/log"
)

type ChatRepo interface {
	// Session
	CreateSession(ctx context.Context, session *model.AIChat) (*model.AIChat, error)
	GetSession(ctx context.Context, id int64) (*model.AIChat, error)
	ListSessions(ctx context.Context, userId int64, page, pageSize int) ([]*model.AIChat, int64, error)
	UpdateSession(ctx context.Context, session *model.AIChat) error
	DeleteSession(ctx context.Context, id int64) error

	// Message
	CreateMessage(ctx context.Context, message *model.AIChatMessage) (*model.AIChatMessage, error)
	ListMessages(ctx context.Context, sessionId int64) ([]*model.AIChatMessage, error)
	DeleteMessagesBySessionID(ctx context.Context, sessionId int64) error
}

type ChatUsecase struct {
	repo ChatRepo
	log  *log.Helper
}

func NewChatUsecase(repo ChatRepo, logger log.Logger) *ChatUsecase {
	return &ChatUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *ChatUsecase) CreateSession(ctx context.Context, userId int64, title string) (*model.AIChat, error) {
	t := &model.AIChat{
		UserID: userId,
		Title:  title,
	}
	t.New()
	return uc.repo.CreateSession(ctx, t)
}

func (uc *ChatUsecase) ListSessions(ctx context.Context, userId int64, page, pageSize int) ([]*model.AIChat, int64, error) {
	return uc.repo.ListSessions(ctx, userId, page, pageSize)
}

func (uc *ChatUsecase) GetSessionMessages(ctx context.Context, sessionId int64) ([]*model.AIChatMessage, error) {
	return uc.repo.ListMessages(ctx, sessionId)
}

func (uc *ChatUsecase) SaveMessage(ctx context.Context, msg *model.AIChatMessage) (*model.AIChatMessage, error) {
	return uc.repo.CreateMessage(ctx, msg)
}

func (uc *ChatUsecase) RenameSession(ctx context.Context, id int64, title string) error {
	session, err := uc.repo.GetSession(ctx, id)
	if err != nil {
		return err
	}
	session.Title = title
	return uc.repo.UpdateSession(ctx, session)
}

func (uc *ChatUsecase) DeleteSession(ctx context.Context, id int64) error {
	// Transaction could be better here, but for now simple deletion
	if err := uc.repo.DeleteMessagesBySessionID(ctx, id); err != nil {
		return err
	}
	return uc.repo.DeleteSession(ctx, id)
}

// Helper to create a new message struct
func (uc *ChatUsecase) NewMessage(sessionId int64, role model.RoleType, content string) *model.AIChatMessage {
	return &model.AIChatMessage{
		SessionID: sessionId,
		Role:      role,
		Content:   content,
		// CreatedAt: time.Now(), // GORM handles this
	}
}

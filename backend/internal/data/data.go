package data

import (
	"github.com/example/aichat/backend/internal/biz/base/loginprovider"
	"github.com/example/aichat/backend/internal/data/ai"
	"github.com/example/aichat/backend/internal/db"
	"github.com/example/aichat/backend/pkg/auth"

	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewSysUserRepo,
	NewSysMenuRepo,
	auth.NewAuthRepo,
	db.NewData,
	NewTransaction,
	NewDictTypeRepo,
	NewDictDataRepo,
	NewTrackerRepo,
	NewChatRepo,
	db.NewRedisRepo,
	loginprovider.NewStateCache,
	ai.NewAIAgentRepo,
	ai.NewAIApplicationRepo,
	ai.NewAIWorkflowRepo,
	ai.NewAIModelRepo,
	ai.NewAIPromptRepo,
	ai.NewAIToolRepo,
)

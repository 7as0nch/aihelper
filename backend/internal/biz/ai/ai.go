/*
 * @Author: chengjiang
 * @Date: 2025-12-09
 * @Description: AI 业务逻辑层
 */
package ai

import (
	"context"
	"fmt"

	pkgai "github.com/example/aichat/backend/pkg/ai"
	"github.com/example/aichat/backend/pkg/ai/adapter"
	"go.uber.org/zap"
)

// AIUsecase AI 业务逻辑
type AIUsecase struct {
	factory *pkgai.DefaultFactory
	log     *zap.Logger
}

// NewAIUsecase 创建 AI Usecase
func NewAIUsecase(log *zap.Logger) *AIUsecase {
	factory := pkgai.NewFactory()

	// 注册 Eino 适配器
	factory.RegisterAdapter(pkgai.AdapterTypeEino, adapter.NewEinoAdapter)

	// 注册 DeepADK 适配器
	factory.RegisterAdapter(pkgai.AdapterTypeDeepAdk, adapter.NewDeepAdkAdapter)

	return &AIUsecase{
		factory: factory,
		// repo:    repo, // TODO cache
		log:     log,
	}
}

// GetAgent 获取 Agent（按需加载，带缓存）
func (uc *AIUsecase) GetAgent(ctx context.Context) (pkgai.Agent, error) {
	agent, err := uc.factory.CreateWithSubAgents(ctx, &pkgai.AgentConfig{
		Name:        "globalAgent",
		Description: `你是一个全能助手: 
		1.可获取系统配置，如当前系统时间。
		2.websearch搜索用户想要搜索的东西并解答。`, // English description
		AdapterType: pkgai.AdapterTypeDeepAdk,
		ModelConfig: pkgai.ModelConfig{
			ModelType: "deepseek",
			ModelName: "deepseek",
			APIKey:    "sk-dc49fec5d27a416f8758ece703aed2ff",
			BaseURL:   "https://api.deepseek.com",
			Thinking:  true,
		},
		WithWebSearchAgent: true,
		MaxIteration: 10,
	}, []*pkgai.AgentConfig{})
	if err != nil {
		uc.log.Error("Failed to get agent", zap.Error(err))
		return nil, fmt.Errorf("failed to get agent: %w", err)
	}
	return agent, nil
}

// StreamChat 执行流式对话
func (uc *AIUsecase) Stream(ctx context.Context, req pkgai.Request) (<-chan pkgai.Response, error) {
	agent, err := uc.GetAgent(ctx)
	if err != nil {
		return nil, err
	}
	return agent.Stream(ctx, req)
}

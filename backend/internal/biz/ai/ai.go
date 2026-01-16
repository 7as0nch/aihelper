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
	"github.com/example/aichat/backend/pkg/ai/chatmodel"
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
	factory.RegisterAdapter(pkgai.AdapterTypeHost, adapter.NewHost)

	return &AIUsecase{
		factory: factory,
		// repo:    repo, // TODO cache
		log:     log,
	}
}

// GetAgent 获取 Agent（按需加载，带缓存）
func (uc *AIUsecase) GetAgent(ctx context.Context) (pkgai.Agent, error) {
	// TODO 不同用户可能选择不同的 Agent，也可能直接@某个智能体。
	agent, err := uc.factory.CreateWithSubAgents(ctx, &pkgai.AgentConfig{
		Name:        "主要智慧帮手",
		Description: `你是一个智慧助手: 
		## 你的功能：
		1.可获取系统配置，如当前系统时间。
		2.搜索用户想要搜索的东西并解答。
		3.业务相关。
		`, // English description
		AdapterType: pkgai.AdapterTypeDeepAdk,
		ModelConfig: pkgai.ModelConfig{
			ModelType: "deepseek",
			ModelName: "deepseek",
			APIKey:    "sk-dc49fec5d27a416f8758ece703aed2ff",
			BaseURL:   "https://api.deepseek.com",
			// ModelType: string(chatmodel.OPENAI_MODEL),
			// ModelName: "mimo-v2-flash",
			// APIKey:    "sk-c5pp6zngjk8gc2jxo0knffvw04g4zofvrrkjas8p5p7j51ml",
			// BaseURL:   "https://api.xiaomimimo.com/v1",
			Thinking:  true,
			TopP:      0.9,
		},
		WithWebSearchAgent: true,
		MaxIteration: 10,
	}, []*pkgai.AgentConfig{{
		Name:        "业务能力",
		Description: `你是一个业务助手: 报表相关需要返回表格形式。
		1.用户报表查询
		2.channel ROI查询。`, // English description
		AdapterType: pkgai.AdapterTypeEino,
		ModelConfig: pkgai.ModelConfig{
			ModelType: string(chatmodel.ARK_MODEL),
			ModelName: "doubao-seed-1-6-250615",
			APIKey:    "9292fee0-90fb-4739-86ff-eb1e886e2823",
			// ModelType: string(chatmodel.OPENAI_MODEL),
			// ModelName: "mimo-v2-flash",
			// APIKey:    "sk-c5pp6zngjk8gc2jxo0knffvw04g4zofvrrkjas8p5p7j51ml",
			// BaseURL:   "https://api.xiaomimimo.com/v1",
			Thinking:  true,
			TopP:      0.9,
		},
		WithWebSearchAgent: false,
		MaxIteration: 10,
	}})
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

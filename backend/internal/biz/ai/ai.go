/*
 * @Author: chengjiang
 * @Date: 2025-12-09
 * @Description: AI 业务逻辑层
 */
package ai

import (
	"context"
	"fmt"

	"github.com/example/aichat/backend/models"
	"github.com/example/aichat/backend/models/generator/model"
	pkgai "github.com/example/aichat/backend/pkg/ai"
	"github.com/example/aichat/backend/pkg/ai/adapter"
	"go.uber.org/zap"
)

// AIUsecase AI 业务逻辑
type AIUsecase struct {
	factory *pkgai.DefaultFactory
	appUC   *AIApplicationUseCase
	agentUC *AIAgentUseCase
	wfUC    *AIWorkflowUseCase
	modelUC *AIModelUseCase
	log     *zap.Logger
}

// NewAIUsecase 创建 AI Usecase
func NewAIUsecase(appUC *AIApplicationUseCase, agentUC *AIAgentUseCase, wfUC *AIWorkflowUseCase, modelUC *AIModelUseCase, log *zap.Logger) *AIUsecase {
	factory := pkgai.NewFactory()

	// 注册适配器
	factory.RegisterAdapter(pkgai.AdapterTypeEino, adapter.NewEinoAdapter)
	factory.RegisterAdapter(pkgai.AdapterTypeDeepAdk, adapter.NewDeepAdkAdapter)
	factory.RegisterAdapter(pkgai.AdapterTypeHost, adapter.NewHost)
	factory.RegisterAdapter(pkgai.AdapterTypeGraph, adapter.NewGraph)

	return &AIUsecase{
		factory: factory,
		appUC:   appUC,
		agentUC: agentUC,
		wfUC:    wfUC,
		modelUC: modelUC,
		log:     log,
	}
}

// GetAgent 获取 Agent（按需加载，带缓存）
func (uc *AIUsecase) GetAgent(ctx context.Context) (pkgai.Agent, error) {
	// 1. 获取应用配置 (实际应根据请求中的 appId 或 code 获取)
	app, err := uc.appUC.GetByCode(ctx, "default_app") // 暂时写死 code
	if err != nil {
		return nil, fmt.Errorf("failed to get application: %w", err)
	}

	// 2. 检查应用状态
	if app.Status != models.Status_Enabled {
		return nil, fmt.Errorf("application is disabled")
	}

	// 3. 根据应用类型构建 AgentConfig
	var agentConfig *pkgai.AgentConfig
	if app.Type == model.ProgramType_Custom {
		// 自定义类型：直接使用 SelfAgent
		if app.SelfAgent == nil {
			return nil, fmt.Errorf("self agent is not configured")
		}
		agentConfig, err = uc.buildAgentConfigFromSelfAgent(ctx, app.SelfAgent)
		if err != nil {
			return nil, fmt.Errorf("failed to build agent config from self agent: %w", err)
		}
	} else {
		// 预定义类型：需要通过 bind 关系获取（TODO: 实现 bind 关系查询）
		// 目前暂时返回错误，提示需要实现 bind 关系
		return nil, fmt.Errorf("predefined type requires agent bind relationship, not implemented yet")
	}

	// 4. 创建 Agent
	// 注意：此处如果需要子 Agent，还需要根据业务逻辑获取并构建
	agent, err := uc.factory.Create(ctx, agentConfig)
	if err != nil {
		uc.log.Error("Failed to get agent", zap.Error(err))
		return nil, fmt.Errorf("failed to get agent: %w", err)
	}
	return agent, nil
}

func (uc *AIUsecase) buildAgentConfig(ctx context.Context, m *model.AIAgent) (*pkgai.AgentConfig, error) {
	aiModel, err := uc.modelUC.GetByID(ctx, m.AIModelID)
	if err != nil {
		return nil, err
	}

	return &pkgai.AgentConfig{
		Name:               m.Name,
		Description:        m.Description,
		AdapterType:        uc.toPkgAdapterType(m.AdapterType),
		MaxIteration:       m.MaxIteration,
		WithWebSearchAgent: m.WithWebSearchAgent,
		WithWriteTODOs:     m.WithWriteTODOs,
		ModelConfig: pkgai.ModelConfig{
			ModelType: string(aiModel.ModelType),
			ModelName: aiModel.ModelName,
			APIKey:    aiModel.APIKey,
			BaseURL:   aiModel.BaseURL,
			Thinking:  true,
			TopP:      0.9,
		},
	}, nil
}

// buildAgentConfigFromSelfAgent 从 SelfAgent 构建 AgentConfig
// SelfAgent 中可能已经包含了 AIModel 的副本（通过 AIModel 字段）
func (uc *AIUsecase) buildAgentConfigFromSelfAgent(ctx context.Context, m *model.AIAgent) (*pkgai.AgentConfig, error) {
	var modelConfig pkgai.ModelConfig

	// 优先使用 SelfAgent 中的 AIModel 副本
	if m.AIModel != nil {
		// 使用 SelfAgent 中存储的 AIModel 副本
		aiModel := m.AIModel
		modelConfig = pkgai.ModelConfig{
			ModelType:   string(aiModel.ModelType),
			ModelName:   aiModel.ModelName,
			APIKey:      aiModel.APIKey,
			BaseURL:     aiModel.BaseURL,
			MaxTokens:   aiModel.MaxTokens,
			Temperature: aiModel.Temperature,
			TopP:        aiModel.TopP,
			Thinking:    true,
		}
	} else if m.AIModelID > 0 {
		// 如果没有 AIModel 副本，通过 AIModelID 获取
		aiModel, err := uc.modelUC.GetByID(ctx, m.AIModelID)
		if err != nil {
			return nil, fmt.Errorf("failed to get model by id %d: %w", m.AIModelID, err)
		}
		modelConfig = pkgai.ModelConfig{
			ModelType:   string(aiModel.ModelType),
			ModelName:   aiModel.ModelName,
			APIKey:      aiModel.APIKey,
			BaseURL:     aiModel.BaseURL,
			MaxTokens:   aiModel.MaxTokens,
			Temperature: aiModel.Temperature,
			TopP:        aiModel.TopP,
			Thinking:    true,
		}
	} else {
		return nil, fmt.Errorf("no model configuration found in self agent")
	}

	return &pkgai.AgentConfig{
		Name:               m.Name,
		Description:        m.Description,
		AdapterType:        uc.toPkgAdapterType(m.AdapterType),
		MaxIteration:       m.MaxIteration,
		WithWebSearchAgent: m.WithWebSearchAgent,
		WithWriteTODOs:     m.WithWriteTODOs,
		ModelConfig:        modelConfig,
	}, nil
}

func (uc *AIUsecase) buildWorkflowConfig(ctx context.Context, wf *model.AIWorkflow) (*pkgai.AgentConfig, error) {
	if wf.Definition == nil {
		return nil, fmt.Errorf("workflow definition is empty")
	}

	// 工作流通常也需要一个基础模型，这里可以从配置中获取或使用默认
	return &pkgai.AgentConfig{
		Name:           wf.Name,
		Description:    wf.Description,
		AdapterType:    pkgai.AdapterTypeGraph,
		WorkflowConfig: wf.Definition,
	}, nil
}

func (uc *AIUsecase) toPkgAdapterType(t model.AdapterType) pkgai.AdapterType {
	switch t {
	case model.AdapterType_ADK:
		return pkgai.AdapterTypeEino
	case model.AdapterType_DeepADK:
		return pkgai.AdapterTypeDeepAdk
	default:
		return pkgai.AdapterTypeEino
	}
}

// StreamChat 执行流式对话
func (uc *AIUsecase) Stream(ctx context.Context, req pkgai.Request) (<-chan pkgai.Response, error) {
	agent, err := uc.GetAgent(ctx)
	if err != nil {
		return nil, err
	}
	return agent.Stream(ctx, req)
}

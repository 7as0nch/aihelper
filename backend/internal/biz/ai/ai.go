/*
 * @Author: chengjiang
 * @Date: 2025-12-09
 * @Description: AI 业务逻辑层
 */
package ai

import (
	"context"
	"fmt"

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

	// 2. 根据用户角色选择配置 (此处逻辑需结合 Auth 模块，暂时模拟为 User)
	role := "user" // admin, user, guest
	var appConfig *pkgai.AppConfig
	switch role {
	case "admin":
		appConfig = app.AdminConfig
	case "user":
		appConfig = app.UserConfig
	default:
		appConfig = app.GuestConfig
	}

	if appConfig == nil {
		return nil, fmt.Errorf("no application config found for role: %s", role)
	}

	// 3. 构建 AgentConfig
	var agentConfig *pkgai.AgentConfig
	if appConfig.ActiveType == pkgai.ActiveTypeWorkflow {
		// 工作流模式
		wf, err := uc.wfUC.GetByCode(ctx, appConfig.TargetCode)
		if err != nil {
			return nil, err
		}
		agentConfig, err = uc.buildWorkflowConfig(ctx, wf)
		if err != nil {
			return nil, err
		}
	} else {
		// 普通 Agent 模式
		agentModel, err := uc.agentUC.GetByCode(ctx, appConfig.TargetCode)
		if err != nil {
			return nil, err
		}
		agentConfig, err = uc.buildAgentConfig(ctx, agentModel)
		if err != nil {
			return nil, err
		}
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

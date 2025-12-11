/* *
 * @Author: chengjiang
 * @Date: 2025-12-11
 * @Description: AI 管理服务层
**/
package ai

import (
	pb "github.com/example/aichat/backend/api/ai"
	bizai "github.com/example/aichat/backend/internal/biz/ai"
)

// AIService AI 管理服务
type AIService struct {
	pb.UnimplementedAIServer
	agentUC  *bizai.AIAgentUseCase
	modelUC  *bizai.AIModelUseCase
	promptUC *bizai.AIPromptUseCase
	toolUC   *bizai.AIToolUseCase
}

// NewAIService 创建 AI 服务
func NewAIService(
	agentUC *bizai.AIAgentUseCase,
	modelUC *bizai.AIModelUseCase,
	promptUC *bizai.AIPromptUseCase,
	toolUC *bizai.AIToolUseCase,
) *AIService {
	return &AIService{
		agentUC:  agentUC,
		modelUC:  modelUC,
		promptUC: promptUC,
		toolUC:   toolUC,
	}
}

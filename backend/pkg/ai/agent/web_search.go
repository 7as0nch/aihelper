/*
 * Copyright 2025 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package agent

import (
	"context"

	"github.com/cloudwego/eino-ext/components/tool/duckduckgo/v2"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/example/aichat/backend/pkg/ai"
	"github.com/example/aichat/backend/pkg/ai/chatmodel"
)

func NewWebSearchAgent(ctx context.Context, config *ai.AgentConfig) (adk.Agent, error) {
	cm, err := chatmodel.NewModel(ctx, chatmodel.ModelConfig{
		ModelType: chatmodel.ModelType(config.ModelConfig.ModelType),
		ModelName: config.ModelConfig.ModelName,
		ApiKey:    config.ModelConfig.APIKey,
		BaseURL:   config.ModelConfig.BaseURL,
		Thinking:  config.ModelConfig.Thinking,
	},
		chatmodel.WithMaxTokens(config.ModelConfig.MaxTokens),
		chatmodel.WithTemperature(config.ModelConfig.Temperature),
		chatmodel.WithTopP(config.ModelConfig.TopP),
		chatmodel.WithDisableThinking(!config.ModelConfig.Thinking),
	)
	if err != nil {
		return nil, err
	}

	searchTool, err := duckduckgo.NewTextSearchTool(ctx, &duckduckgo.Config{
		ToolName:   "duckduckgo_search",
		ToolDesc:   "search for information by duckduckgo",
		Region:     duckduckgo.RegionWT, // The geographical region for results.
		MaxResults: 3,                   // Limit the number of results returned.
	})
	if err != nil {
		return nil, err
	}

	return adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "WebSearchAgent",
		Description: "WebSearchAgent 支持联网搜索",
		Model:       cm,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: []tool.BaseTool{searchTool},
			},
		},
		MaxIterations: 10,
	})
}

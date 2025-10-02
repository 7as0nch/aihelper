package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aichat/backend/models"
	"github.com/aichat/backend/pkg/agent"
	"github.com/aichat/backend/pkg/db"
	"github.com/gin-gonic/gin"
)

// ChatController 聊天控制器
type ChatController struct {
	agentEngine *agent.Engine
}

// NewChatController 创建新的聊天控制器
func NewChatController() *ChatController {
	return &ChatController{
		agentEngine: agent.NewEngine(),
	}
}

// SendMessage 发送消息
func (cc *ChatController) SendMessage(c *gin.Context) {
	// TODO: 从JWT中获取用户ID
	// userID := getUserIDFromToken(c)
	userID := uint(1) // 临时硬编码

	var messageData struct {
		ChatID    uint                 `json:"chat_id"`
		Content   string               `json:"content" binding:"required"`
		Type      string               `json:"type" default:"text"`
		Metadata  map[string]interface{} `json:"metadata"`
	}

	if err := c.ShouldBindJSON(&messageData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var chat models.Chat
	if messageData.ChatID == 0 {
		// 创建新的聊天会话
		chat = models.Chat{
			UserID:  userID,
			Title:   generateChatTitle(messageData.Content),
			Model:   "gpt-3.5-turbo",
			Status:  1,
		}

		result := db.GetDB().Create(&chat)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create chat session"})
			return
		}
	} else {
		// 查找现有聊天会话
		result := db.GetDB().Where("id = ? AND user_id = ?", messageData.ChatID, userID).First(&chat)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Chat session not found"})
			return
		}
	}

	// 保存用户消息
	userMessage := models.Message{
		ChatID:    chat.ID,
		UserID:    userID,
		Content:   messageData.Content,
		Type:      messageData.Type,
		From:      models.MessageFromUser,
		Status:    1,
		Metadata:  messageData.Metadata,
	}

	result := db.GetDB().Create(&userMessage)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save message"})
		return
	}

	// 处理AI回复
	assistantReply, err := cc.processAssistantReply(c, &chat, &userMessage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 保存AI回复消息
	assistantMessage := models.Message{
		ChatID:   chat.ID,
		UserID:   userID,
		Content:  assistantReply,
		Type:     "text",
		From:     models.MessageFromAssistant,
		Status:   1,
	}

	result = db.GetDB().Create(&assistantMessage)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save assistant reply"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Message sent successfully",
		"chat_id": chat.ID,
		"messages": []interface{}{
			gin.H{
				"id":        userMessage.ID,
				"content":   userMessage.Content,
				"type":      userMessage.Type,
				"from":      userMessage.From,
				"created_at": userMessage.CreatedAt,
			},
			gin.H{
				"id":        assistantMessage.ID,
				"content":   assistantMessage.Content,
				"type":      assistantMessage.Type,
				"from":      assistantMessage.From,
				"created_at": assistantMessage.CreatedAt,
			},
		},
	})
}

// GetChatHistory 获取聊天历史
func (cc *ChatController) GetChatHistory(c *gin.Context) {
	// TODO: 从JWT中获取用户ID
	// userID := getUserIDFromToken(c)
	userID := uint(1) // 临时硬编码

	chatID := c.Query("chat_id")
	if chatID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "chat_id is required"})
		return
	}

	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "20")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page size"})
		return
	}
	// TODO: 实现分页查询聊天历史
	var messages []models.Message
	result := db.GetDB().
		Where("chat_id = ? AND user_id = ?", chatID, userID).
		Order("created_at asc").
		Find(&messages).Limit(pageSizeInt).Offset((pageInt - 1) * pageSizeInt)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chat history"})
		return
	}

	// 格式化返回数据
	formattedMessages := make([]interface{}, 0, len(messages))
	for _, msg := range messages {
		formattedMessages = append(formattedMessages, gin.H{
			"id":        msg.ID,
			"content":   msg.Content,
			"type":      msg.Type,
			"from":      msg.From,
			"created_at": msg.CreatedAt,
			"metadata":  msg.Metadata,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"messages": formattedMessages,
		"total":   len(messages),
	})
}

// 生成聊天标题
func generateChatTitle(content string) string {
	if len(content) > 30 {
		return content[:30] + "..."
	}
	return content
}

// 处理AI回复
func (cc *ChatController) processAssistantReply(c *gin.Context, chat *models.Chat, userMessage *models.Message) (string, error) {
	// 获取历史消息，用于上下文
	var historyMessages []models.Message
	result := db.GetDB().
		Where("chat_id = ?", chat.ID).
		Order("created_at desc").
		Limit(5).
		Find(&historyMessages)

	if result.Error != nil {
		return "", result.Error
	}

	// 检查是否需要调用函数工具或工作流
	// 这里简化处理，实际项目中应该根据用户输入和上下文决定是否调用工具

	// 如果用户输入包含特定关键词，尝试调用函数工具
	if containsKeyword(userMessage.Content, []string{"查询", "获取", "天气", "新闻"}) {
		// 示例：调用MCP服务中的函数工具
		executionCtx := &agent.ExecutionContext{
			Context:   c.Request.Context(),
			UserID:    userMessage.UserID,
			Variables: map[string]interface{}{"query": userMessage.Content},
		}

		// 这里使用一个示例的函数工具ID（实际应该根据业务逻辑选择）
		functionResult := cc.agentEngine.ExecuteFunction(executionCtx, 1)
		if functionResult.Success {
			return formatFunctionResult(functionResult.Output), nil
		}
	}

	// 如果用户输入包含特定关键词，尝试调用工作流
	if containsKeyword(userMessage.Content, []string{"流程", "步骤", "报告"}) {
		// 示例：调用工作流
		executionCtx := &agent.ExecutionContext{
			Context:   c.Request.Context(),
			UserID:    userMessage.UserID,
			Variables: map[string]interface{}{"query": userMessage.Content},
		}

		// 这里使用一个示例的工作流ID（实际应该根据业务逻辑选择）
		workflowResult := cc.agentEngine.ExecuteWorkflow(executionCtx, 1)
		if workflowResult.Success {
			return formatFunctionResult(workflowResult.Output), nil
		}
	}

	// 简单的AI回复模拟
	return "这是一个示例回复。在实际应用中，这里会调用真正的AI服务生成回复内容。", nil
}

// 检查内容是否包含关键词
func containsKeyword(content string, keywords []string) bool {
	// TODO: 实现更复杂的关键词匹配逻辑
	return false
}

// 格式化函数执行结果
func formatFunctionResult(result interface{}) string {
	// 尝试将结果转换为JSON字符串
	jsonResult, err := json.Marshal(result)
	if err != nil {
		return "获取到以下结果：" + string(jsonResult)
	}

	// 如果转换失败，返回字符串表示
	return "获取到以下结果：" + string(jsonResult)
}
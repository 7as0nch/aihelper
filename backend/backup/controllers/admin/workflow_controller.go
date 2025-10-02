package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/aichat/backend/models"
	"github.com/aichat/backend/pkg/db"
	"net/http"
	"strconv"
)

// WorkflowController 工作流控制器
type WorkflowController struct {}

// NewWorkflowController 创建新的工作流控制器
func NewWorkflowController() *WorkflowController {
	return &WorkflowController{}
}

// GetWorkflowList 获取工作流列表
func (wc *WorkflowController) GetWorkflowList(c *gin.Context) {
	// TODO: 从JWT中获取用户ID和权限
	// userID := getUserIDFromToken(c)

	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "20")
	keyword := c.Query("keyword")
	enabled := c.Query("enabled")

	// 转换分页参数
	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)
	offset := (pageInt - 1) * pageSizeInt

	// 构建查询
	query := db.GetDB().Model(&models.Workflow{})

	// 关键词筛选
	if keyword != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 启用状态筛选
	if enabled != "" {
		enabledBool, _ := strconv.ParseBool(enabled)
		query = query.Where("enabled = ?", enabledBool)
	}

	// 查询总数
	var total int64
	query.Count(&total)

	// 查询列表
	var workflows []models.Workflow
	query.Order("created_at desc").Offset(offset).Limit(pageSizeInt).Find(&workflows)

	// 格式化返回数据
	formattedWorkflows := make([]interface{}, 0, len(workflows))
	for _, workflow := range workflows {
		formattedWorkflows = append(formattedWorkflows, gin.H{
			"id":          workflow.ID,
			"name":        workflow.Name,
			"description": workflow.Description,
			"enabled":     workflow.Enabled,
			"created_at":  workflow.CreatedAt,
			"updated_at":  workflow.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    formattedWorkflows,
		"total":   total,
		"page":    pageInt,
		"page_size": pageSizeInt,
	})
}

// CreateWorkflow 创建工作流
func (wc *WorkflowController) CreateWorkflow(c *gin.Context) {
	var workflowData struct {
		Name        string             `json:"name" binding:"required"`
		Description string             `json:"description" binding:"required"`
		Enabled     bool               `json:"enabled" default:"true"`
		Steps       []models.WorkflowStep `json:"steps" binding:"required,dive"`
	}

	if err := c.ShouldBindJSON(&workflowData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查名称是否重复
	var existingWorkflow models.Workflow
	result := db.GetDB().Where("name = ?", workflowData.Name).First(&existingWorkflow)
	if result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Workflow name already exists"})
		return
	}

	// 检查步骤中的函数工具是否存在
	for _, step := range workflowData.Steps {
		var functionTool models.FunctionTool
		result := db.GetDB().First(&functionTool, step.FunctionToolID)
		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Function tool not found: " + strconv.FormatUint(uint64(step.FunctionToolID), 10)})
			return
		}
	}

	// 创建工作流
	dbConn := db.GetDB().Begin()
	workflow := models.Workflow{
		Name:        workflowData.Name,
		Description: workflowData.Description,
		Enabled:     workflowData.Enabled,
	}

	result = dbConn.Create(&workflow)
	if result.Error != nil {
		dbConn.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create workflow"})
		return
	}

	// 创建工作流步骤
	for i, step := range workflowData.Steps {
		workflowStep := models.WorkflowStep{
			WorkflowID:     workflow.ID,
			FunctionToolID: step.FunctionToolID,
			StepName:       step.StepName,
			Order:          i + 1, // 从1开始编号
			Condition:      step.Condition,
			VariableMapping: step.VariableMapping,
		}

		result := dbConn.Create(&workflowStep)
		if result.Error != nil {
			dbConn.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create workflow step"})
			return
		}
	}

	// 提交事务
	dbConn.Commit()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Workflow created successfully",
		"data": gin.H{
			"id":          workflow.ID,
			"name":        workflow.Name,
			"description": workflow.Description,
			"enabled":     workflow.Enabled,
			"created_at":  workflow.CreatedAt,
		},
	})
}

// UpdateWorkflow 更新工作流
func (wc *WorkflowController) UpdateWorkflow(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Workflow ID is required"})
		return
	}

	// 转换ID
	workflowID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workflow ID"})
		return
	}

	var workflowData struct {
		Name        string             `json:"name" binding:"required"`
		Description string             `json:"description" binding:"required"`
		Enabled     bool               `json:"enabled"`
		Steps       []models.WorkflowStep `json:"steps" binding:"required,dive"`
	}

	if err := c.ShouldBindJSON(&workflowData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查询工作流
	var workflow models.Workflow
	result := db.GetDB().First(&workflow, workflowID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workflow not found"})
		return
	}

	// 检查名称是否重复（排除当前记录）
	if workflowData.Name != workflow.Name {
		var existingWorkflow models.Workflow
		result := db.GetDB().Where("name = ? AND id != ?", workflowData.Name, workflowID).First(&existingWorkflow)
		if result.Error == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Workflow name already exists"})
			return
		}
	}

	// 检查步骤中的函数工具是否存在
	for _, step := range workflowData.Steps {
		var functionTool models.FunctionTool
		result := db.GetDB().First(&functionTool, step.FunctionToolID)
		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Function tool not found: " + strconv.FormatUint(uint64(step.FunctionToolID), 10)})
			return
		}
	}

	// 更新工作流
	dbConn := db.GetDB().Begin()

	// 更新工作流基本信息
	workflow.Name = workflowData.Name
	workflow.Description = workflowData.Description
	workflow.Enabled = workflowData.Enabled

	result = dbConn.Save(&workflow)
	if result.Error != nil {
		dbConn.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update workflow"})
		return
	}

	// 删除旧的工作流步骤
	result = dbConn.Where("workflow_id = ?", workflowID).Delete(&models.WorkflowStep{})
	if result.Error != nil {
		dbConn.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete old workflow steps"})
		return
	}

	// 创建新的工作流步骤
	for i, step := range workflowData.Steps {
		workflowStep := models.WorkflowStep{
			WorkflowID:     workflow.ID,
			FunctionToolID: step.FunctionToolID,
			StepName:       step.StepName,
			Order:          i + 1, // 从1开始编号
			Condition:      step.Condition,
			VariableMapping: step.VariableMapping,
		}

		result := dbConn.Create(&workflowStep)
		if result.Error != nil {
			dbConn.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create workflow step"})
			return
		}
	}

	// 提交事务
	dbConn.Commit()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Workflow updated successfully",
		"data": gin.H{
			"id":          workflow.ID,
			"name":        workflow.Name,
			"description": workflow.Description,
			"enabled":     workflow.Enabled,
			"updated_at":  workflow.UpdatedAt,
		},
	})
}

// DeleteWorkflow 删除工作流
func (wc *WorkflowController) DeleteWorkflow(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Workflow ID is required"})
		return
	}

	// 转换ID
	workflowID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workflow ID"})
		return
	}

	// 查询工作流
	var workflow models.Workflow
	result := db.GetDB().First(&workflow, workflowID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workflow not found"})
		return
	}

	// 软删除工作流（级联删除步骤）
	dbConn := db.GetDB().Begin()

	// 删除工作流步骤
	result = dbConn.Where("workflow_id = ?", workflowID).Delete(&models.WorkflowStep{})
	if result.Error != nil {
		dbConn.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete workflow steps"})
		return
	}

	// 删除工作流
	result = dbConn.Delete(&workflow)
	if result.Error != nil {
		dbConn.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete workflow"})
		return
	}

	// 提交事务
	dbConn.Commit()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Workflow deleted successfully",
	})
}
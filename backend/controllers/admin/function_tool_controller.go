package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/aichat/backend/models"
	"github.com/aichat/backend/pkg/db"
	"net/http"
	"strconv"
)

// FunctionToolController 函数工具控制器
type FunctionToolController struct {}

// NewFunctionToolController 创建新的函数工具控制器
func NewFunctionToolController() *FunctionToolController {
	return &FunctionToolController{}
}

// GetFunctionList 获取函数工具列表
func (ftc *FunctionToolController) GetFunctionList(c *gin.Context) {
	// TODO: 从JWT中获取用户ID和权限
	// userID := getUserIDFromToken(c)

	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "20")
	keyword := c.Query("keyword")
	functionType := c.Query("type")
	enabled := c.Query("enabled")

	// 转换分页参数
	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)
	offset := (pageInt - 1) * pageSizeInt

	// 构建查询
	query := db.GetDB().Model(&models.FunctionTool{})

	// 关键词筛选
	if keyword != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 类型筛选
	if functionType != "" {
		query = query.Where("type = ?", functionType)
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
	var functionTools []models.FunctionTool
	query.Order("created_at desc").Offset(offset).Limit(pageSizeInt).Find(&functionTools)

	// 格式化返回数据
	formattedTools := make([]interface{}, 0, len(functionTools))
	for _, tool := range functionTools {
		formattedTools = append(formattedTools, gin.H{
			"id":          tool.ID,
			"name":        tool.Name,
			"description": tool.Description,
			"type":        tool.Type,
			"enabled":     tool.Enabled,
			"created_at":  tool.CreatedAt,
			"updated_at":  tool.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    formattedTools,
		"total":   total,
		"page":    pageInt,
		"page_size": pageSizeInt,
	})
}

// CreateFunction 创建函数工具
func (ftc *FunctionToolController) CreateFunction(c *gin.Context) {
	var functionTool models.FunctionTool

	if err := c.ShouldBindJSON(&functionTool); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查名称是否重复
	var existingTool models.FunctionTool
	result := db.GetDB().Where("name = ?", functionTool.Name).First(&existingTool)
	if result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Function tool name already exists"})
		return
	}

	// 创建函数工具
	dbConn := db.GetDB()
	result = dbConn.Create(&functionTool)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create function tool"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Function tool created successfully",
		"data": gin.H{
			"id":          functionTool.ID,
			"name":        functionTool.Name,
			"description": functionTool.Description,
			"type":        functionTool.Type,
			"enabled":     functionTool.Enabled,
			"created_at":  functionTool.CreatedAt,
		},
	})
}

// UpdateFunction 更新函数工具
func (ftc *FunctionToolController) UpdateFunction(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Function tool ID is required"})
		return
	}

	// 转换ID
	toolID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid function tool ID"})
		return
	}

	// 查询函数工具
	dbConn := db.GetDB()
	var functionTool models.FunctionTool
	result := dbConn.First(&functionTool, toolID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Function tool not found"})
		return
	}

	// 绑定更新数据
	var updateData models.FunctionTool
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查名称是否重复（排除当前记录）
	if updateData.Name != functionTool.Name {
		var existingTool models.FunctionTool
		result := dbConn.Where("name = ? AND id != ?", updateData.Name, toolID).First(&existingTool)
		if result.Error == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Function tool name already exists"})
			return
		}
	}

	// 更新函数工具
	functionTool.Name = updateData.Name
	functionTool.Description = updateData.Description
	functionTool.Type = updateData.Type
	functionTool.Config = updateData.Config
	functionTool.Parameters = updateData.Parameters
	functionTool.Enabled = updateData.Enabled

	result = dbConn.Save(&functionTool)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update function tool"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Function tool updated successfully",
		"data": gin.H{
			"id":          functionTool.ID,
			"name":        functionTool.Name,
			"description": functionTool.Description,
			"type":        functionTool.Type,
			"enabled":     functionTool.Enabled,
			"updated_at":  functionTool.UpdatedAt,
		},
	})
}

// DeleteFunction 删除函数工具
func (ftc *FunctionToolController) DeleteFunction(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Function tool ID is required"})
		return
	}

	// 转换ID
	toolID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid function tool ID"})
		return
	}

	// 查询函数工具
	dbConn := db.GetDB()
	var functionTool models.FunctionTool
	result := dbConn.First(&functionTool, toolID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Function tool not found"})
		return
	}

	// 软删除函数工具
	result = dbConn.Delete(&functionTool)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete function tool"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Function tool deleted successfully",
	})
}
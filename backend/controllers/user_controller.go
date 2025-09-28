package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/aichat/backend/models"
	"github.com/aichat/backend/pkg/auth"
	"github.com/aichat/backend/pkg/db"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

// UserController 用户控制器
type UserController struct{}

// NewUserController 创建新的用户控制器
func NewUserController() *UserController {
	return &UserController{}
}

// Login 用户登录
func (uc *UserController) Login(c *gin.Context) {
	var loginData struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查询用户
	var user models.User
	dbConn := db.GetDB()
	result := dbConn.Where("username = ?", loginData.Username).First(&user)
	// 如果用户名不存在，尝试用邮箱登录
	if result.Error != nil {
		result = dbConn.Where("email = ?", loginData.Username).First(&user)
		if result.Error != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}
	}

	// 检查密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// 检查用户状态
	if user.Status == 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "User account is disabled"})
		return
	}

	// 更新最后登录时间
	dbConn.Model(&user).Update("last_login_at", time.Now())

	// 生成JWT令牌
	token, err := auth.GenerateToken(uint64(user.ID), user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login successful",
		"data": gin.H{
			"user_id":   user.ID,
			"username":  user.Username,
			"email":     user.Email,
			"role":      user.Role,
			"created_at": user.CreatedAt,
		},
		"token": token,
	})
}

// Register 用户注册
func (uc *UserController) Register(c *gin.Context) {
	var registerData struct {
		Username string `json:"username" binding:"required,min=3,max=50"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6,max=30"`
	}

	if err := c.ShouldBindJSON(&registerData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查用户名是否已存在
	var existingUser models.User
	result := db.GetDB().Where("username = ?", registerData.Username).First(&existingUser)
	if result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	// 检查邮箱是否已存在
	result = db.GetDB().Where("email = ?", registerData.Email).First(&existingUser)
	if result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	// 加密密码
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(registerData.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	// 创建新用户
	user := models.User{
		Username: registerData.Username,
		Email:    registerData.Email,
		Password: string(passwordHash),
		Role:     models.RoleUser,
		Status:   1,
	}

	result = db.GetDB().Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User registered successfully",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}

// GetUserInfo 获取用户信息
func (uc *UserController) GetUserInfo(c *gin.Context) {
	// 从上下文中获取用户ID
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, ok := userIDInterface.(uint64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user ID"})
		return
	}

	var user models.User
	dbConn := db.GetDB()
	result := dbConn.First(&user, userID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 不返回密码等敏感信息
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"user_id":   user.ID,
			"username":  user.Username,
			"email":     user.Email,
			"role":      user.Role,
			"created_at": user.CreatedAt,
			"last_login_at": user.LastLoginAt,
		},
	})
}

// TODO: 实现JWT生成和验证相关函数
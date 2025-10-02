package main

import (
	"github.com/aichat/backend/controllers"
	"github.com/aichat/backend/controllers/admin"
	"github.com/aichat/backend/pkg/auth"
	"github.com/aichat/backend/pkg/db"
	// "github.com/aichat/backend/pkg/initdata"
	"github.com/aichat/backend/pkg/redis"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

func main() {
	// 初始化配置
	initConfig()

	// 初始化JWT
	auth.InitJWT()

	// 初始化数据库
	// initDB()

	// 初始化Redis
	// initRedis()

	// 初始化管理员账户
	// initdata.InitAdminUser()

	// 创建路由
	router := gin.Default()

	// 初始化控制器
	userController := controllers.NewUserController()
	chatController := controllers.NewChatController()
	functionToolController := admin.NewFunctionToolController()
	workflowController := admin.NewWorkflowController()

	// 设置CORS
	router.Use(corsMiddleware())

	// 健康检查
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// API路由分组
	api := router.Group("/api/v1")
	{
		// 用户相关路由
		user := api.Group("/user")
		{
			user.POST("/login", userController.Login)
			user.POST("/register", userController.Register)
			user.GET("/info", auth.AuthMiddleware(), userController.GetUserInfo)
		}

		// 聊天相关路由
		chat := api.Group("/chat")
		{
			chat.POST("/send", auth.AuthMiddleware(), chatController.SendMessage)
			chat.GET("/history", auth.AuthMiddleware(), chatController.GetChatHistory)
		}

		// 管理后台相关路由
		admin := api.Group("/admin")
		admin.Use(auth.AuthMiddleware(), auth.AdminMiddleware())
		{
			// 配置管理
			config := admin.Group("/config")
			{
				config.GET("/mcp", handleGetMCPConfig)
				config.POST("/mcp", handleUpdateMCPConfig)
			}

			// 函数工具管理
			function := admin.Group("/function")
			{
				function.GET("/list", functionToolController.GetFunctionList)
				function.POST("/create", functionToolController.CreateFunction)
				function.PUT("/:id", functionToolController.UpdateFunction)
				function.DELETE("/:id", functionToolController.DeleteFunction)
			}

			// 工作流管理
			workflow := admin.Group("/workflow")
			{
				workflow.GET("/list", workflowController.GetWorkflowList)
				workflow.POST("/create", workflowController.CreateWorkflow)
				workflow.PUT("/:id", workflowController.UpdateWorkflow)
				workflow.DELETE("/:id", workflowController.DeleteWorkflow)
			}
		}
	}

	// 启动服务器
	port := viper.GetString("server.port")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// 初始化配置
func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../config")
	viper.AddConfigPath("/etc/aichat/")

	// 设置默认值
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("database.driver", "mysql")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "3306")
	viper.SetDefault("database.username", "root")
	viper.SetDefault("database.password", "password")
	viper.SetDefault("database.dbname", "aichat")
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", "6379")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Config file not found, using defaults: %v", err)
	}
}

// 初始化数据库
func initDB() {
	db.InitDB()
}

// 初始化Redis
func initRedis() {
	redis.InitRedis()
}

// CORS中间件
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// 认证中间件 - 已在auth包中实现，此处为兼容性保留
func authMiddleware() gin.HandlerFunc {
	return auth.AuthMiddleware()
}

// 管理员中间件 - 已在auth包中实现，此处为兼容性保留
func adminMiddleware() gin.HandlerFunc {
	return auth.AdminMiddleware()
}

// MCP配置管理函数（待实现）
func handleGetMCPConfig(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Not implemented yet"})
}

func handleUpdateMCPConfig(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Not implemented yet"})
}
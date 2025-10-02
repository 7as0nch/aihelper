package initdata

import (
	"github.com/aichat/backend/models"
	"github.com/aichat/backend/pkg/db"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

// InitAdminUser 初始化管理员用户
func InitAdminUser() {
	db := db.GetDB()

	// 检查是否已存在管理员用户
	var count int64
	db.Model(&models.User{}).Where("role = ?", "admin").Count(&count)

	if count > 0 {
		log.Println("Admin user already exists, skipping initialization")
		return
	}

	// 加密默认密码
	password, err := bcrypt.GenerateFromPassword([]byte("admin123456"), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		return
	}

	// 创建管理员用户
	now := time.Now()
	adminUser := models.User{
		Username:    "admin",
		Password:    string(password),
		Email:       "admin@example.com",
		Role:        "admin",
		Status:      1, // 1表示启用
		CreatedAt:   now,
		UpdatedAt:   now,
		LastLoginAt: &now,
	}

	result := db.Create(&adminUser)
	if result.Error != nil {
		log.Printf("Failed to create admin user: %v", err)
		return
	}

	log.Println("Admin user created successfully")
	log.Println("Username: admin")
	log.Println("Password: admin123456")
	log.Println("Please change the default password after first login")
}
package db

import (
	"fmt"
	"github.com/aichat/backend/models"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

var DB *gorm.DB

// 初始化数据库连接
func InitDB() {
	// 从配置文件获取数据库连接信息
	host := viper.GetString("database.host")
	port := viper.GetString("database.port")
	username := viper.GetString("database.username")
	password := viper.GetString("database.password")
	dbname := viper.GetString("database.dbname")
	charset := viper.GetString("database.charset")
	parseTime := viper.GetBool("database.parseTime")
	loc := viper.GetString("database.loc")

	// 构建DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%v&loc=%s",
		username, password, host, port, dbname, charset, parseTime, loc)

	// 配置日志级别
	logLevel := logger.Silent
	if viper.GetString("log.level") == "debug" {
		logLevel = logger.Info
	}

	// 连接数据库
	var err error
	DB, err = gorm.Open(mysql.New(mysql.Config{
		DSN: dsn,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 获取通用数据库对象 sql.DB，用于设置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Database connection established successfully")

	// 自动迁移表结构
	autoMigrate()
}

// 自动迁移表结构
func autoMigrate() {
	// 自动迁移所有模型
	modelsList := []interface{}{}
	modelsList = append(modelsList, &models.User{})
	modelsList = append(modelsList, &models.Chat{})
	modelsList = append(modelsList, &models.Message{})
	modelsList = append(modelsList, &models.FunctionTool{})
	modelsList = append(modelsList, &models.Workflow{})
	modelsList = append(modelsList, &models.WorkflowStep{})

	for _, model := range modelsList {
		if err := DB.AutoMigrate(model); err != nil {
			log.Printf("Failed to auto migrate model %T: %v", model, err)
		}
	}

	log.Println("Database tables auto migrated successfully")
}

// 获取数据库连接
func GetDB() *gorm.DB {
	if DB == nil {
		InitDB()
	}
	return DB
}
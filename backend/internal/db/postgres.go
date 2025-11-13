package db

import (
	"fmt"
	"time"

	"github.com/example/aichat/backend/internal/conf"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// GormZapWriter 实现 gorm 日志接口的 zap writer 适配器
type GormZapWriter struct {
	log *zap.Logger
}

// Printf 实现 gorm logger.Writer 接口
func (w GormZapWriter) Printf(format string, args ...interface{}) {
	// 使用 Info 级别记录 SQL 日志
	w.log.Info(fmt.Sprintf(format, args...))
}

// NewPostgresDB 创建 PostgreSQL 数据库连接，支持 SQL 日志打印
func NewPostgresDB(conf *conf.Bootstrap, log *zap.Logger) (*gorm.DB, error) {
	pgConfig := conf.Data.PgDatabase
	// 构建 DSN (Data Source Name)
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		pgConfig.Host,
		pgConfig.Port,
		pgConfig.User,
		pgConfig.Password,
		pgConfig.Dbname,
		pgConfig.Sslmode,
	)

	// 自定义 GORM 日志配置，支持 SQL 日志打印
	newLogger := logger.New(
		GormZapWriter{log: log},
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  1,    // 日志级别，控制 SQL 日志的详细程度
			IgnoreRecordNotFoundError: true,        // 忽略记录未找到错误
			Colorful:                  true,        // 彩色打印，增强可读性
		},
	)

	// 配置 GORM
	gormConfig := &gorm.Config{
		Logger: newLogger, // 启用 SQL 日志
	}

	// 连接数据库
	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	// 配置连接池参数
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// 设置连接池参数，优化数据库性能
	sqlDB.SetMaxIdleConns(10)           // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)          // 最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大生命周期

	// 记录连接信息
	log.Info("PostgreSQL connection established",
		zap.String("host", pgConfig.Host),
		zap.Int("port", int(pgConfig.Port)),
		zap.String("dbname", pgConfig.Dbname),
	)

	return db, nil
}

// MustNewPostgresDB 创建 PostgreSQL 数据库连接，如果出错则 panic
// 适用于应用启动时的数据库初始化，确保数据库连接成功
func MustNewPostgresDB(conf *conf.Bootstrap, log *zap.Logger) *gorm.DB {
	db, err := NewPostgresDB(conf, log)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to PostgreSQL: %v", err))
	}
	return db
}
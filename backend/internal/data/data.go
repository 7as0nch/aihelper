package data

import (
	"github.com/aichat/backend/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewUserRepo)

// Data .
type Data struct {
	db *gorm.DB
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	// 连接数据库
	db, err := gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{})
	if err != nil {
		log.NewHelper(logger).Errorf("failed to connect database: %v", err)
		return nil, nil, err
	}

	// 获取通用数据库对象 sql.DB ，然后使用其提供的功能
	sqlDB, err := db.DB()
	if err != nil {
		log.NewHelper(logger).Errorf("failed to get sql.DB: %v", err)
		return nil, nil, err
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	cleanup := func() {
		if sqlDB != nil {
			sqlDB.Close()
		}
		log.NewHelper(logger).Info("closing the data resources")
	}
	
	return &Data{db: db}, cleanup, nil
}

// DB 返回数据库连接
func (d *Data) DB() *gorm.DB {
	return d.db
}

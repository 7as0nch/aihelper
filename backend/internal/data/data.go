package data

import (
	"github.com/example/aichat/backend/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo, NewUserFeedbackRepo)

// Data .
type Data struct {
	// TODO wrapped database client
	db *gorm.DB
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	// 初始化mysqldb
	db := NewMysqlDB(c.Database.Source)
	if db == nil {
		log.NewHelper(logger).Error("NewMysqlDB failed")
		return nil, nil, nil
	}
	return &Data{}, cleanup, nil
}

func NewMysqlDB(source string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(source), &gorm.Config{

	})
	if err != nil {
		log.Error(err)
		panic(err)
	}
	return db
}

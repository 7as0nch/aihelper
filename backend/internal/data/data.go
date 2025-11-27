package data

import (
	"github.com/example/aichat/backend/internal/conf"
	"github.com/example/aichat/backend/internal/db"
	"github.com/example/aichat/backend/pkg/auth"
	"go.uber.org/zap"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewSysUserRepo,
	NewSysMenuRepo,
	auth.NewAuthRepo,
	NewData, NewUserFeedbackRepo,
	NewDictTypeRepo,
	NewDictDataRepo,
	NewTrackerRepo)

type DataRepo interface {
	GetDB() *gorm.DB
	Redis()
	// TODO OSS
}

// Data .
type Data struct {
	// TODO wrapped database client
	db *gorm.DB
}

// GetDB implements DataRepo.
func (d *Data) GetDB() *gorm.DB {
	return d.db
}

// Redis implements DataRepo.
func (d *Data) Redis() {
	panic("unimplemented")
}

// NewData .
func NewData(c *conf.Bootstrap, logger *zap.Logger) (DataRepo, func(), error) {
	// 初始化mysqldb
	db := db.MustNewPostgresDB(c, logger)
	if db == nil {
		logger.Error("NewMysqlDB failed")
		return nil, nil, nil
	}
	cleanup := func() {
		logger.Info("closing the data resources")
	}
	return &Data{
		db: db,
	}, cleanup, nil
}

func NewMysqlDB(source string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(source), &gorm.Config{})
	if err != nil {
		log.Error(err)
		panic(err)
	}
	return db
}

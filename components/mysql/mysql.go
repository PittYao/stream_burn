package mysql

import (
	"github.com/PittYao/stream_burn/components/config"
	"github.com/PittYao/stream_burn/components/log"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var Instance *gorm.DB

func Load() error {
	if config.C.Mysql.Dsn == "" {
		panic("配置文件Mysql.dsn为空,不能启动服务")
	}

	// 配置zap为gorm的日志
	logger := log.New(log.L)
	db, err := gorm.Open(mysql.Open(config.C.Mysql.Dsn), &gorm.Config{
		Logger: logger,
	})
	if err != nil {
		log.L.Error("MySQL连接异常", zap.Error(err))
		panic("MySQL启动异常: " + err.Error())
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.L.Error("MySQL连接异常", zap.Error(err))
		panic("MySQL连接异常: " + err.Error())
	}
	sqlDB.SetMaxIdleConns(config.C.Mysql.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.C.Mysql.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(config.C.Mysql.ConnMaxLifetime)

	log.L.Info("MySQL连接初始化成功")

	Instance = db
	return nil
}

func WithoutConfig(dsn string) {
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic("MySQL启动异常: " + err.Error())
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic("MySQL连接异常: " + err.Error())
	}
	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetConnMaxLifetime(1 * time.Hour)

	Instance = db

}

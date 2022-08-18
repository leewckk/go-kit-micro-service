package gorm

import (
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	DateFmt     = "2006-01-02"
	DateTimeFmt = "2006-01-02 15:04:05"
)

func NewOrmDB(conn string, driver string, debug bool) *gorm.DB {

	logrus.Info("try to connect ", driver, " conn: ", conn)
	db, err := gorm.Open(mysql.Open(conn), &gorm.Config{
		QueryFields: true,
	})
	if nil != err {
		panic(err)
	}

	/// 连接池配置
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(256)
	sqlDB.SetConnMaxLifetime(60 * time.Second)

	if debug {
		db = db.Debug()
	}
	return db
}

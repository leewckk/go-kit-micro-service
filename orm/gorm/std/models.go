package std

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

// var (
// 	__db__ *gorm.DB
// )

// func init() {
// 	__db__ = createDB()
// }

// func createDB() *gorm.DB {

// 	conf := configure.GetDefault().Config()
// 	conn := conf.Db.Conn
// 	driver := conf.Db.Driver
// 	debug := conf.Db.Debug

// 	if "" == conn {
// 		panic(fmt.Sprintf("Failed found database conn from config file ,config path: %v", configure.GetDefault().Uri))
// 	}
// 	return NewDB(conn, driver, debug)
// }

func NewDB(conn string, driver string, debug bool) *gorm.DB {

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

// func GetDB() *gorm.DB {
// 	sqlDB, err := __db__.DB()
// 	if nil != err {
// 		__db__ = createDB()
// 	}
// 	if err := sqlDB.Ping(); err != nil {
// 		sqlDB.Close()
// 		__db__ = createDB()
// 	}
// 	return __db__
// }

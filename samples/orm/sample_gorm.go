package main

import (
	"time"

	std "github.com/leewckk/go-kit-micro-service/orm/gorm"
)

type User struct {
	Id       int64     `json:"-" gorm:"primaryKey;autoIncrement"`
	UserId   int64     `json:"userId"`
	UserName string    `json:"userName"`
	Weight   float64   `json:"weight"`
	Birth    time.Time `json:"birth" gorm:"type:date"`
}

func (u *User) TableName() string {
	return "user"
}

func main() {
	driver := "mysql"
	conn := `gokit-user:gokit-passwd@tcp(localhost)/gokit?charset=utf8&&loc=Local&&parseTime=true`
	debug := true

	db := std.NewOrmDB(conn, driver, debug)
	autoMigrates := []interface{}{
		new(User),
	}
	err := db.AutoMigrate(autoMigrates...)
	if nil != err {
		panic(err)
	}
}

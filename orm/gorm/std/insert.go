package std

import (
	"go-kit-micro-service/logs"

	"gorm.io/gorm/clause"
)

func InsertLock(list interface{}) (interface{}, error) {

	db := GetDB()
	tx := db.Begin()
	defer func() {
		if r := recover(); nil != r {
			tx.Rollback()
		}
	}()
	result := db.Clauses(
		clause.Locking{
			Strength: "UPDATE",
		},
		clause.OnConflict{
			UpdateAll: true,
		}).Create(list)

	if nil != result.Error {
		logs.Errorf("Error insert record, err : %v ", result.Error)
		tx.Rollback() //// 事务回滚
		return list, result.Error
	}

	/// 提交，释放锁
	err := tx.Commit().Error
	return list, err
}

func Insert(list interface{}) (interface{}, error) {

	db := GetDB()
	result := db.Clauses(
		clause.OnConflict{
			UpdateAll: true,
		}).Create(list)

	if nil != result.Error {
		logs.Errorf("Error insert record, err : %v ", result.Error)
		return list, result.Error
	}

	/// 提交，释放锁
	return list, result.Error
}

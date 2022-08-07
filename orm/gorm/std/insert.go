package std

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func InsertLock(db *gorm.DB, list interface{}) (interface{}, error) {

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
		logrus.Errorf("Error insert record, err : %v ", result.Error)
		tx.Rollback() //// 事务回滚
		return list, result.Error
	}

	/// 提交，释放锁
	err := tx.Commit().Error
	return list, err
}

func Insert(db *gorm.DB, list interface{}) (interface{}, error) {

	result := db.Clauses(
		clause.OnConflict{
			UpdateAll: true,
		}).Create(list)

	if nil != result.Error {
		logrus.Errorf("Error insert record, err : %v ", result.Error)
		return list, result.Error
	}

	/// 提交，释放锁
	return list, result.Error
}

package std

import (
	"reflect"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type QueryObject interface {
	NewModel() interface{}
	NewContainer() interface{}
}

func NewQs(db *gorm.DB, object QueryObject, filters map[string]interface{}) (*gorm.DB, error) {
	if len(filters) == 0 {
		return db, nil
	}

	cond, vals, err := WhereBuild(filters)
	if nil != err {
		logrus.Error(err)
		return nil, err
	}
	return db.Model(object.NewModel()).Where(cond, vals...), nil
}

func QueryModel(object QueryObject, filters map[string]interface{}, orderby []string, limit, offset int) (interface{}, error) {

	container := object.NewContainer()
	qs, err := NewQs(object, filters)
	if nil != err {
		return container, err
	}

	for _, order := range orderby {
		qs = qs.Order(order)
	}
	qs = qs.Limit(limit).Offset(offset)
	result := qs.Find(container)
	if nil != result.Error {
		logrus.Error(result.Error)
	}
	return container, result.Error
}

func QueryModelWithTotal(object QueryObject, filters map[string]interface{}, orderby []string, limit, offset int) (int64, interface{}, error) {

	var total int64
	container := object.NewContainer()
	qs, err := NewQs(object, filters)
	if nil != err {
		return 0, container, err
	}

	qs.Count(&total)
	for _, order := range orderby {
		qs = qs.Order(order)
	}
	qs = qs.Limit(limit).Offset(offset)
	result := qs.Find(container)
	if nil != result.Error {
		logrus.Error(result.Error)
	}
	return total, container, result.Error
}

func QueryModelFromChan(object QueryObject, filters map[string]interface{}, orderby []string, limit, offset int) chan interface{} {

	ch := make(chan interface{}, 0)
	go func() {
		defer func() {
			close(ch)
		}()

		limitmax := 1000
		qs, err := NewQs(object, filters)
		if nil != err {
			return
		}

		for _, order := range orderby {
			qs = qs.Order(order)
		}

		if limit != 0 && limit < limitmax {
			limitmax = limit
		}
		total := limit

		for {
			// fmt.Printf("total : %v, limitmax: %v, limit : %v, offset : %v \n", total, limitmax, limit, offset)
			container := object.NewContainer()

			qs = qs.Limit(limitmax).Offset(offset)
			result := qs.Find(container)
			if result.Error != nil {
				logrus.Errorf("Error query table : %v, err : %v", reflect.TypeOf(object.NewModel()), err)
			}
			ch <- container
			if result.RowsAffected < int64(limitmax) {
				break
			}
			offset += limitmax
			total -= limitmax

			if total <= 0 && limit > 0 {
				break
			}

			if limit > 0 && total < limitmax {
				limitmax = total
			}
		}
	}()
	return ch
}

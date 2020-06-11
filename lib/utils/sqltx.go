package utils

import (
	"strings"

	"github.com/google/uuid"

	"github.com/jinzhu/gorm"
	"github.com/juetun/dashboard-api-car/basic/utils/hashid"
)

type HashModel interface {
	SaltForHID() string
	GetID() int
	StartHidInit()
}

// 生成GUID
// @param prefix string 数据库名+表名
func Guid(prefix string) string {
	guid := uuid.New()
	s := guid.String()
	s = strings.Join(strings.Split(s, "-"), "")
	return s
}
func OnPanic(f func()) {
	if err := recover(); err != nil {
		f()
	}
}

func CreateForHID(db *gorm.DB, model HashModel, haveStartTransaction ...bool) (err error) {
	// 开启事物
	startTransaction := false
	if len(haveStartTransaction) > 0 {
		startTransaction = haveStartTransaction[0]
	}
	tx := db

	// 如果外部开启了事务,则无需执行本动作
	if !startTransaction {
		tx = db.Begin()
		if err = tx.Error; err != nil {
			return err
		}
		defer OnPanic(
			func() {
				tx.Rollback()
			},
		)
		defer func() {
			if err != nil {
				tx.Rollback()
			}
		}()
	}

	model.StartHidInit()
	// 创建
	if err = tx.Create(model).Error; err != nil {
		return err
	}

	// 根据id 设置hid
	hid, err := hashid.Encode(model.SaltForHID(), model.GetID())
	if err != nil {
		return err
	}

	// 更新设置hid后的model
	if err := tx.Model(model).Update("hid", hid).Error; err != nil {
		return err
	}

	// 如果外部开启了事务,则无需执行本动作
	if !startTransaction {
		err = tx.Commit().Error
	}
	return
}

// 创建并返回Hid
func CreateForHIDAndReHid(db *gorm.DB, model HashModel) (s string, err error) {
	// 开启事物

	tx := db.Begin()
	if err = tx.Error; err != nil {
		return "", err
	}

	defer OnPanic(
		func() {
			tx.Rollback()
		},
	)

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 创建
	if err = tx.Create(model).Error; err != nil {
		return "", err
	}

	// 根据id 设置hid
	hid, err := hashid.Encode(model.SaltForHID(), model.GetID())
	if err != nil {
		return "", err
	}

	// 更新设置hid后的model
	if err := tx.Model(model).Update("hid", hid).Error; err != nil {
		return "", err
	}
	err = tx.Commit().Error

	return hid, err
}

// 应用于事务中
func CreateForHIDTx(tx *gorm.DB, model HashModel) (err error) {
	if err = tx.Error; err != nil {
		return err
	}

	defer OnPanic(
		func() {
			tx.Rollback()
		},
	)

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 创建
	if err = tx.Create(model).Error; err != nil {
		return err
	}

	// 根据id 设置hid
	hid, err := hashid.Encode(model.SaltForHID(), model.GetID())
	if err != nil {
		return err
	}

	// 更新设置hid后的model
	if err := tx.Model(model).Update("hid", hid).Error; err != nil {
		return err
	}
	return
}

// // 通用事物操作
// func TxProcess(db *gorm.DB, txFuncGroup ...func(tx *gorm.DB) error) (err error) {
//	// 开启事物
//	//var err error
//
//	tx := db.Begin()
//	if err = tx.Error; err != nil {
//		return err
//	}
//
//	defer OnPanic(
//		func() {
//			tx.Rollback()
//		},
//	)
//
//	defer func() {
//		if err != nil {
//			tx.Rollback()
//		}
//	}()
//
//	for _, doTx := range txFuncGroup {
//		if err := doTx(tx); err != nil {
//			panic(err)
//		}
//	}
//	tx.Commit()
//	return
// }

package dao

import (
	"casinoDemo/api/casino/model"
	"fmt"

	"gorm.io/gorm"
)

type GlobalDao struct {
}

func NewGlobalDao() *GlobalDao {
	return &GlobalDao{}
}

// FilterRec 过滤单条记录
func (d GlobalDao) FilterRec(tx *gorm.DB, condEq, condNeq map[string]any, order []string) (*model.Global, error) {
	getRec := &model.Global{}
	ses := tx.Table(getRec.TableName())
	for k, v := range condEq {
		ses.Where(fmt.Sprintf("%s = ?", k), v)
	}
	for k, v := range condNeq {
		ses.Where(fmt.Sprintf("%s <> ?", k), v)
	}
	for _, v := range order {
		ses.Order(v)
	}
	err := ses.First(&getRec).Error
	if err != nil {
		return nil, err
	}
	return getRec, nil
}

// Updates 更新
func (d GlobalDao) Updates(tx *gorm.DB, data *model.Global, sel []any) (*model.Global, error) {
	ses := tx.Table(data.TableName())
	if len(sel) > 0 {
		ses.Select(sel[0], sel[1:]...)
	}
	err := ses.Updates(data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Create 创建
func (t GlobalDao) Create(tx *gorm.DB, data *model.Global) (*model.Global, error) {
	err := tx.Table(data.TableName()).Create(data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

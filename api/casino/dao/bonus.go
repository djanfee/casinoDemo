package dao

import (
	"casinoDemo/api/casino/model"
	"fmt"

	"gorm.io/gorm"
)

type BonusDao struct {
	model.Bonus
}

func NewBonusDao() *BonusDao {
	return &BonusDao{}
}

// FilterRec 过滤单条记录
func (d *BonusDao) FilterRec(tx *gorm.DB, condEq, condNeq map[string]any, order []string) (*model.Bonus, error) {
	getRec := &model.Bonus{}
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
func (d *BonusDao) Updates(tx *gorm.DB, data *model.Bonus, sel []any) (*model.Bonus, error) {
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
func (t *BonusDao) Create(tx *gorm.DB, data *model.Bonus) (*model.Bonus, error) {
	err := tx.Table(data.TableName()).Create(data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

// CreateBatch 批量创建
func (t *BonusDao) CreateBatch(tx *gorm.DB, data []*model.Bonus) ([]*model.Bonus, error) {
	err := tx.Table(model.DefaultBonus.TableName()).Create(data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}
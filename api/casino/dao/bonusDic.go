package dao

import (
	"casinoDemo/api/casino/model"
	"fmt"

	"gorm.io/gorm"
)

type BonusDicDao struct {
}

func NewBonusDicDao() *BonusDicDao {
	return &BonusDicDao{}
}

// FilterRec 过滤单条记录
func (d *BonusDicDao) FilterRec(tx *gorm.DB, condEq, condNeq map[string]any, order []string) (*model.BonusDic, error) {
	getRec := &model.BonusDic{}
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

// FilterRecs 过滤多条记录
func (d *BonusDicDao) FilterRecs(tx *gorm.DB, condEq, condNeq, condIn map[string]interface{}, dbOrder []string) ([]*model.BonusDic, error) {
	getRecs := make([]*model.BonusDic, 0)
	ses := tx.Table(model.DefaultBonusDic.TableName())
	for k, v := range condEq {
		ses.Where(fmt.Sprintf("%s = ?", k), v)
	}
	for k, v := range condNeq {
		ses.Where(fmt.Sprintf("%s <> ?", k), v)
	}
	for k, v := range condIn {
		ses.Where(fmt.Sprintf("%s IN ?", k), v)
	}
	parentCode, ok := condEq["parentCode"]
	if ok && fmt.Sprintf("%v", parentCode) != "top" {
		ses.Or(fmt.Sprintf("( platform_api_code = '%v' AND parentCode = '%v')", parentCode, "top"))
	}
	for _, v := range dbOrder {
		ses.Order(v)
	}
	err := ses.Find(&getRecs).Error
	if err != nil {
		return nil, err
	}
	return getRecs, nil
}

// Updates 更新
func (d *BonusDicDao) Updates(tx *gorm.DB, data *model.BonusDic, sel []any) (*model.BonusDic, error) {
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
func (t *BonusDicDao) Create(tx *gorm.DB, data *model.BonusDic) (*model.BonusDic, error) {
	err := tx.Table(data.TableName()).Create(data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Count 总数
func (d *BonusDicDao) Count(tx *gorm.DB, condEq, condNeq, condIn map[string]any) (int64, error) {
	total := int64(0)
	ses := tx.Table(model.DefaultBonusDic.TableName())
	for k, v := range condEq {
		ses.Where(fmt.Sprintf("%v = ?", k), v)
	}
	for k, v := range condNeq {
		ses.Where(fmt.Sprintf("%v <> ?", k), v)
	}
	for k, v := range condIn {
		ses.Where(fmt.Sprintf("%v IN ?", k), v)
	}
	err := ses.Count(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}

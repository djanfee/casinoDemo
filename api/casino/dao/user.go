package dao

import (
	"casinoDemo/api/casino/model"
	"fmt"

	"gorm.io/gorm"
)

type UserDao struct {
}

func NewUserDao() *UserDao {
	return &UserDao{}
}

// FilterRec 过滤单条记录
func (d *UserDao) FilterRec(tx *gorm.DB, condEq, condNeq map[string]any, order []string) (*model.User, error) {
	getRec := &model.User{}
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
func (d *UserDao) Updates(tx *gorm.DB, data *model.User, sel []any) (*model.User, error) {
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
func (t *UserDao) Create(tx *gorm.DB, data *model.User) (*model.User, error) {
	err := tx.Table(data.TableName()).Create(data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Count 总数
func (d *UserDao) Count(tx *gorm.DB, condEq, condNeq, condIn map[string]any) (int64, error) {
	total := int64(0)
	ses := tx.Table(model.DefaultUser.TableName())
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

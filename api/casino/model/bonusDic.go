package model

var DefaultBonusDic = &BonusDic{}

type BonusDic struct {
	ID      uint64 `gorm:"primaryKey;column:id;type:bigint unsigned;not null" json:"-"`  // 主键
	Index   int64  `gorm:"column:index;type:bigint;not null;default:0" json:"index"`     // 索引
	BonusID uint64 `gorm:"column:bonusId;type:bigint unsigned;not null" json:"bonus_id"` // 分红id
}

// TableName get sql table name.获取数据库表名
func (m *BonusDic) TableName() string {
	return "bonusDic"
}

package model

var DefaultBonus = &Bonus{}

type Bonus struct {
	ID           uint64 `gorm:"primaryKey;column:id;type:bigint unsigned;not null" json:"-"`          // 主键
	Bonus        int64  `gorm:"column:bonus;type:bigint;not null;default:0" json:"bonus"`             // 分红
	BeginRound   int64  `gorm:"column:beginRound;type:bigint;not null;default:0" json:"begin_round"`  // 起始期号
	ContainItems int64  `gorm:"column:containItems;type:int;not null;default:0" json:"contain_items"` // 共包含多少期
}

// TableName get sql table name.获取数据库表名
func (m *Bonus) TableName() string {
	return "bonus"
}

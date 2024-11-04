package model

var DefaultUser = &User{}

type User struct {
	ID            uint64 `gorm:"primaryKey;column:id;type:bigint unsigned;not null" json:"-"`               // 主键
	Address       string `gorm:"column:address;type:varchar(255);not null;default:''" json:"address"`       // 地址
	DepositAmount int64  `gorm:"column:depositAmount;type:bigint;not null;default:0" json:"deposit_amount"` // 质押金额
	ClaimedRound  int64  `gorm:"column:claimedRound;type:int;not null;default:0" json:"claimed_round"`      // 已领取轮次
}

// TableName get sql table name.获取数据库表名
func (m *User) TableName() string {
	return "user"
}

package model

var DefaultGlobal = &Global{}

type Global struct {
	ID                  uint64 `gorm:"primaryKey;column:id;type:bigint unsigned;not null" json:"-"`                             // 主键
	DepositAmount       int64  `gorm:"column:depositAmount;type:bigint;not null;default:0" json:"deposit_amount"`               // 总质押金额
	Round               int64  `gorm:"column:round;type:bigint;not null;default:0" json:"round"`                                // 轮次
	NextRoundIncrAmount int64  `gorm:"column:nextRoundIncrAmount;type:bigint;not null;default:0" json:"next_round_incr_amount"` // 下一轮次添加金额
	PresentIncome       int64  `gorm:"column:presentIncome;type:bigint;not null;default:0" json:"present_income"`               // 当前收益
	TotalIncome         int64  `gorm:"column:totalIncome;type:bigint;not null;default:0" json:"total_income"`                   // 总收益
	NextRoundDecrAmount int64  `gorm:"column:nextRoundDecrAmount;type:bigint;not null;default:0" json:"next_round_decr_amount"` // 下一轮次减少金额
}

// TableName get sql table name.获取数据库表名
func (m *Global) TableName() string {
	return "global"
}

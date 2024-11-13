package casino_svc

// GloablData 全局数据
type GloablData struct {
	DepositAmount  int64 // 总存款
	CompletedRound int64 // 已完结的轮次
	PresentIncome  int64 // 当前轮次总盈利
	TotalIncome    int64 // 总盈利
}

// UserData 用户数据
type UserData struct {
	Address       string // 地址
	DepositAmount int64  // 质押金额
	ClaimedRound  int64  // 已领取轮次
}

// Bonus 分红
type Bonus struct {
	BeginRound    int64   // 起始轮次
	BonusAmount   float64 // 分红
	ContainRounds int64   // 包含轮次
}

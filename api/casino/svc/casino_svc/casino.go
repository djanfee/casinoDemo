package casino_svc

import (
	"casinoDemo/api/casino/dao"
	"casinoDemo/api/casino/model"
	"math"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type CasinoSvc struct {
	CasinoDb    *gorm.DB
	UserDao     *dao.UserDao
	GlobalDao   *dao.GlobalDao
	BonusDao    *dao.BonusDao
	BonusDicDao *dao.BonusDicDao
}

func NewCasinoSvc(casinoDb *gorm.DB) *CasinoSvc {
	return &CasinoSvc{
		CasinoDb:    casinoDb,
		UserDao:     dao.NewUserDao(),
		GlobalDao:   dao.NewGlobalDao(),
		BonusDao:    dao.NewBonusDao(),
		BonusDicDao: dao.NewBonusDicDao(),
	}
}

func (s *CasinoSvc) CalculateBonus(blockSeq int64, global *model.Global) error {
	// 当前期号
	curRound := math.Ceil(float64(blockSeq) / float64(OneDayBlocks))

	// 未分红期号
	unBonusRound := curRound - float64(global.Round) - 1
	if unBonusRound <= 0 {
		return nil
	}

	// 计算分红
	bonusAmount := float64(global.PresentIncome) / float64(global.DepositAmount)

	// 存储分红
	creatingBonus := &model.Bonus{
		BeginRound:   global.Round + 1,
		Bonus:        int64(bonusAmount),
		ContainItems: int64(unBonusRound),
	}
	createdBonus, err := s.BonusDao.Create(s.CasinoDb, creatingBonus)
	if err != nil {
		logx.Error("create bonus error", logx.Field("error", err))
		return err
	}

	// 存储分红字典
	for i := int64(0); i < int64(unBonusRound); i++ {
		creatingBonusDic := &model.BonusDic{
			Index:   global.Round + i + 1,
			BonusID: createdBonus.ID,
		}
		_, err = s.BonusDicDao.Create(s.CasinoDb, creatingBonusDic)
		if err != nil {
			logx.Error("create bonusDic error", logx.Field("error", err))
			return err
		}
	}

	// 存储数据
	global.Round += int64(unBonusRound)
	global.DepositAmount += global.NextRoundAmount
	global.NextRoundAmount = 0
	_, err = s.GlobalDao.Updates(s.CasinoDb, global, nil)
	if err != nil {
		logx.Error("update global error", logx.Field("error", err))
		return err
	}

	return nil
}

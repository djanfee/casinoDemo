package logic

import (
	"casinoDemo/api/casino/internal/svc"
	"casinoDemo/api/casino/internal/types"
	"casinoDemo/api/casino/model"
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ClaimBonusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewClaimBonusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClaimBonusLogic {
	return &ClaimBonusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ClaimBonusLogic) ClaimBonus(req *types.ClaimBonusReq) (resp *types.ClaimBonusResp, err error) {
	// 查询全局数据
	global, err := l.svcCtx.CasinoSvc.GlobalDao.FilterRec(l.svcCtx.CasinoDb, nil, nil, nil)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logx.WithContext(l.ctx).Errorf("get global data error: %v", err)
		return nil, err
	}
	if global == nil {
		creatingGlobal := &model.Global{
			DepositAmount:       0,
			Round:               0,
			NextRoundIncrAmount: 0,
		}
		global, err = l.svcCtx.CasinoSvc.GlobalDao.Create(l.svcCtx.CasinoDb, creatingGlobal)
		if err != nil {
			logx.WithContext(l.ctx).Errorf("create global data error: %v", err)
			return nil, err
		}
	}

	// 查询用户数据
	user, err := l.svcCtx.CasinoSvc.UserDao.FilterRec(l.svcCtx.CasinoDb, map[string]any{"address": req.Address}, nil, nil)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logx.WithContext(l.ctx).Errorf("get user data error: %v", err)
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}
	unclaimedRounds := global.Round - user.ClaimedRound
	if unclaimedRounds < 1 {
		return nil, errors.New("no unclaimed rounds")
	}

	// 查询用户未领取的分红轮次
	claimingRounds := make([]int64, 0)
	for i := 0; i < int(unclaimedRounds); i++ {
		claimingRounds = append(claimingRounds, user.ClaimedRound+int64(i)+1)
	}

	// 查询用户分红数据
	bonusDic, err := l.svcCtx.CasinoSvc.BonusDicDao.FilterRecs(l.svcCtx.CasinoDb, nil, nil, map[string]interface{}{"index": claimingRounds}, nil)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("get bonusDic data error: %v", err)
		return nil, err
	}
	bounsIds := make([]uint64, 0)
	for _, v := range bonusDic {
		bounsIds = append(bounsIds, v.BonusID)
	}
	bonuses, err := l.svcCtx.CasinoSvc.BonusDao.FilterRecs(l.svcCtx.CasinoDb, nil, nil, map[string]interface{}{"id": bounsIds}, nil)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("get bonus data error: %v", err)
		return nil, err
	}
	userBounsAmount := int64(0)
	for _, v := range bonuses {
		userBounsAmount += v.Bonus * user.DepositAmount
	}

	// 更新用户数据

	return
}

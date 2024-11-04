package logic

import (
	"context"
	"errors"

	"casinoDemo/api/casino/internal/svc"
	"casinoDemo/api/casino/internal/types"
	"casinoDemo/api/casino/model"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type ParticipateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewParticipateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ParticipateLogic {
	return &ParticipateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ParticipateLogic) Participate(req *types.ParticipateReq) (resp *types.ParticipateResp, err error) {
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
		creatingUser := &model.User{
			Address:       req.Address,
			DepositAmount: req.Value,
			ClaimedRound:  global.Round,
		}
		user, err = l.svcCtx.CasinoSvc.UserDao.Create(l.svcCtx.CasinoDb, creatingUser)
		if err != nil {
			logx.WithContext(l.ctx).Errorf("create user data error: %v", err)
			return nil, err
		}
	}

	// 判断是否触发分红
	err = l.svcCtx.CasinoSvc.CalculateBonus(req.BlockSeq, global)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("calculate bonus error: %v", err)
		return nil, err
	}

	// 更新用户数据
	user.DepositAmount += req.Value

	// 更新全局数据
	global.NextRoundIncrAmount += req.Value

	// 存储用户数据
	_, err = l.svcCtx.CasinoSvc.UserDao.Updates(l.svcCtx.CasinoDb, user, nil)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("update user data error: %v", err)
		return nil, err
	}

	// 存储全局数据
	_, err = l.svcCtx.CasinoSvc.GlobalDao.Updates(l.svcCtx.CasinoDb, global, nil)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("update global data error: %v", err)
		return nil, err
	}

	return &types.ParticipateResp{}, nil
}

package logic

import (
	"context"
	"errors"

	"casinoDemo/api/casino/internal/svc"
	"casinoDemo/api/casino/internal/types"
	"casinoDemo/api/casino/svc/casino_svc"

	"github.com/zeromicro/go-zero/core/logx"
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
	// 加载全局数据
	err = l.svcCtx.CasinoSvc.LoadGlobalData(l.ctx)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("load global data error. err:%v", err)
		return nil, err
	}

	// 获取当前轮次和未分红轮次以及是否可以开启下一轮
	_, unBonusRound, canNextRound, err := l.svcCtx.CasinoSvc.GetCurrentRoundAndUnBonusRound(l.ctx, req.BlockSeq)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("get current round and un bonus round error: %v", err)
		return nil, err
	}

	// 是否开启下一轮
	if canNextRound {
		err = l.svcCtx.CasinoSvc.StartNextRound(l.ctx, unBonusRound)
		if err != nil {
			logx.WithContext(l.ctx).Errorf("start next round error: %v", err)
			return nil, err
		}
	}

	// 检查是否已在用户列表中
	existedUser, err := l.svcCtx.CasinoSvc.GetUserData(l.ctx, req.Address)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("get user data error. err:%v", err)
		return nil, err
	}
	if existedUser != nil {
		return nil, errors.New("user already exists")
	}

	// 从下一轮用户列表中查询用户数据并更新
	user, err := l.svcCtx.CasinoSvc.GetUserFromNextRoundUserList(l.ctx, req.Address)
	if err != nil {
		logx.WithContext(l.ctx).Error("get user data error. err:%v", err)
		return nil, err
	}
	if user == nil {
		user = &casino_svc.UserData{
			Address:       req.Address,
			DepositAmount: 0,
			ClaimedRound:  l.svcCtx.CasinoSvc.GlobalData.CompletedRound,
		}
	}
	user.ClaimedRound = l.svcCtx.CasinoSvc.GlobalData.CompletedRound
	user.DepositAmount += req.Value
	err = l.svcCtx.CasinoSvc.AddUserToNextRoundUserList(l.ctx, user)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("add user to next round user list error: %v", err)
		return nil, err
	}

	// 存储全局数据
	err = l.svcCtx.CasinoSvc.SaveGlobalData(l.ctx)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("save global data error: %v", err)
		return nil, err
	}

	return &types.ParticipateResp{}, nil
}

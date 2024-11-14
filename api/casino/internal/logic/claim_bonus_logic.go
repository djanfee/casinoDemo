package logic

import (
	"casinoDemo/api/casino/internal/svc"
	"casinoDemo/api/casino/internal/types"
	"casinoDemo/api/casino/svc/casino_svc"
	"context"
	"errors"

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
		logx.Info("==============can next round")
		err = l.svcCtx.CasinoSvc.StartNextRound(l.ctx, unBonusRound)
		if err != nil {
			logx.WithContext(l.ctx).Errorf("start next round error: %v", err)
			return nil, err
		}
	}

	// 检查用户是否存在
	existedUser, err := l.svcCtx.CasinoSvc.GetUserData(l.ctx, req.Address)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("get user error. err:%v", err)
		return nil, err
	}
	if existedUser == nil {
		return nil, errors.New("user not found")
	}

	// 将用户添加到提现用户列表
	err = l.svcCtx.CasinoSvc.AddUserToWithdrawUserList(l.ctx, &casino_svc.UserData{
		Address:       req.Address,
		DepositAmount: 0,
		ClaimedRound:  l.svcCtx.CasinoSvc.GlobalData.CompletedRound,
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("add user to withdraw user list error. err:%v", err)
		return nil, err
	}
	return &types.ClaimBonusResp{}, nil
}

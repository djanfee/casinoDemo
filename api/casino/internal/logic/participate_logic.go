package logic

import (
	"context"

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

	// 当前轮次
	remainder := req.BlockSeq % casino_svc.OneDayBlocks
	curRound := int64(0)
	if remainder == 0 {
		curRound = req.BlockSeq / casino_svc.OneDayBlocks
	} else {
		curRound = req.BlockSeq/casino_svc.OneDayBlocks + 1
	}

	// 未分红轮次
	canNextRound := false
	unBonusRound := curRound - l.svcCtx.CasinoSvc.GlobalData.CompletedRound - 1
	if unBonusRound >= 1 {
		canNextRound = true
	}

	// 开启下一轮
	if canNextRound {
		// 计算分红
		err = l.svcCtx.CasinoSvc.CalculateBonus(l.ctx, unBonusRound)
		if err != nil {
			logx.WithContext(l.ctx).Errorf("calculate bonus error: %v", err)
			return nil, err
		}

		// 处理新增质押
		err = l.svcCtx.CasinoSvc.HandleNewDeposit(l.ctx)
		if err != nil {
			logx.WithContext(l.ctx).Errorf("handle new deposit error: %v", err)
			return nil, err
		}

		// 处理取款
		err = l.svcCtx.CasinoSvc.HandleWithdraw(l.ctx)
		if err != nil {
			logx.WithContext(l.ctx).Errorf("handle withdraw error: %v", err)
			return nil, err
		}
	}

	// 存储全局数据
	err = l.svcCtx.CasinoSvc.SaveGlobalData(l.ctx)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("save global data error: %v", err)
		return nil, err
	}

	return &types.ParticipateResp{}, nil
}

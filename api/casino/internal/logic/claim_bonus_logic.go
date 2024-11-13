package logic

import (
	"casinoDemo/api/casino/internal/svc"
	"casinoDemo/api/casino/internal/types"
	"context"

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
	// 更新用户数据

	return
}

package handler

import (
	"net/http"

	"casinoDemo/api/casino/internal/logic"
	"casinoDemo/api/casino/internal/svc"
	"casinoDemo/api/casino/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ClaimBonusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ClaimBonusReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewClaimBonusLogic(r.Context(), svcCtx)
		resp, err := l.ClaimBonus(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

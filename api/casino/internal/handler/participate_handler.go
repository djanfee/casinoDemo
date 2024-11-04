package handler

import (
	"net/http"

	"casinoDemo/api/casino/internal/logic"
	"casinoDemo/api/casino/internal/svc"
	"casinoDemo/api/casino/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ParticipateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ParticipateReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewParticipateLogic(r.Context(), svcCtx)
		resp, err := l.Participate(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

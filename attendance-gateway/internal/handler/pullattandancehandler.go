package handler

import (
	"net/http"

	"attendance-gateway/internal/logic"
	"attendance-gateway/internal/svc"
	"attendance-gateway/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func PullAttandanceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PullAttRequest
		if err := httpx.ParseJsonBody(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewPullAttandanceLogic(r.Context(), svcCtx)
		resp, err := l.PullAttandance(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

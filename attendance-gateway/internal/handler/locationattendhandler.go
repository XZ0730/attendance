package handler

import (
	"net/http"

	"attendance-gateway/internal/logic"
	"attendance-gateway/internal/svc"
	"attendance-gateway/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func LocationAttendHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LocationAttReq
		if err := httpx.ParseJsonBody(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewLocationAttendLogic(r.Context(), svcCtx)
		resp, err := l.LocationAttend(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

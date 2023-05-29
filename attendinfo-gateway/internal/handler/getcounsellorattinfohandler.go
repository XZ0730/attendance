package handler

import (
	"net/http"

	"attendinfo-gateway/internal/logic"
	"attendinfo-gateway/internal/svc"
	"attendinfo-gateway/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetCounsellorAttInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CounsellorInfoRequest
		if err := httpx.ParseJsonBody(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewGetCounsellorAttInfoLogic(r.Context(), svcCtx)
		resp, err := l.GetCounsellorAttInfo(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

package handler

import (
	"appeal-gateway/internal/common/errorx"
	"appeal-gateway/internal/logic"
	"appeal-gateway/internal/svc"
	"appeal-gateway/internal/types"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetAppealListBySidHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AppealListRequest
		if err := httpx.ParseForm(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, errorx.NewCodeError(10002, err.Error()))
			return
		}
		l := logic.NewGetAppealListBySidLogic(r.Context(), svcCtx)
		resp, err := l.GetAppealListBySid(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, errorx.NewCodeError(int(resp.Status), resp.Message))
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

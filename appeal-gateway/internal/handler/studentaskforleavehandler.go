package handler

import (
	"net/http"

	"appeal-gateway/internal/common/errorx"
	"appeal-gateway/internal/logic"
	"appeal-gateway/internal/svc"
	"appeal-gateway/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func StudentAskforLeaveHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AppealRequest
		if err := httpx.ParseJsonBody(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, errorx.NewCodeError(10002, err.Error()))
			return
		}

		l := logic.NewStudentAskforLeaveLogic(r.Context(), svcCtx)
		resp, err := l.StudentAskforLeave(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, errorx.NewCodeError(int(resp.Status), resp.Message))
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

package handler

import (
	"net/http"

	"appeal-gateway/internal/logic"
	"appeal-gateway/internal/svc"
	"appeal-gateway/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetComplainTablesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ComplainGetRequest
		v := r.URL.Query()
		req.CounsellorID = v.Get("coun_id")
		req.StudentName = v.Get("stu_name")
		req.SupervisorName = v.Get("sup_name")
		req.Major = v.Get("major")
		req.College = v.Get("college")
		req.SupervisorID = v.Get("sup_id")
		l := logic.NewGetComplainTablesLogic(r.Context(), svcCtx)
		resp, err := l.GetComplainTables(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

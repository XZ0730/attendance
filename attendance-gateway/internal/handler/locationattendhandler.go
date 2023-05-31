package handler

import (
	"attend/common/errorx"
	"encoding/json"
	"fmt"
	"net/http"

	"attendance-gateway/internal/logic"
	"attendance-gateway/internal/model"
	"attendance-gateway/internal/svc"
	"attendance-gateway/internal/types"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func LocationAttendHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LocationAttReq
		if err := httpx.ParseJsonBody(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, errorx.NewCodeError(11111, err.Error()))
			return
		}

		l := logic.NewLocationAttendLogic(r.Context(), svcCtx)
		resp, err := l.LocationAttend(&req)
		fmt.Println("ok1")
		if resp.Status == 200 {
			supervisor_client, ok := model.Manager.Clients.Load(req.SupervisorID)
			fmt.Println("ok1")
			fmt.Println("sid1:", req.SupervisorID)
			if ok {
				locReply := &model.LocationReply{
					StudentID: req.StudentID,
					Message:   "定位签到",
				}
				fmt.Println("ok2")
				locr, _ := json.Marshal(locReply)
				supervisor_client.(*model.Client).Socket.WriteMessage(websocket.TextMessage, locr)
				fmt.Println("ok3")
				// supervisor_client.(*model.Client).Send<-
			}
		}

		if err != nil {
			httpx.ErrorCtx(r.Context(), w, errorx.NewCodeError(int(resp.Status), err.Error()))
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

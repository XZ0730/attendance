package handler

import (
	"attendance-gateway/internal/errorx"
	"attendance-gateway/internal/logic"
	"attendance-gateway/internal/model"
	"attendance-gateway/internal/svc"
	"attendance-gateway/internal/types"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func NormalAttendHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := (&websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}).Upgrade(w, r, nil) //ws升级
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, errorx.NewCodeError(10002, err.Error()))
		}
		v := r.URL.Query()
		var req types.NormalAttReq
		req.SupervisroID = v.Get("supervisor_id")
		req.CourseID = v.Get("course_id")
		req.University = v.Get("university")
		if req.CourseID == "" || req.University == "" || req.SupervisroID == "" {
			httpx.ErrorCtx(r.Context(), w, errorx.NewCodeError(10004, "参数为空"))
		}
		cli := &model.Client{
			ID:        req.SupervisroID,
			Socket:    conn,
			HeartBeat: time.Now().Unix(),
			Done:      make(chan struct{}),
			Send:      make(chan []byte),
		}
		fmt.Println("sid:", req.SupervisroID)
		model.Manager.Register <- cli
		go cli.Readmsg(&req, svcCtx)
		go cli.Writemsg(&req, svcCtx)
		go cli.LinkWith()
		l := logic.NewNormalAttendLogic(r.Context(), svcCtx)
		resp, err := l.NormalAttend(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, errorx.NewCodeError(int(resp.Status), err.Error()))
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

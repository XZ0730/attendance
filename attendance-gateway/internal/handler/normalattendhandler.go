package handler

import (
	"attendance-gateway/internal/logic"
	"attendance-gateway/internal/model"
	"attendance-gateway/internal/svc"
	"attendance-gateway/internal/types"
	"errors"
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
			httpx.ErrorCtx(r.Context(), w, err)
		}
		var req types.NormalAttReq
		cli := &model.Client{
			ID:        "1",
			Socket:    conn,
			HeartBeat: time.Now().Unix(),
			Done:      make(chan struct{}),
			Send:      make(chan []byte),
		}
		model.Manager.Register <- cli
		go cli.Readmsg(&req, svcCtx)

		v := r.URL.Query()
		req.CourseID = v.Get("course_id")
		req.University = v.Get("university")
		req.StudentID = v.Get("student_id")
		if req.CourseID == "" || req.University == "" || req.StudentID == "" {
			httpx.ErrorCtx(r.Context(), w, errors.New("参数为空"))
		}
		go cli.Writemsg(&req, svcCtx)
		l := logic.NewNormalAttendLogic(r.Context(), svcCtx)
		resp, err := l.NormalAttend(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

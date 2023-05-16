package logic

import (
	"attend/attendservice"
	"context"
	"strconv"

	"attendance-gateway/internal/svc"
	"attendance-gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LocationAttendLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLocationAttendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LocationAttendLogic {

	return &LocationAttendLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LocationAttendLogic) LocationAttend(req *types.LocationAttReq) (resp *types.AttResponse, err error) {
	// todo: add your logic here and delete this line
	rsp, err := l.svcCtx.Attendservice.LocationAttend(l.ctx, &attendservice.LocationAttRequest{
		StudentId:    req.StudentID,
		CouseMain:    strconv.Itoa(int(req.CourseMain)),
		Longitude:    req.Longitude,
		Latitude:     req.Latitude,
		SupervisorId: req.SupervisorID,
	})
	return &types.AttResponse{
		Status:  rsp.Status,
		Message: rsp.Message,
	}, err
}

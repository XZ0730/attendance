package logic

import (
	"attend/attendserviceclient"
	"context"
	"fmt"

	"attendance-gateway/internal/svc"
	"attendance-gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PullAttandanceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPullAttandanceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PullAttandanceLogic {
	return &PullAttandanceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PullAttandanceLogic) PullAttandance(req *types.PullAttRequest) (resp *types.AttResponse, err error) {
	// todo: add your logic here and delete this line
	// par := &attendservice.PullAttRequest{
	// SupervisorID: req.SupervisorID,
	// CourseID:     req.CourseID,
	// Longitude:    req.Longitude,
	// Latitude:     req.Latitude,
	// }
	rsp, err := l.svcCtx.Attendservice.PullAttendance(l.ctx, &attendserviceclient.PullAttRequest{
		SupervisorID: req.SupervisorID,
		CourseID:     req.CourseID,
		Longitude:    req.Longitude,
		Latitude:     req.Latitude,
		CourseMain:   req.CourseMain,
	})
	fmt.Println("rsp:", rsp)
	return &types.AttResponse{
		Status:  rsp.Status,
		Message: rsp.Message,
		Error:   rsp.Error,
	}, err
}

package logic

import (
	"context"

	"attendance-gateway/rpc/attend/attendservice"
	"attendance-gateway/rpc/attend/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type PullAttendanceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPullAttendanceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PullAttendanceLogic {
	return &PullAttendanceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PullAttendanceLogic) PullAttendance(in *attendservice.PullAttRequest) (*attendservice.AttResponse, error) {
	// todo: add your logic here and delete this line

	return &attendservice.AttResponse{}, nil
}

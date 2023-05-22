package logic

import (
	"attend/attendservice"
	"context"

	"attendance-gateway/internal/svc"
	"attendance-gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAttListByCourseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAttListByCourseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAttListByCourseLogic {
	return &GetAttListByCourseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAttListByCourseLogic) GetAttListByCourse(req *types.GetAttListByCourseReq) (resp *types.AttResponse, err error) {
	// todo: add your logic here and delete this line
	rsp, err := l.svcCtx.Attendservice.GetAttendListByCourse(l.ctx, &attendservice.GetAttListByCourseReq{
		CourseMain: req.CourseMain,
		Week:       req.Week,
	})
	if err != nil {
		return &types.AttResponse{
			Status:  rsp.Status,
			Message: rsp.Message,
		}, nil
	}
	return &types.AttResponse{
		Status:  rsp.Status,
		Data:    rsp.CourseMember,
		Total:   rsp.Total,
		Message: rsp.Message,
	}, nil
}

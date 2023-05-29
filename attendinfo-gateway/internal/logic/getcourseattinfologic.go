package logic

import (
	"context"
	"info/attendinfo"

	"attendinfo-gateway/internal/svc"
	"attendinfo-gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCourseAttInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCourseAttInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCourseAttInfoLogic {
	return &GetCourseAttInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCourseAttInfoLogic) GetCourseAttInfo(req *types.CourseInfoRequest) (resp *types.CourseInfoResponse, err error) {
	// todo: add your logic here and delete this line
	rsp, err := l.svcCtx.InfoCli.GetCourseAttInfo(l.ctx, &attendinfo.CourseInfoRequest{
		University:  req.University,
		Courseid:    req.Course_id,
		Course_Main: req.Course_main,
	})
	return &types.CourseInfoResponse{
		Status:   rsp.Status,
		DataInfo: rsp.CourseAttInfos,
		Total:    rsp.Total,
		Error:    rsp.Error,
	}, nil
}

package logic

import (
	"attend/attendservice"
	"context"

	"attendance-gateway/internal/svc"
	"attendance-gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type NormalAttendLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewNormalAttendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NormalAttendLogic {
	return &NormalAttendLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *NormalAttendLogic) NormalAttend(req *types.NormalAttReq) (resp *types.AttResponse, err error) {
	// todo: add your logic here and delete this line
	if req.Type != 1 {
		ar, err2 := l.svcCtx.Attendservice.NormalAttend(l.ctx, &attendservice.NormalReqest{
			University: req.University,
			CourseId:   req.CourseID,
			StudentId:  req.StudentID,
		})
		if err2 != nil {
			return &types.AttResponse{
				Status: 00011,
				Error:  err2.Error(),
			}, nil
		}
		return &types.AttResponse{
			Status:  ar.Status,
			Data:    ar.CourseMember,
			Total:   ar.Total,
			Message: ar.Message,
			Error:   ar.Error,
		}, nil
	} else {
		ar, err2 := l.svcCtx.Attendservice.AttMember(l.ctx, &attendservice.NormalReqest{
			University: req.University,
			CourseId:   req.CourseID,
			StudentId:  req.StudentID,
			CourseMain: req.CourseMain,
		})
		return &types.AttResponse{
			Status:  ar.Status,
			Message: ar.Message,
			Error:   ar.Error,
		}, err2
	}

}

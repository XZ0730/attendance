package logic

import (
	"context"

	"appeal-gateway/internal/svc"
	"appeal-gateway/internal/types"
	"appeal-gateway/rpc/appeal/appeal"

	"github.com/zeromicro/go-zero/core/logx"
)

type StudentAskforLeaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStudentAskforLeaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StudentAskforLeaveLogic {
	return &StudentAskforLeaveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StudentAskforLeaveLogic) StudentAskforLeave(req *types.AppealRequest) (resp *types.AppealResponse, err error) {
	// todo: add your logic here and delete this line
	res, err := l.svcCtx.Appealer.StudentAskforLeave(l.ctx, &appeal.AppealRequest{
		StudentID:      req.StudentID,
		ContactPhone:   req.ContactPhone,
		EmergencyName:  req.EmergencyName,
		EmergencyPhone: req.EmergencyPhone,
		//辅导员信息
		CounsellorName: "zzb",
		CounsellorID:   req.CounsellorID,
		//申诉-请假理由
		LeaveReason: req.LeaveReason,
		//申诉-请假课程
		CourseName:      req.CourseName,
		CourseID:        req.CourseID,
		LeaveCourseFrom: int32(req.LeaveCourseFrom),
		LeaveCourseTo:   int32(req.LeaveCourseTo),
		//申诉表-请假条区分
		TagAs: uint32(req.TagAs),
	})
	if err != nil {
		return &types.AppealResponse{
			Status:  res.Status,
			Message: res.Message,
			Error:   res.Error,
		}, err
	}

	return &types.AppealResponse{
		Status:  res.Status,
		Message: res.Message,
		Error:   res.Error,
	}, nil
}

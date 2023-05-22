package logic

import (
	"attend/common/errorx"
	"attend/model"
	"context"
	"encoding/json"
	"strconv"

	"attend/attendservice"
	"attend/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAttendListByCourseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAttendListByCourseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAttendListByCourseLogic {
	return &GetAttendListByCourseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAttendListByCourseLogic) GetAttendListByCourse(in *attendservice.GetAttListByCourseReq) (*attendservice.AttNorResponse, error) {
	// todo: add your logic here and delete this line
	//根据课程主id查询  然后查询redis映射
	cr := &model.Course{}
	err := l.svcCtx.DB.Table("course").Where("id=?", in.GetCourseMain()).First(&cr).Error
	if err != nil {
		return &attendservice.AttNorResponse{
			Status:  errorx.CourseISNoExist,
			Message: errorx.GetERROR(errorx.CourseISNoExist),
			Error:   err.Error(),
		}, nil
	}
	// cm := model.GetCourseMem(cr.CourseId, cr.University)
	memMap := make(map[string]*attendservice.CourseMember, 0)
	res, err2 := l.svcCtx.RDB6.HGet(l.ctx, strconv.Itoa(int(cr.Id)), strconv.Itoa(int(in.GetWeek()))).Result()
	if err2 != nil {
		return &attendservice.AttNorResponse{
			Status:  37788,
			Message: "此周未点名，考勤信息为空",
			Error:   err2.Error(),
		}, nil
	}
	err3 := json.Unmarshal([]byte(res), &memMap)
	if err3 != nil {
		return &attendservice.AttNorResponse{
			Status:  37789,
			Message: "考勤信息获取失败",
			Error:   err3.Error(),
		}, nil
	}
	cm := make([]*attendservice.CourseMember, 0)
	for _, v := range memMap {
		cm = append(cm, v)
	}
	return &attendservice.AttNorResponse{
		Status:       errorx.SUCCESS,
		CourseMember: cm,
		Message:      errorx.GetERROR(errorx.SUCCESS),
	}, nil
}

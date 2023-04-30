package logic

import (
	"context"
	"fmt"

	"appeal/appeal"
	"appeal/internal/svc"
	"appeal/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type ComplainToSupervisorLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewComplainToSupervisorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ComplainToSupervisorLogic {
	return &ComplainToSupervisorLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// role --> 1.学生 2.督导员 3.辅导员
func (l *ComplainToSupervisorLogic) ComplainToSupervisor(in *appeal.ComplainRequest) (*appeal.ComplainResponse, error) {
	stu := &model.Character{}
	//还要增加的功能--->辅导员确认-->扣除学生督导信誉分-->删除
	//辅导员认证不属实--->删除

	err := l.svcCtx.MysqlDB.Table("character_msg").
		Where("code=? AND role=?", in.GetStudentID(), 1).
		First(&stu).Error
	if err != nil {
		return &appeal.ComplainResponse{
			Status:  40001,
			Message: "当前用户不存在",
		}, nil
	}
	ct := &model.ComplainTable{
		StudentID:     in.StudentID,
		Student_Major: stu.Major,
		StudentName:   stu.Name,
		College:       stu.College,
	}
	//判断投诉的督导员是否存在
	stu1 := &model.Character{}
	err = l.svcCtx.MysqlDB.Table("character_msg").
		Where("code=? AND name=? AND role=?", in.GetSupervisorID(), in.GetSupervisorName(), 2).
		First(&stu1).Error
	if err != nil {
		return &appeal.ComplainResponse{
			Status:  40002,
			Message: "投诉的督导员不存在",
		}, nil
	}
	fmt.Println("stu1", stu1)
	ct.SupervisorID = stu1.Code
	ct.SupervisorName = in.GetSupervisorName()
	ct.Supervisor_Major = stu1.Major
	ct.Supervisor_College = stu1.College
	ct.SchoolName = stu1.University
	stu2 := &model.Character{}
	err = l.svcCtx.MysqlDB.Table("character_msg").
		Where("code=? AND name=? AND role=?", in.GetCounsellorID(), in.GetCounsellorName(), 3).
		First(&stu2).Error
	if err != nil {
		return &appeal.ComplainResponse{
			Status:  40003,
			Message: "辅导员不存在",
		}, nil
	}
	ct.CounsellorName = stu2.Name
	ct.CounsellorID = in.GetCounsellorID()
	ct.Reason = in.GetReason()
	// todo: add your logic here and delete this line
	err2 := l.svcCtx.MysqlDB.Model(&model.ComplainTable{}).Create(&ct).Error
	if err2 != nil {
		return &appeal.ComplainResponse{
			Status:  40033,
			Message: "投诉发送失败",
			Error:   err2.Error(),
		}, err2
	}
	return &appeal.ComplainResponse{
		Status:  200,
		Message: "successful",
	}, nil
}

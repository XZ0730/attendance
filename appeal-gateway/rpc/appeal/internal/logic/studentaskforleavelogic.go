package logic

import (
	"context"
	"strconv"

	"appeal-gateway/rpc/appeal/appeal"
	"appeal-gateway/rpc/appeal/internal/svc"
	"appeal-gateway/rpc/appeal/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type StudentAskforLeaveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewStudentAskforLeaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StudentAskforLeaveLogic {
	return &StudentAskforLeaveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *StudentAskforLeaveLogic) StudentAskforLeave(in *appeal.AppealRequest) (*appeal.AppealResponse, error) {
	// todo: add your logic here and delete this line
	status, err := l.svcCtx.RDB.Get(context.TODO(), strconv.Itoa(int(in.CourseID))).Result()
	if err == nil && status != "1" { //等于1代表正在上课
		//查得到数据并且正在上课，那么回直接返回，不允许进行申诉或者请假
		return &appeal.AppealResponse{
			Status:  200,
			Message: "当前课程正在上课,请下课后再重试",
		}, nil
	}
	//这边要有一个逻辑:等待课时结构完成-查询该学生是否有选这个课，如果没有就直接返回
	//查询到的课程信息的话进行匹配一下,查询该课程所有课时
	//比较请求中的from,如果是请假条，则from不能从已经结束的课程进行请假
	//如果是申诉表，from只能是上次结束的课程，且两种情况，from和to都不能超过课时选择上限
	//---
	//课程存在且不在上课
	if in.LeaveCourseFrom > in.LeaveCourseTo || in.LeaveCourseFrom <= 0 {
		return &appeal.AppealResponse{
			Status:  3000,
			Message: "课时选择错误",
		}, nil
	}
	if in.TagAs == 2 {
		if in.GetLeaveCourseFrom() != in.GetLeaveCourseTo() {
			return &appeal.AppealResponse{
				Status:  3001,
				Message: "申诉表:课时选择错误",
			}, nil
		}
	}
	//申诉表---还要查询appeal表查看是否已经有申诉未审核的记录，如果有则不允许在继续申诉，
	//请假也需要判断当前是否已经存在form1<from2<to1  from1<to2<to1
	//存在这种情况则直接返回
	var cnt1 int64
	var cnt2 int64
	l.svcCtx.MysqlDB.Model(model.LeaveTable{}).
		Where("student_id=? AND course_id=? AND leave_course_from<=? AND leave_course_to>=? AND is_audit!=?",
			in.GetStudentID(), in.GetCourseID(), in.GetLeaveCourseFrom(), in.GetLeaveCourseFrom(), 3).
		Count(&cnt1)
	l.svcCtx.MysqlDB.Model(model.LeaveTable{}).
		Where("student_id=? AND course_id=? AND leave_course_from<=? AND leave_course_to>=? AND is_audit!=?",
			in.GetStudentID(), in.GetCourseID(), in.GetLeaveCourseTo(), in.GetLeaveCourseTo(), 3).
		Count(&cnt2)
	if cnt1 != 0 || cnt2 != 0 {
		return &appeal.AppealResponse{
			Status:  30034,
			Message: "您选择的课时已经审核通过或者还在审核中",
		}, nil
	}
	lt := &model.LeaveTable{
		StudentID: in.StudentID,

		StudentName:  "xx",
		StudentMajor: "计算机",
		StudentClass: "2",
		College:      "计算机于大数据学院",
		//手机号联系方式
		ContactPhone:   in.GetContactPhone(),
		EmergencyName:  in.GetEmergencyName(),
		EmergencyPhone: in.GetEmergencyPhone(),
		//学校名称
		SchoolName: "福州大学",
		//辅导员信息
		CounsellorName: "zzb",
		CounsellorID:   in.GetCounsellorID(),
		//申诉-请假理由
		LeaveReason: in.LeaveReason,
		//申诉-请假课程
		CourseName:      in.GetCourseName(),
		CourseID:        in.GetCourseID(),
		LeaveCourseFrom: int(in.GetLeaveCourseFrom()),
		LeaveCourseTo:   int(in.GetLeaveCourseTo()),
		//申诉表-请假条区分
		TagAs: uint(in.GetTagAs()),
	}
	//这边--等待mq模块完成--放入mq中，然后返回完成，后续交给mq去处理
	// tx := l.svcCtx.MysqlDB.Begin(&sql.TxOptions{//开启事务
	// 	Isolation: sql.LevelReadCommitted,
	// })
	err2 := l.svcCtx.MysqlDB.Model(model.LeaveTable{}).Create(&lt).Error
	if err2 != nil {
		return &appeal.AppealResponse{
			Status:  3004,
			Message: "提交失败",
		}, err2
	}
	return &appeal.AppealResponse{
		Status:  200,
		Message: strconv.Itoa(int(in.StudentID)),
	}, nil
}

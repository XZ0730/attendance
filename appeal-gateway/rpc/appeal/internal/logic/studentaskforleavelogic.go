package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"appeal/appeal"
	"appeal/internal/svc"
	"appeal/model"
	"mq_server/mq"

	"github.com/cyfckl2012/disgo"
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
	lock, _ := disgo.GetLock(l.svcCtx.RDB3, "test")
	succ, err3 := lock.TryLock(l.svcCtx.RDB3.Context(), 5*time.Second, 10*time.Second)
	if !succ {
		return &appeal.AppealResponse{
			Status:  30017,
			Message: "当前系统繁忙，稍后再试",
			Error:   err3.Error(),
		}, err3
	}
	fmt.Println("-------------------测试位置22-----------------------")
	status, err := l.svcCtx.RDB3.
		HGet(context.TODO(), in.GetUniversity(), in.GetCourseID()).Result()
	if err == nil && status == "1" { //等于1代表正在上课
		//查得到数据并且正在上课，那么回直接返回，不允许进行申诉或者请假
		return &appeal.AppealResponse{
			Status:  30014,
			Message: "当前课程正在上课,请下课后再重试",
		}, nil
	}
	if err != nil {
		fmt.Println("errrrrrrrrr：", err)
	}
	fmt.Println("status:", status)
	fmt.Println("-------------------测试位置23------------------------")
	succ, _ = lock.Release(l.svcCtx.RDB3.Context())
	if !succ {
		return &appeal.AppealResponse{
			Status:  30018,
			Message: "当前系统繁忙，稍后再试",
			Error:   err3.Error(),
		}, err3
	}
	//这边要有一个逻辑:等待课时结构完成-查询该学生是否有选这个课，如果没有就直接返回
	//查询到的课程信息的话进行匹配一下,查询该课程所有课时
	//课程信息不能为空 ---- 前端会搜索课程，点击会发送请求带有课程信息-- 所以这边不用判断是否学生有无选择这个课程
	//一门课程对应 两个课时，
	//比较请求中的from,如果是请假条，则from不能从已经结束的课程进行请假
	//如果是申诉表，from只能是上次结束的课程，且两种情况，from和to都不能超过课时选择上限
	//---
	//课程存在且不在上课
	fmt.Println("-------------------测试位置1------------------------")
	course := &model.Course{}
	err1 := l.svcCtx.MysqlDB.Table("course").
		Where("course_id=?", in.GetCourseID()).
		First(&course).Error
	fmt.Println("-------------------测试位置2------------------------")
	if err1 != nil {
		fmt.Println("ERR1:", err1.Error())
		return &appeal.AppealResponse{
			Status:  30005,
			Message: "课程貌似不存在捏",
		}, err1
	}
	fmt.Println("-------------------测试位置3------------------------")
	if in.GetLeaveCourseFrom() > in.GetLeaveCourseTo() ||
		in.GetLeaveCourseFrom() < course.WeekStart ||
		in.GetLeaveCourseTo() > course.WeekEnd {
		return &appeal.AppealResponse{
			Status:  30000,
			Message: "课时选择错误",
		}, nil
	}
	if in.TagAs == 2 { //2是申诉表
		latest := 0
		fmt.Println("course:", course)
		res, err := l.svcCtx.RDB.HGet(l.svcCtx.RDB.Context(), course.University, in.GetCourseID()).Result()
		if err != nil {
			fmt.Println("errrrrr--", err)
		} else {
			latest, _ = strconv.Atoi(res)
			fmt.Println("la:", latest)
		}

		if latest != int(in.GetLeaveCourseFrom()) {
			return &appeal.AppealResponse{
				Status:  30001,
				Message: "申诉课时已经结束或者未开始，只能选择最近上完的课时捏",
			}, nil
		}
		if in.GetLeaveCourseFrom() != in.GetLeaveCourseTo() {
			return &appeal.AppealResponse{
				Status:  30001,
				Message: "申诉表:课时选择错误",
			}, nil
		}
		var cnt3 int64 = 0
		//这边到时候补个逻辑，
		l.svcCtx.MysqlDB.Table("leave_table").
			Where("student_id=? AND course_id=? AND leave_course_from=? AND course_name=?",
				in.GetStudentID(), course.CourseId, in.GetLeaveCourseFrom(), course.Name).
			Count(&cnt3)
		if cnt3 != 0 {
			return &appeal.AppealResponse{
				Status:  30008,
				Message: "该课时已经申诉过了",
			}, nil
		}
	}
	//申诉表---还要查询appeal表查看是否已经有申诉未审核的记录，如果有则不允许在继续申诉，
	//请假也需要判断当前是否已经存在form1<from2<to1  from1<to2<to1
	//存在这种情况则直接返回
	var cnt1 int64 = 0
	var cnt2 int64 = 0
	l.svcCtx.MysqlDB.Model(model.LeaveTable{}).
		Where("student_id=? AND course_id=? AND leave_course_from<=? AND leave_course_to>=? AND is_audit!=? AND course_id=? AND course_name=?",
			in.GetStudentID(), in.GetCourseID(), in.GetLeaveCourseFrom(), in.GetLeaveCourseFrom(), 3, in.GetCourseID(), in.GetCourseName()).
		Count(&cnt1)
	l.svcCtx.MysqlDB.Model(model.LeaveTable{}).
		Where("student_id=? AND course_id=? AND leave_course_from<=? AND leave_course_to>=? AND is_audit!=? AND course_id=? AND course_name=?",
			in.GetStudentID(), in.GetCourseID(), in.GetLeaveCourseTo(), in.GetLeaveCourseTo(), 3, in.GetCourseID(), in.GetCourseName()).
		Count(&cnt2)
	if cnt1 != 0 || cnt2 != 0 {
		return &appeal.AppealResponse{
			Status:  30034,
			Message: "您选择的课时已经审核通过或者还在审核中",
		}, nil
	}
	stu := &model.Character{}
	l.svcCtx.MysqlDB.Table("character_msg").
		Where("code=?", in.GetStudentID()).
		First(&stu)
	lt := &model.LeaveTable{
		StudentID: in.StudentID,

		StudentName:  stu.Name,
		StudentMajor: stu.Major,
		StudentClass: stu.ClassNum,
		College:      stu.College,
		//手机号联系方式
		ContactPhone:   in.GetContactPhone(),
		EmergencyName:  in.GetEmergencyName(),
		EmergencyPhone: in.GetEmergencyPhone(),
		//学校名称
		SchoolName: "福州大学",
		//辅导员信息
		CounsellorName: in.GetCounsellorName(),
		CounsellorID:   in.GetCounsellorID(),
		//申诉-请假理由
		LeaveReason: in.LeaveReason,
		//申诉-请假课程
		CourseName:      course.Name,
		CourseID:        course.CourseId,
		LeaveCourseFrom: int(in.GetLeaveCourseFrom()),
		LeaveCourseTo:   int(in.GetLeaveCourseTo()),
		//申诉表-请假条区分
		TagAs: uint(in.GetTagAs()),
	}
	//这边--等待mq模块完成--放入mq中，然后返回完成，后续交给mq去处理
	// tx := l.svcCtx.MysqlDB.Begin(&sql.TxOptions{//开启事务
	// 	Isolation: sql.LevelReadCommitted,
	// })
	ltt, err4 := json.Marshal(*lt)
	if err4 != nil {
		return &appeal.AppealResponse{
			Status: 30014,
			Error:  err4.Error(),
		}, err4
	}
	// l.svcCtx.MQ.PublishTo()
	// err2 := l.svcCtx.MysqlDB.Model(model.LeaveTable{}).Create(&lt).Error
	// if err2 != nil {
	// 	return &appeal.AppealResponse{
	// 		Status:  3004,
	// 		Message: "提交失败",
	// 	}, err2
	// }
	req := &mq.LeaveRequest{
		Queuename: "leavemq",
		Data:      string(ltt),
	}
	rsp, err := l.svcCtx.MQ.PublishLeave(l.ctx, req)
	return &appeal.AppealResponse{
		Status:  200,
		Message: rsp.Message,
	}, nil
}

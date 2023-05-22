package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"appeal/appeal"
	"appeal/common/errorx"
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
	fmt.Println("测试2")
	lock, _ := disgo.GetLock(l.svcCtx.RDB3, "ttt")
	fmt.Println("测试3")
	succ, err3 := lock.TryLock(l.ctx, 5*time.Second, 10*time.Second)
	if !succ {
		return &appeal.AppealResponse{
			Status:  errorx.BusySysError,
			Message: errorx.GetERROR(errorx.BusySysError),
			Error:   err3.Error(),
		}, nil
	}
	fmt.Println("测试1")
	fmt.Println("-------------------测试位置22-----------------------")
	status, err := l.svcCtx.RDB3.
		HGet(context.TODO(), in.GetUniversity(), in.GetCourseID()).Result()
	if err == nil && status == "1" { //等于1代表正在上课
		//查得到数据并且正在上课，那么回直接返回，不允许进行申诉或者请假
		lock.Release(l.ctx)
		return &appeal.AppealResponse{
			Status:  errorx.COURSE_ING,
			Message: errorx.GetERROR(errorx.COURSE_ING),
		}, nil
	}
	if err != nil {
		fmt.Println("errrrrrrrrr：", err)
	}
	fmt.Println("status:", status)
	fmt.Println("-------------------测试位置23------------------------")
	succ, _ = lock.Release(l.ctx)
	if !succ {
		return &appeal.AppealResponse{
			Status:  errorx.BusySysError,
			Message: errorx.GetERROR(errorx.BusySysError),
			Error:   err3.Error(),
		}, nil
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
		Where("course_id=? AND university=?", in.GetCourseID(), in.GetUniversity()).
		First(&course).Error
	fmt.Println("-------------------测试位置2------------------------")
	if err1 != nil {
		fmt.Println("ERR1:", err1.Error())
		return &appeal.AppealResponse{
			Status:  errorx.CourseISNoExist,
			Message: errorx.GetERROR(errorx.CourseISNoExist),
			Error:   err1.Error(),
		}, nil
	}
	fmt.Println("-------------------测试位置3------------------------")
	if in.GetLeaveCourseFrom() > in.GetLeaveCourseTo() ||
		in.GetLeaveCourseFrom() < course.WeekStart ||
		in.GetLeaveCourseTo() > course.WeekEnd {
		return &appeal.AppealResponse{
			Status:  errorx.CourseSelectError1,
			Message: errorx.GetERROR(errorx.CourseSelectError1),
		}, nil
	}
	fmt.Println("测试2")
	latest := 0 //最近结束周
	res, err := l.svcCtx.RDB.HGet(l.svcCtx.RDB.Context(), course.University, strconv.Itoa(int(course.Id))).Result()
	if err != nil {
		fmt.Println("errrrrr--", err)
		return &appeal.AppealResponse{
			Status:  errorx.RedisInitError,
			Message: errorx.GetERROR(errorx.RedisInitError),
			Error:   err.Error(),
		}, nil
	} else {
		latest, _ = strconv.Atoi(res)
		fmt.Println("la:", latest)
	}
	fmt.Println("测试3")
	if in.TagAs == 2 { //2是申诉表

		fmt.Println("course:", course)

		if latest != int(in.GetLeaveCourseFrom()) {
			return &appeal.AppealResponse{
				Status:  errorx.CourseSelectError2,
				Message: errorx.GetERROR(errorx.CourseSelectError2),
			}, nil
		}
		if in.GetLeaveCourseFrom() != in.GetLeaveCourseTo() {
			return &appeal.AppealResponse{
				Status:  errorx.CourseSelectError1,
				Message: errorx.GetERROR(errorx.CourseSelectError1),
			}, nil
		}
		//判断上周是否点名了--没点名--直接返回  点名了继续
		id := strconv.Itoa(int(course.Id))
		result, err := l.svcCtx.RDB7.HGet(l.ctx, id, res).Result()
		if err != nil {
			return &appeal.AppealResponse{
				Status: errorx.RedisInitError,
				Error:  errorx.GetERROR(errorx.RedisInitError),
			}, nil
		}
		fmt.Println("result:", result)
		if result == "0" {
			return &appeal.AppealResponse{
				Status: errorx.LastWeekError,
				Error:  errorx.GetERROR(errorx.LastWeekError),
			}, nil
		}
		var cnt3 int64 = 0
		//这边到时候补个逻辑，
		l.svcCtx.MysqlDB.Table("leave_table").
			Where("student_id=? AND course_id=? AND leave_course_from=? AND course_name=? AND school_name=?",
				in.GetStudentID(), course.CourseId, in.GetLeaveCourseFrom(), course.Name, in.GetUniversity()).
			Count(&cnt3)
		if cnt3 != 0 {
			return &appeal.AppealResponse{
				Status:  errorx.RepeteAppealError,
				Message: errorx.GetERROR(errorx.RepeteAppealError),
			}, nil
		}

	} else {
		var cnt3 int64 = 0
		l.svcCtx.MysqlDB.Table("leave_table").
			Where("student_id=? AND course_id=? AND leave_course_from=? AND course_name=? AND school_name=?",
				in.GetStudentID(), course.CourseId, in.GetLeaveCourseFrom(), course.Name, in.GetUniversity()).
			Count(&cnt3)
		if cnt3 != 0 {
			return &appeal.AppealResponse{
				Status:  errorx.RepeteAppealError,
				Message: errorx.GetERROR(errorx.RepeteAppealError),
			}, nil
		}
		if in.GetLeaveCourseFrom() <= uint32(latest) ||
			in.GetLeaveCourseTo() > course.WeekEnd {
			return &appeal.AppealResponse{
				Status:  errorx.CourseRepeteError,
				Message: errorx.GetERROR(errorx.CourseRepeteError),
			}, nil
		}
	}
	//申诉表---还要查询appeal表查看是否已经有申诉未审核的记录，如果有则不允许在继续申诉，
	//请假也需要判断当前是否已经存在form1<from2<to1  from1<to2<to1
	//存在这种情况则直接返回
	var cnt1 int64 = 0
	var cnt2 int64 = 0
	l.svcCtx.MysqlDB.Model(model.LeaveTable{}).
		Where("student_id=? AND course_id=? AND leave_course_from<=? AND leave_course_to>=? AND is_audit!=? AND school_name=? AND course_name=?",
			in.GetStudentID(), in.GetCourseID(), in.GetLeaveCourseFrom(), in.GetLeaveCourseFrom(), 3, in.GetUniversity(), in.GetCourseName()).
		Count(&cnt1)
	l.svcCtx.MysqlDB.Model(model.LeaveTable{}).
		Where("student_id=? AND course_id=? AND leave_course_from<=? AND leave_course_to>=? AND is_audit!=? AND school_name=? AND course_name=?",
			in.GetStudentID(), in.GetCourseID(), in.GetLeaveCourseTo(), in.GetLeaveCourseTo(), 3, in.GetUniversity(), in.GetCourseName()).
		Count(&cnt2)
	if cnt1 != 0 || cnt2 != 0 {
		return &appeal.AppealResponse{
			Status:  errorx.AUDIT_ING,
			Message: errorx.GetERROR(errorx.AUDIT_ING),
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
		SchoolName: course.University,
		//辅导员信息
		CounsellorName: in.GetCounsellorName(),
		CounsellorID:   in.GetCounsellorID(),
		//申诉-请假理由
		LeaveReason: in.GetLeaveReason(),
		//申诉-请假课程
		CourseName:      course.Name,
		CourseID:        course.CourseId,
		LeaveCourseFrom: int(in.GetLeaveCourseFrom()),
		LeaveCourseTo:   int(in.GetLeaveCourseTo()),
		//申诉表-请假条区分
		TagAs:   uint(in.GetTagAs()),
		IsAudit: 1,
	}
	//这边--等待mq模块完成--放入mq中，然后返回完成，后续交给mq去处理
	// tx := l.svcCtx.MysqlDB.Begin(&sql.TxOptions{//开启事务
	// 	Isolation: sql.LevelReadCommitted,
	// })
	ltt, err4 := json.Marshal(*lt)
	if err4 != nil {
		return &appeal.AppealResponse{
			Status:  errorx.JSON_MARSHAL_ERROR,
			Message: errorx.GetERROR(errorx.JSON_MARSHAL_ERROR),
			Error:   err4.Error(),
		}, nil
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
	if err != nil {
		return &appeal.AppealResponse{
			Status:  errorx.MQ_RETURN_ERROR,
			Message: rsp.Message,
			Error:   err.Error(),
		}, nil
	}
	return &appeal.AppealResponse{
		Status:  errorx.SUCCESS,
		Message: errorx.GetERROR(errorx.SUCCESS),
	}, nil
}

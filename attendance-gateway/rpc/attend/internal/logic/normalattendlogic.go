package logic

import (
	"context"
	"fmt"
	"strconv"

	"attend/attendservice"
	"attend/common/errorx"
	"attend/internal/svc"
	"attend/model"

	"github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/core/logx"
)

type NormalAttendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNormalAttendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NormalAttendLogic {
	return &NormalAttendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *NormalAttendLogic) NormalAttend(in *attendservice.NormalReqest) (*attendservice.AttNorResponse, error) {
	// todo: add your logic here and delete this line
	// model.DB.Table("course_group").
	// 		Select("character_msg.code as code,character_msg.name as student_name,course_group.course_id as course_id").Joins("left join character_msg on course_group.student_id=character_msg.code where course_id=?", "10000").
	// 		Scan(&results)
	mem := model.GetCourseMem(in.GetCourseId(), in.GetUniversity())
	memMap := make(map[string]*attendservice.CourseMember, 0)
	cr := &model.Course{}
	err := model.DB.Table("course").Where("university=? AND course_id=?", in.GetUniversity(), in.GetCourseId()).First(&cr).Error
	if err != nil {
		return &attendservice.AttNorResponse{
			Status:  38998,
			Message: "课程不存在",
			Error:   err.Error(),
		}, nil
	}
	week, err2 := model.RDB.HGet(l.ctx, in.GetUniversity(), strconv.Itoa(int(cr.Id))).Result()
	if err2 != nil {
		return &attendservice.AttNorResponse{
			Status:  errorx.RedisInitError,
			Message: errorx.GetERROR(errorx.RedisInitError),
			Error:   err2.Error(),
		}, nil
	}
	weeknow, _ := strconv.Atoi(week)
	for _, v := range mem {
		if v == nil {
			break
		} else {
			memMap[v.Code] = v
			v.Week = uint32(weeknow + 1)
			v.MissAttend = 1
		}
	}
	sta, err := l.svcCtx.RDB5.ZRangeByScore(l.ctx, strconv.Itoa(int(cr.Id)), &redis.ZRangeBy{
		Max: "1",
		Min: "1",
	}).Result()
	fmt.Println("id:", cr.Id)
	if err != nil {
		fmt.Println("err:", err)
	}
	fmt.Println("sta", sta)
	for _, v := range sta { //出勤设置为2，缺勤设置为1
		memMap[v].MissAttend = 2
	}
	memSlice := make([]*attendservice.CourseMember, 0)
	for _, v := range memMap {
		memSlice = append(memSlice, v)
	}
	fmt.Println("memslice:", memSlice)
	return &attendservice.AttNorResponse{
		Status:       errorx.SUCCESS,
		CourseMember: memSlice,
		Total:        uint32(len(mem)),
		Message:      errorx.GetERROR(errorx.SUCCESS),
	}, nil
}

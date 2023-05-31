package logic

import (
	"attend/common/errorx"
	"attend/model"
	"attend/pkg"
	"context"
	"encoding/json"
	"fmt"
	"mq_server/mqclient"
	"strconv"
	"time"

	"attend/attendservice"
	"attend/internal/svc"

	"github.com/cyfckl2012/disgo"
	"github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/core/logx"
)

type PullAttendanceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPullAttendanceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PullAttendanceLogic {
	return &PullAttendanceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PullAttendanceLogic) PullAttendance(in *attendservice.PullAttRequest) (*attendservice.AttResponse, error) {
	// todo: add your logic here and delete this line
	cr := &model.Course{}
	fmt.Println("--------------------------1---------")
	err6 := l.svcCtx.DB.Table("course").Where("id=?", in.GetCourseMain()).First(&cr).Error
	if err6 != nil {
		return &attendservice.AttResponse{
			Status:  errorx.RecordGetError,
			Message: errorx.GetERROR(errorx.RecordGetError),
		}, nil
	}
	fmt.Println("---------------------------2--------")
	ok := pkg.JudgeTime(cr)
	if !ok {
		return &attendservice.AttResponse{
			Status:  errorx.PuLLAttTimeERROR,
			Message: errorx.GetERROR(errorx.PuLLAttTimeERROR),
		}, nil
	}

	dl, err2 := disgo.GetLock(l.svcCtx.RDB3, "rdb3")
	if err2 != nil {
		return &attendservice.AttResponse{
			Status:  errorx.BusySysError,
			Message: errorx.GetERROR(errorx.BusySysError),
			Error:   err2.Error(),
		}, err2
	}
	fmt.Println("----------测试1----------")
	_, err3 := dl.TryLock(l.ctx, 5*time.Second, 10*time.Second)
	if err3 != nil {
		return &attendservice.AttResponse{
			Status:  errorx.BusySysError,
			Message: errorx.GetERROR(errorx.BusySysError),
			Error:   err3.Error(),
		}, err3
	}
	fmt.Println("----------测试2----------")
	//获取点名状态 检测是否正在点名
	status, err := l.svcCtx.RDB3.ZScore(l.ctx, cr.University, strconv.Itoa(int(cr.Id))).Result()
	if err != nil {
		return &attendservice.AttResponse{
			Status:  errorx.RedisInitError,
			Message: errorx.GetERROR(errorx.RedisInitError),
			Error:   err.Error(),
		}, nil
	}
	fmt.Println("----------测试3----------")
	dl.Release(l.ctx)
	if status == 1 { //正在点名中
		return &attendservice.AttResponse{
			Status:  errorx.COURSE_ING,
			Message: errorx.GetERROR(errorx.COURSE_ING),
		}, nil
	}
	err = l.svcCtx.RDB3.ZAdd(l.ctx, cr.University, &redis.Z{
		Score:  1,
		Member: cr.Id,
	}).Err()
	if err != nil {
		return &attendservice.AttResponse{
			Status:  errorx.RedisInitError,
			Message: errorx.GetERROR(errorx.RedisInitError),
			Error:   err.Error(),
		}, nil
	}
	fmt.Println("----------测试4----------")
	err5 := l.svcCtx.RDB4.GeoAdd(l.ctx, strconv.Itoa(int(in.GetCourseMain())), &redis.GeoLocation{
		Longitude: in.GetLongitude(),
		Latitude:  in.GetLatitude(),
		Name:      in.GetSupervisorID(),
	}).Err()
	if err5 != nil {
		return &attendservice.AttResponse{
			Status:  errorx.GeoERROR,
			Message: errorx.GetERROR(errorx.GeoERROR),
			Error:   err5.Error(),
		}, nil
	}
	mar := &model.MarshalPull{
		SupervisorID: in.GetSupervisorID(),
		CourseID:     in.GetCourseID(),
		University:   cr.University,
		Longitude:    in.GetLongitude(),
		Latitude:     in.GetLatitude(),
	}
	data, err4 := json.Marshal(*mar)
	if err4 != nil {
		return &attendservice.AttResponse{
			Status:  errorx.JSON_MARSHAL_ERROR,
			Message: errorx.GetERROR(errorx.JSON_MARSHAL_ERROR),
		}, nil
	}
	rsp, err := l.svcCtx.MQ.PublishPull(l.ctx, &mqclient.AttRequest{
		Tag:  1,
		Data: string(data),
	})
	if err != nil {
		return &attendservice.AttResponse{
			Status:  errorx.MQ_RETURN_ERROR,
			Message: rsp.Message,
			Error:   err.Error(),
		}, nil
	}
	return &attendservice.AttResponse{
		Status:  errorx.SUCCESS,
		Message: rsp.Message,
	}, nil
}

package logic

import (
	"context"
	"fmt"

	"attend/attendservice"
	"attend/common/errorx"
	"attend/internal/svc"

	"github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/core/logx"
)

type LocationAttendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLocationAttendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LocationAttendLogic {
	return &LocationAttendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LocationAttendLogic) LocationAttend(in *attendservice.LocationAttRequest) (*attendservice.AttResponse, error) {
	// todo: add your logic here and delete this line
	fmt.Println("-----------------------")
	l.svcCtx.RDB4.GeoAdd(l.ctx, in.GetCouseMain(), &redis.GeoLocation{
		Longitude: in.GetLongitude(),
		Latitude:  in.GetLatitude(),
		Name:      in.GetStudentId(),
	})
	flo, err := l.svcCtx.RDB4.GeoDist(l.ctx, in.GetCouseMain(), in.GetSupervisorId(), in.GetStudentId(), "").Result()
	if err != nil {
		fmt.Println("err:", err.Error())
		return &attendservice.AttResponse{
			Status:  errorx.LocationDistERROR,
			Message: errorx.GetERROR(errorx.LocationDistERROR),
			Error:   err.Error(),
		}, nil
	}
	fmt.Println("flo:", flo)
	if flo > 15 {
		return &attendservice.AttResponse{
			Status:  errorx.OverERROR,
			Message: errorx.GetERROR(errorx.OverERROR),
		}, nil
	}
	err = l.svcCtx.RDB5.ZAdd(l.ctx, in.GetCouseMain(), &redis.Z{
		Score:  1,
		Member: in.GetStudentId(),
	}).Err()
	if err != nil {
		return &attendservice.AttResponse{
			Status:  errorx.FailedAttend,
			Message: errorx.GetERROR(errorx.FailedAttend),
		}, nil
	}
	return &attendservice.AttResponse{
		Status:  errorx.SUCCESS,
		Message: errorx.GetERROR(errorx.SUCCESS),
	}, nil
}

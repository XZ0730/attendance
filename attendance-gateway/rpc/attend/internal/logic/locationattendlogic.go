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
	if flo > 15 {
		fmt.Println("flo:", flo)
		return &attendservice.AttResponse{
			Status:  errorx.OverERROR,
			Message: errorx.GetERROR(errorx.OverERROR),
		}, nil
	}

	return &attendservice.AttResponse{
		Status:  errorx.SUCCESS,
		Message: errorx.GetERROR(errorx.SUCCESS),
	}, nil
}

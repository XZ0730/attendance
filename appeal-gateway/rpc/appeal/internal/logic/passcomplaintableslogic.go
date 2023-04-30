package logic

import (
	"context"
	"fmt"
	"time"

	"appeal/appeal"
	"appeal/internal/svc"

	// "appeal/model"
	"mq_server/mq"

	"github.com/cyfckl2012/disgo"
	"github.com/zeromicro/go-zero/core/logx"
)

type PassComplainTablesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPassComplainTablesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PassComplainTablesLogic {
	return &PassComplainTablesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PassComplainTablesLogic) PassComplainTables(in *appeal.ComplainPassRequest) (*appeal.AppealResponse, error) {
	// todo: add your logic here and delete this line
	//删除逻辑---通过用户token携带的辅导员id传入
	//通过的话，应该要有个扣除积分的操作
	//未通过则普通删除
	if in.GetPass() == true {
		fmt.Println("扣除积分")

		lock, err := disgo.GetLock(l.svcCtx.RDB, "test")
		if err != nil {
			fmt.Println("err:", err)
		}
		success, err2 := lock.TryLock(l.svcCtx.RDB.Context(), 5*time.Second, 10*time.Second)
		if !success {
			return &appeal.AppealResponse{
				Status:  39908,
				Message: "系统繁忙稍后重试",
			}, err2
		} //这边应该token传入 账户基本信息
		//key:学校名称 member:code  score:信誉分
		err3 := l.svcCtx.RDB2.
			ZIncrBy(l.svcCtx.RDB2.Context(), in.GetSchoolName(), -2, in.GetSupervisorID()).
			Err()
		// l.svcCtx.RDB.GeoAdd(l.svcCtx.RDB.Context(), "car", &redis.GeoLocation{
		// 	Longitude: 22.22,
		// 	Latitude:  22.22222,
		// 	Name:      "supervisor",
		// })
		if err3 != nil {
			fmt.Println("err3:", err3)
		}
		_, err2 = lock.Release(l.svcCtx.RDB.Context())
		if err2 != nil {
			return &appeal.AppealResponse{
				Status:  39908,
				Message: "系统繁忙稍后重试",
			}, err2
		}
	}
	req := &mq.Request{
		Cid:          in.GetCid(),
		CounsellorId: in.GetConsellorID(),
	}
	_, err := l.svcCtx.MQ.Publish(l.ctx, req)
	if err != nil {
		return &appeal.AppealResponse{
			Status:  39901,
			Message: "删除出错",
			Error:   err.Error(),
		}, err
	}
	return &appeal.AppealResponse{
		Status:  200,
		Message: "successful",
	}, nil
}

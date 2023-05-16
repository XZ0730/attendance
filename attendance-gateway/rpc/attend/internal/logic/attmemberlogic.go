package logic

import (
	"appeal/common/errorx"
	"context"
	"strconv"

	"attend/attendservice"
	"attend/internal/svc"

	"github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/core/logx"
)

type AttMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAttMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AttMemberLogic {
	return &AttMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 传入stuid courseid university
func (l *AttMemberLogic) AttMember(in *attendservice.NormalReqest) (*attendservice.AttResponse, error) {
	// todo: add your logic here and delete this line
	//修改redis数据
	courseMain := strconv.Itoa(int(in.GetCourseMain()))
	l.svcCtx.RDB5.ZAdd(l.ctx, courseMain, &redis.Z{
		Score:  1,
		Member: in.GetStudentId(),
	})
	return &attendservice.AttResponse{
		Status:  errorx.SUCCESS,
		Message: errorx.GetERROR(errorx.SUCCESS),
	}, nil
}

package logic

import (
	"context"

	"attendance-gateway/internal/svc"
	"attendance-gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PullAttandanceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPullAttandanceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PullAttandanceLogic {
	return &PullAttandanceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PullAttandanceLogic) PullAttandance(req *types.PullAttRequest) (resp *types.AttResponse, err error) {
	// todo: add your logic here and delete this line

	return
}

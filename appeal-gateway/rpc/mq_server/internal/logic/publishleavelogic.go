package logic

import (
	"context"

	"mq_server/internal/svc"
	"mq_server/mq"
	"mq_server/rabbitmq"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishLeaveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublishLeaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishLeaveLogic {
	return &PublishLeaveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PublishLeaveLogic) PublishLeave(in *mq.LeaveRequest) (*mq.Response, error) {
	// todo: add your logic here and delete this line
	rabbitmq.NewLeaveMQ(in.Queuename).Publish(in.Data)
	return &mq.Response{
		Message: "successful",
	}, nil
}

package logic

import (
	"context"

	"mq_server/internal/svc"
	"mq_server/mq"
	"mq_server/rabbitmq"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishPullLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublishPullLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishPullLogic {
	return &PublishPullLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PublishPullLogic) PublishPull(in *mq.AttRequest) (*mq.Response, error) {
	// todo: add your logic here and delete this line

	rabbitmq.NewDelayMQ("delay1").Publish(in.Data)
	return &mq.Response{
		Message: "successful",
	}, nil
}

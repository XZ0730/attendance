package logic

import (
	"context"
	"encoding/json"

	"mq_server/internal/svc"
	"mq_server/mq"
	"mq_server/rabbitmq"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishLogic {
	return &PublishLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PublishLogic) Publish(in *mq.Request) (*mq.Response, error) {
	// todo: add your logic here and delete this line
	req, err := json.Marshal(in)
	if err != nil {
		return &mq.Response{
			Message: "json marshal error",
		}, err
	}
	rabbitmq.NewPassMQ("passmq").Publish(string(req))
	return &mq.Response{}, nil
}

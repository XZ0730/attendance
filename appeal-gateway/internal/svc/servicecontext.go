package svc

import (
	"appeal-gateway/internal/config"
	"appeal/appealclient"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config   config.Config
	Appealer appealclient.Appeal
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:   c,
		Appealer: appealclient.NewAppeal(zrpc.MustNewClient(c.AppealCli)),
	}
}

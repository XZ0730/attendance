package svc

import (
	"attendinfo-gateway/internal/config"
	"info/attendinfoclient"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config  config.Config
	InfoCli attendinfoclient.Attendinfo
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		InfoCli: attendinfoclient.NewAttendinfo(zrpc.MustNewClient(c.InfoCli)),
	}
}

package svc

import (
	"attend/attendserviceclient"
	"attendance-gateway/internal/config"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config        config.Config
	Attendservice attendserviceclient.Attendservice
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:        c,
		Attendservice: attendserviceclient.NewAttendservice(zrpc.MustNewClient(c.AttendCli)),
	}
}

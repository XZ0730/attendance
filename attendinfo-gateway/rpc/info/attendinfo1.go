package main

import (
	"flag"
	"fmt"

	"info/attendinfo"
	"info/internal/config"
	"info/internal/logic"
	"info/internal/server"
	"info/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/attendinfo.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		attendinfo.RegisterAttendinfoServer(grpcServer, server.NewAttendinfoServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()
	logic.Init()
	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}

package main

import (
	"flag"
	"fmt"

	"mq_server/internal/config"
	"mq_server/internal/server"
	"mq_server/internal/svc"
	"mq_server/mq"
	"mq_server/rabbitmq"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/mq.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		mq.RegisterMqServer(grpcServer, server.NewMqServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	rabbitmq.InitPassMQ()
	rabbitmq.InitLeaveMQ()
	rabbitmq.InitDelayMQ()
	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}

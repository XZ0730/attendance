package main

import (
	"flag"
	"fmt"
	"strconv"

	"attendance-gateway/internal/config"
	"attendance-gateway/internal/handler"
	"attendance-gateway/internal/model"
	"attendance-gateway/internal/svc"

	"github.com/SimonWang00/goeureka"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/attendance-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)
	go model.Manager.Start()
	opt := make(map[string]string)
	opt["username"] = "admin"
	opt["password"] = "123456"
	// 加载配置
	goeureka.RegisterClient("http://47.113.216.236:9737", "43.136.122.18",
		"APP-attend", "8889",
		strconv.Itoa(43), opt)
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

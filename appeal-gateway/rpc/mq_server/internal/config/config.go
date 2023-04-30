package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	MQURL string
	Mysql struct {
		DataSource string
	}
	Addr string
}

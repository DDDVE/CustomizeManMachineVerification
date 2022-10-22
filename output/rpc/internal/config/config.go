package config

import (
	"pkg/model"

	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Mysql model.MysqlConfig
	Redis model.RedisConfig
}

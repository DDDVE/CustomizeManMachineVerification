package config

import (
	"pkg/model"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type (
	Config struct {
		rest.RestConf
		LoginRpc zrpc.RpcClientConf
		Redis    model.RedisConfig
	}
)

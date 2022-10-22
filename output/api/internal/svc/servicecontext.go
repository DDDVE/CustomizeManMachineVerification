package svc

import (
	"api/internal/config"
	"pkg/model"
	"rpc/outputclient"

	"github.com/garyburd/redigo/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config    config.Config
	OutputRpc outputclient.Output
	RedisPool *redis.Pool
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		OutputRpc: outputclient.NewOutput(zrpc.MustNewClient(c.OutputRpc)),
		RedisPool: model.NewRedis(&c.Redis),
	}
}

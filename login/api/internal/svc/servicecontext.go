package svc

import (
	"api/internal/config"
	"pkg/model"
	"rpc/loginclient"

	"github.com/garyburd/redigo/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config    config.Config
	LoginRpc  loginclient.Login
	RedisPool *redis.Pool
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		LoginRpc:  loginclient.NewLogin(zrpc.MustNewClient(c.LoginRpc)),
		RedisPool: model.NewRedis(&c.Redis),
	}
}

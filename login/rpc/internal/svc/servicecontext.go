package svc

import (
	"pkg/model"
	"rpc/internal/config"

	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
)

type ServiceContext struct {
	Config    config.Config
	GormDB    *gorm.DB
	RedisPool *redis.Pool
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		GormDB:    model.NewMysql(&c.Mysql),
		RedisPool: model.NewRedis(&c.Redis),
	}
}

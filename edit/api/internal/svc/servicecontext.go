package svc

import (
	"api/internal/config"

	"rpc/editclient"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config  config.Config
	EditRpc editclient.Edit
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		EditRpc: editclient.NewEdit(zrpc.MustNewClient(c.EditRpc)),
	}
}

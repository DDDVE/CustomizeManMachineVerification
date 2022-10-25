package main

import (
	"flag"
	"fmt"
	"pkg/apiregist"

	"api/internal/config"
	"api/internal/handler"
	"api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/edit.yaml", "the config file")

func init() {
	//api网关注册
	if err := apiregist.ApiRegist(); err != nil {
		panic("api注册失败:" + err.Error())
	}
}

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

// fresh edit.go -f etc/edit.yaml

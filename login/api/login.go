package main

import (
	"flag"
	"fmt"

	"api/internal/config"
	"api/internal/handler"
	"api/internal/svc"
	"pkg/apiregist"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/login.yaml", "the config file")

func init() {
	logx.DisableStat()
	logx.SetLevel(1)
	logx.SetUp(logx.LogConf{
		Mode: "file",
		Path: "../logs",
	})
	// api网关注册
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

	// 微服务网关设置了跨域，这里就不用设置跨域了，不然会报错
	// server.Use(middleware.NewCrossDomain().Handler)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

// fresh login.go -f etc/login.yaml

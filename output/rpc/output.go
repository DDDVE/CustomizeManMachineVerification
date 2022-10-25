package main

import (
	"flag"
	"log"

	"rpc/internal/config"
	"rpc/internal/server"
	"rpc/internal/svc"
	"rpc/types/output"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/output.yaml", "the config file")

func init() {
	logx.DisableStat()
	// logx.SetLevel(1)
	// logx.SetUp(logx.LogConf{
	// 	Mode: "file",
	// 	Path: "../logs",
	// })
}

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		output.RegisterOutputServer(grpcServer, server.NewOutputServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	log.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}

//fresh output.go -f etc/output.yaml

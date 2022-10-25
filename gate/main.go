package main

import (
	"gate/handler"
	"gate/initialize"
	"gate/middleware"
	"gate/utils/log"
	"net/http"
)

func init() {}

func main() {
	// 初始化网关相关任务
	initialize.InitGate()

	// 使用全局中间件黑名单，限流,跨域
	middleware.Use(middleware.CrossDomain, middleware.BlackIP, middleware.Limit)

	// api网关注册
	http.HandleFunc("/apiRegist", middleware.Handler(handler.ApiRegist))

	// 对外提供人机验证问题
	http.HandleFunc("/output", middleware.Handler(handler.Output))

	// 登录相关请求
	http.HandleFunc("/login/", middleware.Handler(handler.Login))

	// 使用token校验中间件
	middleware.Use(middleware.TokenCheck)

	// 其余请求进行转发
	http.HandleFunc("/", middleware.Handler(handler.DistributeReq))

	log.Println("服务已启动")
	if err := http.ListenAndServeTLS(":8888", "../pkg/conf/cmmvplat.com_bundle.crt", "../pkg/conf/cmmvplat.com.key", nil); err != nil {
		log.Panic("启动监听失败", err)
		return
	}
}

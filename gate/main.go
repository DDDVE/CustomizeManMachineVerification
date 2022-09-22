package main

import (
	"log"
	"net/http"

	"gate/schedule"
	"gate/utils"

	"gate/handler"
)

func main() {
	// 初始化周期生成平台密钥任务
	go schedule.InitGenerateKeyScheTask()

	// 初始化第一天的密钥
	utils.InitFirstKey()

	// 初始化api网关信息
	handler.InitApiGate()

	// 初始化api网关周期检查任务
	go schedule.InitApiTestScheTask()

	// 初始化黑名单数据
	handler.InitBlackIp()

	// 初始化令牌桶
	utils.InitLimit()

	// 登录
	http.HandleFunc("/login/employee", handler.Login)
	// 注册
	//http.HandleFunc("/regist", handler.Regist)

	// api网关注册
	http.HandleFunc("/apiRegist/*", handler.ApiRegist)

	// 进入拦截器判断
	http.HandleFunc("/", handler.Intercept)
	//fmt.Println(os.Getwd())
	if err := http.ListenAndServe("localhost:8080", nil); err != nil {
		log.Panic("启动监听失败")
		return
	}
}

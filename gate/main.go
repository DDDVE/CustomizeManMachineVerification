package main

import (
	"log"
	"net/http"

	"gate/schedule"
	"gate/utils"

	"gate/handler"
)

func main() {
	// 初始化周期任务
	go schedule.InitGenerateKeyScheTask()

	// 初始化第一天的密钥
	utils.InitFirstKey()
	// todo: 从磁盘读取黑名单数据

	// 登录
	http.HandleFunc("/login", handler.Login)
	// 注册
	//http.HandleFunc("/regist", handler.Regist)

	// 进入拦截器判断
	http.HandleFunc("/", handler.Intercept)

	if err := http.ListenAndServe("localhost:8080", nil); err != nil {
		log.Panic("启动监听失败")
		return
	}
	log.Println("项目启动成功，监听中...")
}

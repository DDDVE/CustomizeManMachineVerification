package main

import (
	"net/http"

	"gate/handler"
	"gate/schedule"
	"gate/utils"
)

func main() {
	// 初始化周期任务
	go schedule.InitGenerateKeyScheTask()

	// 初始化第一天的密钥
	utils.InitFirstKey()

	http.HandleFunc("/", handler.Index)
	http.HandleFunc("/login", handler.Login)
}

package initialize

import (
	"gate/handler"
	"gate/initialize/schedule"
	"gate/middleware"
	"gate/utils/log"
)

func InitGate() {
	//设置日志输出路径和级别
	log.SetOutPath("./conf/log")
	log.SetLevel(log.LEVEL_INFO)

	// 初始化服务降级
	// utils.InitLevel()

	// 初始化周期生成平台密钥任务
	go schedule.InitGenerateKeyScheTask()

	// 初始化api网关信息
	handler.InitApiGate()

	// 初始化api网关周期检查任务
	go schedule.InitApiTestScheTask()

	// 初始化黑名单数据
	middleware.InitBlackIp()

	// 初始化令牌桶
	middleware.InitLimit()
}

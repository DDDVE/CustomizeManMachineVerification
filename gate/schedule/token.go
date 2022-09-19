package schedule

import (
	"gate/utils"
	"log"
	"time"
)

const (
	// 一天的秒数
	SecondsOfDay = 86400
)

// 初始任务是先在第二天零点触发定时任务
// 同时开启周期性任务
func InitGenerateKeyScheTask() {
	// 获取今天的年月日信息
	year, month, day := time.Now().Date()
	location, _ := time.LoadLocation("Asia/Shanghai")
	// 获取项目启动的第二天零点的时间戳
	t := time.Date(year, month, day, 0, 0, 0, 0, location).Unix() + SecondsOfDay
	// 获取当前时间距离第二天早上零点的秒数
	duration := t - time.Now().Unix()
	D := duration * int64(time.Second)
	// 开启定时任务
	log.Println("距离第二天还有：", D)
	// 注意time.Duration是纳秒
	sche := time.NewTimer(time.Duration(D))

	// // 测试定时任务，每五秒打印一次情况
	// test := time.NewTicker(10 * time.Second)
	// var i int64 = 1
	// for range test.C {
	// 	log.Println("还剩", duration-i*10, "秒")
	// 	i++
	// 	if i*10 >= duration {
	// 		break
	// 	}
	// }
	// test.Stop()
	// log.Println("倒计时结束, 马上执行定时任务")

	// 当时间到了
	<-sche.C
	log.Println("定时任务来了!")
	firstGenerateKeyTask()

	// 第一次定时任务过后开启周期任务
	ticker := time.NewTicker(SecondsOfDay * time.Second)
	defer func() {
		log.Println("周期性生成密钥任务停止")
		ticker.Stop()
	}()

	for range ticker.C {
		keyForToday := utils.GetRandomString(utils.LenOfKey)
		log.Println("生成今天的平台随机密钥: ", keyForToday)
		// 将头一天的密钥移动至上一个位置
		utils.MySignedKey[0] = utils.MySignedKey[1]
		// 最新的位置放上最新的密钥
		utils.MySignedKey[1] = keyForToday
	}
}

// 项目初始化时的定时任务
func firstGenerateKeyTask() {
	// 为今天生成随机平台密钥
	keyForToday := utils.GetRandomString(utils.LenOfKey)
	log.Println("生成今天的平台随机密钥: ", keyForToday)
	// 放入密钥
	utils.MySignedKey[1] = keyForToday
}

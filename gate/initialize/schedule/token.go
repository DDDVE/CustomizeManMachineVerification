package schedule

import (
	"gate/utils/log"
	"gate/utils/token"
	"time"
)

const (
	// 一天的秒数
	GENERATE_KEY_PERIOD = 24 //单位h
)

// 初始任务是先在第二天零点触发定时任务
// 同时开启周期性任务
func InitGenerateKeyScheTask() {
	token.GetRandomKey()
	// 获取今天的年月日信息
	year, month, day := time.Now().Date()
	location, _ := time.LoadLocation("Asia/Shanghai")
	// 凌晨3点开始生成密钥,获取项目启动的第二天零点的时间戳
	t := time.Date(year, month, day, 23, 59, 59, 0, location).Unix()
	// 获取当前时间距离第二天早上零点的秒数
	duration := t - time.Now().Unix()
	// 开启定时任务
	log.Debugf("距离第二天还有%v分钟", duration/60)
	// 注意time.Duration是纳秒
	sche := time.NewTimer(time.Duration(duration * int64(time.Second)))

	// 当时间到了
	<-sche.C
	log.Println("生成密钥定时任务开启!")
	token.GetRandomKey()

	// 第一次定时任务过后开启周期任务
	ticker := time.NewTicker(GENERATE_KEY_PERIOD * time.Hour)
	defer func() {
		log.Println("周期性生成密钥任务停止")
		ticker.Stop()
	}()

	for range ticker.C {
		token.GetRandomKey()
	}
}

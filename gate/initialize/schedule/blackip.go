package schedule

import (
	"gate/middleware"
	"gate/utils/log"
	"gate/utils/syncmap"
	"time"
)

func InitBlackIpScheTask() {
	log.Println("开始黑名单定时周期任务")
	ticker := time.NewTicker(middleware.BLACK_PERIOD_TASK * time.Second)
	defer func() {
		log.Println("黑名单周期任务退出")
		ticker.Stop()
	}()
	for range ticker.C {
		if middleware.SecondStatistic.Lenght() != 0 {
			middleware.SecondStatistic = syncmap.NewSyncMap()
		}
	}
}

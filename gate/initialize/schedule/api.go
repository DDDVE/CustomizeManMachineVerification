package schedule

import (
	"context"
	"errors"
	"gate/handler"
	"gate/utils"
	"gate/utils/log"
	"time"
)

/**
定期向各个模块的api网关发送测试信息，调用ping接口，回复pong则ok
*/

//项目启动时初始化周期性检查api网关的任务
func InitApiTestScheTask() {
	log.Println("开始启动周期性检查api网关的任务")
	ticker := time.NewTicker(handler.API_GATE_CHECK_PERIOD * time.Second)
	defer func() {
		log.Println("周期性检查api网关任务退出")
		ticker.Stop()
	}()
	for range ticker.C {
		flag := true
		for i := 0; i < handler.MODULE_COUNT; i++ {
			if handler.ApiData[i].ApiCount > 0 {
				flag = false
				break
			}
		}
		log.Debug("api flag///", flag)
		if flag {
			log.Debug("无任何api网关，已阻塞///")
			<-handler.Button
			log.Debug("有api网关注册，开始周期性检查///")
		}
		ApiTestScheTask()
	}
}

func ApiTestScheTask() {
	for i := 0; i < handler.MODULE_COUNT; i++ {
		if handler.ApiData[i].ApiCount > 0 {
			handler.ApiRWLock.RLock()
			log.Debug("开始发送心跳。。。", handler.ApiData[i].ModuleName)
			for j := 0; j < handler.ApiData[i].ApiCount; j++ {
				go testEachApi(i, handler.ApiData[i].ApiAddrs[j])
			}
			handler.ApiRWLock.RUnlock()
		}
	}
}

func testEachApi(moduleID int, address string) {
	url := handler.REQUEST_URL_PREFIX + address + handler.REQUEST_URL_SUFFIX
	log.Debug("发送心跳的URL：", url)
	err := errors.New("")
	for i := 0; i < handler.RETRY_TIMES; i++ {
		log.Debugf("开始发送第%v次请求，adress:%v", i, address)
		ctx, cancel := context.WithTimeout(context.Background(), handler.REQUEST_TIMEOUT*time.Second)
		err = utils.SendHttpGet(ctx, url)
		cancel()
		if err == nil {
			break
		}
		time.Sleep(handler.API_REQUEST_PERIOD * time.Second)
	}
	if err != nil {
		log.Debug("开始删除此地址，address:", address)
		deleteApiAddr(moduleID, address)
	}
}

func deleteApiAddr(moduleID int, address string) {
	for i := 0; i < handler.ApiData[moduleID].ApiCount; i++ {
		if address == handler.ApiData[moduleID].ApiAddrs[i] {
			handler.ApiRWLock.Lock()
			handler.ApiData[moduleID].ApiAddrs[i], handler.ApiData[moduleID].ApiAddrs[handler.ApiData[moduleID].ApiCount-1] =
				handler.ApiData[moduleID].ApiAddrs[handler.ApiData[moduleID].ApiCount-1], handler.ApiData[moduleID].ApiAddrs[i]
			handler.ApiData[moduleID].ApiCount--
			handler.ApiRWLock.Unlock()
		}
	}
	handler.PersistenceApi()
	log.Printf("addr:%v 已从模板%v中删除", address, handler.ApiData[moduleID].ModuleName)
}

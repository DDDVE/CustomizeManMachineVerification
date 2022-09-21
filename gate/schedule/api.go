package schedule

import (
	"gate/handler"
	"gate/utils"
	"log"
	"sync"
	"time"
)

/**
定期向各个模块的api网关发送测试信息
*/

// 等待检查各个api网关的协程统一处理完毕
var wg = &sync.WaitGroup{}

// 记录哪种api网关的第几个服务器需要解除注册
// api网关种类 -> 下标切片
var ApiCheckoutMap = map[string][]int{}

// 检查api网关周期定为10分钟
const ApiGateCheckPeriod = 600

// 项目启动时初始化周期性检查api网关的任务
func InitApiTestScheTask() {
	log.Println("开始启动周期性检查api网关的任务")
	ticker := time.NewTicker(ApiGateCheckPeriod * time.Second)
	defer func() {
		log.Println("周期性检查api网关任务退出")
		ticker.Stop()
	}()
	for range ticker.C {
		ApiTestScheTask()
	}
}

func ApiTestScheTask() {
	//检查ApiCheckoutMap是否已清零，如果没有则等待
	for len(ApiCheckoutMap) > 0 {
		time.Sleep(time.Second)
	}
	// 先向ApiMap加读锁，拷贝一份内容
	copyApiMap := map[string][]*handler.ApiGate{}
	// 记录有多少个api网关被检查
	apiNum := 0
	handler.ApiMapRWMutex.RLock()
	for k, v := range handler.ApiMap {
		copyApiMap[k] = []*handler.ApiGate{}
		for i := 0; i < len(v); i++ {
			// 并发环境下，可能会出现某个api网关被标记了不可用，但负责清理的协程还没来得及清理
			if v[i].Status != 0 {
				continue
			}
			apiNum++
			copy := &handler.ApiGate{
				Address: v[i].Address,
				Port:    v[i].Port,
				Type:    v[i].Type,
				Index:   v[i].Index,
			}
			copyApiMap[k] = append(copyApiMap[k], copy)
		}
	}
	handler.ApiMapRWMutex.RUnlock()
	// 对于每一个类型的模块切片，起一个协程，检查每一个api网关是否存活
	for _, v := range copyApiMap {
		wg.Add(1)
		go testEachApi(v)
	}
	// 这里需要等待各个协程执行完毕
	wg.Wait()

	// 执行完毕后，把ApiCheckoutMap中的记录在ApiMap中标记为不可用
	i := 0
	handler.ApiMapRWMutex.Lock()
	for k, v := range ApiCheckoutMap {
		for len(ApiCheckoutMap[k]) > 0 {
			i++
			handler.ApiMap[k][v[0]].Status = 1
			ApiCheckoutMap[k] = ApiCheckoutMap[k][1:]
		}
		delete(ApiCheckoutMap, k)
	}
	handler.ApiMapRWMutex.Unlock()
	log.Printf("定时检查api网关任务完成, 共检查了%d个api网关, 有%d个api网关被标记\n", apiNum, i)
}

func testEachApi(toTest []*handler.ApiGate) {
	for i := 0; i < len(toTest); i++ {
		alive := false
		// 拿到ip端口，发送http请求
		apiGate := toTest[i]
		url := "http://" + apiGate.Address + ":" + apiGate.Port + "/test"
		res, err := utils.SendHttpGet(url)
		// 如果发送请求时报错或无响应
		if err != nil || res == "" {
			// 重试三次
			for i := 0; i < 3; i++ {
				res, _ := utils.SendHttpGet(url)
				if res != "" {
					alive = true
					break
				}
			}
			// 如果仍没有响应则认为该api网关宕机
			if !alive {
				// 在ApiCheckoutMap中记录, 所有协程全部检查完后统一标记为不可用
				// 这里不用加锁，因为一个种类的api网关只有一个协程在检查和记录
				if ApiCheckoutMap[apiGate.Type] == nil {
					ApiCheckoutMap[apiGate.Type] = []int{apiGate.Index}
				} else {
					ApiCheckoutMap[apiGate.Type] = append(ApiCheckoutMap[apiGate.Type], apiGate.Index)
				}
				continue
			}
		}
	}
	// 执行完毕后waitGroup减一
	wg.Done()
}

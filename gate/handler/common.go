package handler

import (
	"gate/utils"
	"net/http"
	"strings"
)

/**
此文件用于汇总handler包下的公共代码，避免冗余
*/

// 循环遍历ApiMap，得到一个可用的或者报错
func FindValidApiGate(start int, end int, apiType string) int {
	for i := start; i < len(ApiMap[apiType]); i = (i + 1) % len(ApiMap[apiType]) {
		// 如果遍历到终点还没有说明这种类型没有api网关可用
		if i == end {
			break
		}
		if ApiMap[apiType][i].Status == 0 {
			return i
		}
	}
	return -1
}

// 公共的重定向方法
func CommonRedirct(w http.ResponseWriter, r *http.Request, apiType string) {
	// 通过负载均衡算法得到一个api网关的IP端口
	ip := strings.Split(r.RemoteAddr, ":")[0]
	pos := utils.FindApiGateToRedirect(apiType, ip)

	// 加读锁
	ApiMapRWMutex.RLock()
	pos = pos % len(ApiMap[apiType])
	// 如果该位置的api网关不可用，则往后遍历到一个可用的再转发
	if ApiMap[apiType][pos].Status != 0 {
		pos = FindValidApiGate(pos, pos, apiType)
	}
	// 如果返回-1表示没有可用的api网关
	if pos == -1 {
		utils.WriteData(w, &utils.HttpRes{
			Status: utils.HttpRefuse,
			Data:   nil,
		})
		return
	}
	redirectApiGate := ApiMap[apiType][pos]
	// 解读锁
	ApiMapRWMutex.RUnlock()
	// 重定向
	w.Header().Set("location", redirectApiGate.Address+":"+redirectApiGate.Port)
	w.WriteHeader(http.StatusFound)
}

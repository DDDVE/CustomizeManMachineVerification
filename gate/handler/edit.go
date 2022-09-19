package handler

import (
	"gate/utils"
	"log"
	"net/http"
	"strings"
)

func EditGet(w http.ResponseWriter, r *http.Request, apiType string) {
	log.Printf("进入%s重定向板块, %v方法\n", apiType, r.Method)
	// 通过负载均衡算法得到一个api网关的IP端口
	ip := strings.Split(r.RemoteAddr, ":")[0]
	pos := utils.FindApiGateToRedirect(apiType, ip)
	// 加读锁
	ApiMapRWMutex.RLock()
	pos = pos % len(ApiMap[apiType])
	redirectApiGate := ApiMap[apiType][pos]
	// 解读锁
	ApiMapRWMutex.RUnlock()
	// 重定向
	w.Header().Set("location", redirectApiGate.Address+":"+redirectApiGate.Port)
	w.WriteHeader(http.StatusFound)
}

func EditPost(w http.ResponseWriter, r *http.Request, apiType string) {
	log.Printf("进入%s重定向板块, %v方法\n", apiType, r.Method)
	// 通过负载均衡算法得到一个api网关的IP端口
	ip := strings.Split(r.RemoteAddr, ":")[0]
	pos := utils.FindApiGateToRedirect(apiType, ip)
	// 加读锁
	ApiMapRWMutex.RLock()
	pos = pos % len(ApiMap[apiType])
	redirectApiGate := ApiMap[apiType][pos]
	// 解读锁
	ApiMapRWMutex.RUnlock()
	// 重定向
	w.Header().Set("location", redirectApiGate.Address+":"+redirectApiGate.Port)
	w.WriteHeader(http.StatusFound)
}

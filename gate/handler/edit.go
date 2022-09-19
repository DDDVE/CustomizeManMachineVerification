package handler

import (
	"log"
	"net/http"
)

func EditGet(w http.ResponseWriter, r *http.Request) {
	log.Printf("进入edit重定向板块, %v方法\n", r.Method)
	w.Header().Set("Location", testUrl)
	w.WriteHeader(http.StatusFound)
	// 通过负载均衡算法得到一个api网关的IP端口
}

func EditPost(w http.ResponseWriter, r *http.Request) {
	log.Printf("进入edit重定向板块, %v方法\n", r.Method)
	// 通过负载均衡算法得到一个api网关的IP端口
}

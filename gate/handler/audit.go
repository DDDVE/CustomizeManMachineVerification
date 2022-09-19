package handler

import (
	"log"
	"net/http"
)

func AuditGet(w http.ResponseWriter, r *http.Request) {
	log.Printf("进入audit重定向板块, %v方法\n", r.Method)
	// 通过负载均衡算法得到一个api网关的IP端口
}

func AuditPost(w http.ResponseWriter, r *http.Request) {
	log.Printf("进入audit重定向板块, %v方法\n", r.Method)
	// 通过负载均衡算法得到一个api网关的IP端口
}

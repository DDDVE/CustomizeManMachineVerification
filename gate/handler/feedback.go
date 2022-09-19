package handler

import (
	"log"
	"net/http"
)

func FeedbackGet(w http.ResponseWriter, r *http.Request) {
	log.Printf("进入feedback重定向板块, %v方法\n", r.Method)
	// 通过负载均衡算法得到一个api网关的IP端口
}

func FeedbackPost(w http.ResponseWriter, r *http.Request) {
	log.Printf("进入feedback重定向板块, %v方法\n", r.Method)
	// 通过负载均衡算法得到一个api网关的IP端口
}

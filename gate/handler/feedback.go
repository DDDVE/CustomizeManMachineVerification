package handler

import (
	"log"
	"net/http"
)

func FeedbackGet(w http.ResponseWriter, r *http.Request, apiType string) {
	log.Printf("进入%s重定向板块, %v方法\n", apiType, r.Method)
	CommonRedirct(w, r, apiType)
}

func FeedbackPost(w http.ResponseWriter, r *http.Request, apiType string) {
	log.Printf("进入%s重定向板块, %v方法\n", apiType, r.Method)
	CommonRedirct(w, r, apiType)
}

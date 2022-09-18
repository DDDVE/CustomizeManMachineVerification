package handler

import (
	"log"
	"net/http"
)

func EditGet(w http.ResponseWriter, r *http.Request) {
	log.Printf("进入edit重定向板块, %v方法", r.Method)
	w.Header().Set("Location", testUrl)
	w.WriteHeader(http.StatusFound)
}

func EditPost(w http.ResponseWriter, r *http.Request) {

}

package handler

import (
	"gate/utils"
	"log"
	"net/http"
	"strings"
)

type LoginData struct {
	Token string `json:"token"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	log.Println("进入登录模块")
	//TODO: 重定向至登录模块

	// 以下暂时保留方便测试
	// 根据最新的密钥生成token返回给用户
	auth := strings.Split(r.Header.Get("Authorization"), "@==@")
	if len(auth) < 2 {
		log.Println("请求头不合法")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userMobile := auth[1]
	userToken, err := utils.GenerateTokenString(userMobile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	utils.WriteData(w, &utils.HttpRes{
		Status: utils.HttpSucceed,
		Data: LoginData{
			Token: userToken,
		},
	})

}

func LoginGet(w http.ResponseWriter, r *http.Request, apiType string) {
	log.Printf("进入%s重定向板块, %v方法\n", apiType, r.Method)
	CommonRedirct(w, r, apiType)
}

func LoginPost(w http.ResponseWriter, r *http.Request, apiType string) {
	log.Printf("进入%s重定向板块, %v方法\n", apiType, r.Method)
	CommonRedirct(w, r, apiType)
}

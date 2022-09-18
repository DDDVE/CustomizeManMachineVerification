package handler

import (
	"gate/utils"
	"log"
	"net/http"
	"strings"
)

func Regist(w http.ResponseWriter, r *http.Request) {
	log.Println("进入注册模块")
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

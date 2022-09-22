package handler

import (
	"context"
	"encoding/json"
	"gate/utils"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type LoginData struct {
	Token string `json:"token"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	log.Println("进入登录模块")
	//判断请求的ip是否在黑名单
	ip := strings.Split(r.RemoteAddr, ":")[0]
	if CheckBlackIp(ip) {
		log.Printf("该地址%s在黑名单中, 已拦截\n", ip)
		return
	}

	// 限流
	ctx, _ := context.WithTimeout(context.Background(), utils.LinkTimeOut*time.Second)
	err := utils.Limiter.Wait(ctx)
	if err != nil {
		log.Printf("主机%s因限流, 获取令牌失败: %+v", ip, err)
		utils.WriteData(w, &utils.HttpRes{
			Status: utils.HttpRefuse,
			Data:   nil,
		})
		return
	}

	// 查找可用的login类型api网关
	pos := utils.FindApiGateToRedirect(utils.TypeOfApiLogin, ip)
	ApiMapRWMutex.RLock()
	pos = pos % len(ApiMap[utils.TypeOfApiLogin])
	if ApiMap[utils.TypeOfApiLogin][pos].Status != 0 {
		pos = FindValidApiGate(pos, pos, utils.TypeOfApiLogin)
	}
	if pos == -1 {
		utils.WriteData(w, &utils.HttpRes{
			Status: utils.HttpRefuse,
			Data:   nil,
		})
		return
	}
	r.Host = ApiMap[utils.TypeOfApiLogin][pos].Address + ":" + ApiMap[utils.TypeOfApiLogin][pos].Port
	ApiMapRWMutex.RUnlock()
	// 发送post请求
	context, err := utils.RetransmissionPost(r)
	// 接收响应数据并判断
	if err != nil {
		utils.WriteData(w, &utils.HttpRes{
			Status: utils.HttpRefuse,
			Data:   nil,
		})
		return
	}
	res := map[string]string{}
	err = json.Unmarshal(context, &res)
	if err != nil {
		log.Println("解析转发响应体报错: ", err)
		utils.WriteData(w, &utils.HttpRes{
			Status: utils.HttpRefuse,
			Data:   nil,
		})
		return
	}
	// 如果响应体为空则认为登陆失败
	if res["status"] != strconv.Itoa(utils.HttpSucceed) {
		utils.WriteData(w, &utils.HttpRes{
			Status: utils.HttpRefuse,
			Data:   nil,
		})
		return
	}
	// 根据最新的密钥生成token返回给用户
	auth := strings.Split(r.Header.Get("Authorization"), "@==@")
	if len(auth) > 1 {
		log.Printf("主机%s的请求头不合法\n", ip)
		utils.WriteData(w, &utils.HttpRes{
			Status: utils.HttpTokenCheckFalse,
			Data:   nil,
		})
		return
	}
	userMobile := auth[0]
	userToken, err := utils.GenerateTokenString(userMobile)
	if err != nil {
		utils.WriteData(w, &utils.HttpRes{
			Status: utils.HttpRefuse,
			Data:   nil,
		})
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

// 转发用户的登录请求并判断是否登录成功
func SendLoginRequest() {

}

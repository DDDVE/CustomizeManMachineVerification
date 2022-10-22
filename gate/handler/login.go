package handler

import (
	"gate/utils"
	"gate/utils/log"
	"gate/utils/token"
	"net/http"
	"strings"
)

type LoginData struct {
	Token string `json:"token"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	requestID := r.Header.Get("requestID")
	userIP := r.Header.Get("userIP")
	log.Debugf("requestID:%v/// userIP:%v///", requestID, userIP)
	//不合法的URL直接在降级中间件过滤掉，所以这里暂不做细致校验
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		utils.RespFormat(w, utils.REQUEST_METHOD_ERROR, nil)
	}
	//如果是登录提交表单请求成功则签发token
	path := strings.Split(r.URL.Path, "/")
	p := ""
	if len(path) >= 3 {
		p = path[2]
	}
	if p == "employee" {
		r.ParseForm()
		log.Debug("r.PostForm[\"mobileNum\"][0]:////", r.Form["mobileNum"][0])
		token := token.GenerateToken(r.Form["mobileNum"][0])
		log.Debug("生成的token:///", token)
		w.Header().Add("Access-Control-Expose-Headers", "Authorization, Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		w.Header().Add("Authorization", token)
	}
	//解析请求的模块名
	moduleID := getModuleID(r.URL.Path)
	if moduleID == -1 {
		utils.RespFormat(w, utils.REQUEST_PATH_ERROR, nil)
	}

	//通过对userAddr哈希取余方式负载均衡
	id := loadBalance(r.RemoteAddr, moduleID)

	//转发请求
	resp, err := forwardReq(r, ApiData[moduleID].ApiAddrs[id])
	if err != nil {
		utils.RespUnknownErr(w, err)
	}

	w.Write(resp)

	// 查找可用的login类型api网关
	// pos := utils.FindApiGateToRedirect(utils.TypeOfApiLogin, ip)
	// ApiMapRWMutex.RLock()
	// pos = pos % len(ApiMap[utils.TypeOfApiLogin])
	// if ApiMap[utils.TypeOfApiLogin][pos].Status != 0 {
	// 	pos = FindValidApiGate(pos, pos, utils.TypeOfApiLogin)
	// }
	// if pos == -1 {
	// 	utils.WriteData(w, &utils.HttpRes{
	// 		Status: utils.HttpRefuse,
	// 		Data:   nil,
	// 	})
	// 	return
	// }
	// r.Host = ApiMap[utils.TypeOfApiLogin][pos].Address + ":" + ApiMap[utils.TypeOfApiLogin][pos].Port
	// ApiMapRWMutex.RUnlock()
	// // 发送post请求
	// context, err := utils.RetransmissionPost(r)
	// // 接收响应数据并判断
	// if err != nil {
	// 	utils.WriteData(w, &utils.HttpRes{
	// 		Status: utils.HttpRefuse,
	// 		Data:   nil,
	// 	})
	// 	return
	// }
	// res := map[string]string{}
	// err = json.Unmarshal(context, &res)
	// if err != nil {
	// 	log.Println("解析转发响应体报错: ", err)
	// 	utils.WriteData(w, &utils.HttpRes{
	// 		Status: utils.HttpRefuse,
	// 		Data:   nil,
	// 	})
	// 	return
	// }
	// // 如果响应体为空则认为登陆失败
	// if res["status"] != strconv.Itoa(utils.HttpSucceed) {
	// 	utils.WriteData(w, &utils.HttpRes{
	// 		Status: utils.HttpRefuse,
	// 		Data:   nil,
	// 	})
	// 	return
	// }
	// // 根据最新的密钥生成token返回给用户
	// auth := strings.Split(r.Header.Get("Authorization"), "@==@")
	// if len(auth) > 1 {
	// 	log.Printf("主机%s的请求头不合法\n", ip)
	// 	utils.WriteData(w, &utils.HttpRes{
	// 		Status: utils.HttpTokenCheckFalse,
	// 		Data:   nil,
	// 	})
	// 	return
	// }
	// userMobile := auth[0]
	// userToken, err := utils.GenerateTokenString(userMobile)
	// if err != nil {
	// 	utils.WriteData(w, &utils.HttpRes{
	// 		Status: utils.HttpRefuse,
	// 		Data:   nil,
	// 	})
	// 	return
	// }
	// utils.WriteData(w, &utils.HttpRes{
	// 	Status: utils.HttpSucceed,
	// 	Data: LoginData{
	// 		Token: userToken,
	// 	},
	// })
}

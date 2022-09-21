package handler

import (
	"fmt"
	"gate/utils"
	"log"
	"net/http"
	"os"
	"strings"
)

func Intercept(w http.ResponseWriter, r *http.Request) {
	log.Println("进入拦截器")
	//判断请求的ip是否在黑名单
	ip := strings.Split(r.RemoteAddr, ":")[0]
	if CheckBlackIp(ip) {
		log.Printf("该地址%s在黑名单中, 已拦截\n", ip)
		return
	}
	// 验证token
	autho := strings.Split(r.Header.Get("Authorization"), "@==@")
	if len(autho) < 2 {
		log.Println("请求头不合法")
		utils.WriteData(w, &utils.HttpRes{
			Status: utils.HttpTokenCheckFalse,
			Data:   nil,
		})
		// 这里打印普通信息后会往下继续执行，所以需要return
		return
	}
	userToken := autho[0]
	userMobile := autho[1]
	if !utils.CheckTokenString(userToken, userMobile) {
		log.Println("请求头不合法")
		utils.WriteData(w, &utils.HttpRes{
			Status: utils.HttpTokenCheckFalse,
			Data:   nil,
		})
		return
	}
	log.Println("用户token校验通过")
	// token验证通过后对url做匹配，并分配给不同的handler
	// 首先对请求类型做判断
	switch r.Method {
	case "GET":
		HandleGet(w, r)
	case "POST":
		HandlePost(w, r)
	default:
		utils.WriteData(w, &utils.HttpRes{
			Status: utils.HttpMethodCheckFalse,
			Data:   nil,
		})
	}
}

const (
	// 测试用的重定向地址
	testUrl = "http://www.baidu.com"
)

func HandleGet(w http.ResponseWriter, r *http.Request) {
	log.Println("进入get方法处理")
	// 获取第一级的请求目录，如edit，audit，feedback
	path := strings.Split(r.URL.Path, "/")
	path = path[1:]
	if path[len(path)-1] == "" {
		path = path[:len(path)-1]
	}
	// 对父级目录做判断
	switch path[0] {
	case utils.TypeOfApiEdit:
		EditGet(w, r, path[0])
	case utils.TypeOfApiAudit:
		AuditPost(w, r, path[0])
	case utils.TypeOfApiFeedback:
		FeedbackGet(w, r, path[0])
	case utils.TypeOfApiLogin:
		LoginGet(w, r, path[0])
	default:
		utils.WriteData(w, &utils.HttpRes{
			Status: utils.HttpUrlCheckFalse,
			Data:   nil,
		})
	}
}

func HandlePost(w http.ResponseWriter, r *http.Request) {
	log.Println("进入post方法处理")
	// 获取第一级的请求目录，如edit，audit，feedback
	path := strings.Split(r.URL.Path, "/")
	path = path[1:]
	if path[len(path)-1] == "" {
		path = path[:len(path)-1]
	}
	// 对父级目录做判断
	switch path[0] {
	case utils.TypeOfApiEdit:
		EditPost(w, r, path[0])
	case utils.TypeOfApiAudit:
		AuditPost(w, r, path[0])
	case utils.TypeOfApiFeedback:
		FeedbackPost(w, r, path[0])
	case utils.TypeOfApiLogin:
		LoginPost(w, r, path[0])
	default:
		utils.WriteData(w, &utils.HttpRes{
			Status: utils.HttpUrlCheckFalse,
			Data:   nil,
		})
	}
}

// 黑名单
var BlackIpMap = map[string]bool{}

func InitBlackIp() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Panicln("获取工作目录报错: ", err, " 终止程序")
	}
	context, err := utils.ReadFile(pwd + "\\blackIp.txt")
	if err != nil {
		log.Panicln("读取黑名单文件报错: ", err, " 终止程序")
	}
	ips := strings.Split(context, "\n")
	for i := 0; i < len(ips); i++ {
		BlackIpMap[ips[i]] = true
	}
	log.Println("读取黑名单完成")
	for k, _ := range BlackIpMap {
		fmt.Println(k)
	}
}

func CheckBlackIp(ip string) bool {
	return BlackIpMap[ip]
}

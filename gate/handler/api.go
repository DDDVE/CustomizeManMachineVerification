package handler

import (
	"gate/utils"
	"log"
	"net/http"
	"strings"
	"sync"
)

/**
关于api网关注册
*/

// api网关对象
// 注册时才创建
type ApiGate struct {
	Address string `json:"address"`
	// 请求来自于对方哪个端口
	Port string `json:"port"`

	// 类型，指该api网关服务于哪个类型的模块
	Type string `json:"type"`

	// 状态，表示可用或不可用
	// 之所以引入状态是因为每次宕机都需要将该api网关删除，再将后面的移上来太麻烦
	// 0=可用，1=不可用
	Status int `json:"status"`

	// 该api网关在该类的第几个
	Index int `json:"index"`
}

// api网关注册时的响应信息
type ApiRegistRes struct {
	// 加密后的字符串
	Token string `json:"token"`
}

var (
	// 类型名称 -> api网关集群信息
	ApiMap = map[string][]*ApiGate{}

	// api网关种类及对应的公钥
	ApiToPublicKey = map[string]string{}
	// api网关种类及对应的私钥
	ApiToPrivateKey = map[string]string{}

	// api网关地址：address:8080
	// 随机字符串到api网关地址的映射
	RandomStringToApi = map[string]string{}
	// api网关地址到随机字符串的映射
	ApiToRandomString = map[string]string{}
)

// ApiMap在被读的时候可能有新的注册或者宕机发生
// 整体属于读多写少，故采用读写锁
var (
	ApiMapRWMutex            = new(sync.RWMutex)
	RandomStringToApiRWMutex = new(sync.RWMutex)
	ApiToRandomStringRWMutex = new(sync.RWMutex)
)

// 初始化保存api网关的切片
func InitApiGate() {
	for i := 0; i < len(utils.ApiGateSlice); i++ {
		ApiMap[utils.ApiGateSlice[i]] = []*ApiGate{}
	}

	// 初始化各个api网关类型对应的公钥私钥
	//TODO: 从文件中读取
}

func ApiRegist(w http.ResponseWriter, r *http.Request) {
	// 首先得到请求的地址
	// 参数校验
	// 获取api网关的端口号
	log.Println("进入网关注册模块")
	remotePort := r.Header.Get("port")
	if remotePort == "" {
		utils.WriteData(w, &utils.HttpRes{
			Status: utils.HttpParamCheckFalse,
			Data:   nil,
		})
		return
	}
	log.Println("获取api网关IP地址")
	remoteAddr := strings.Split(r.RemoteAddr, ":")
	// 获取该api网关ip
	remoteIp := remoteAddr[0]
	apiAddres := remoteIp + ":" + remotePort
	log.Println("开始解析路由")
	path := strings.Split(r.URL.Path, "/")
	path = path[1:]
	if path[len(path)-1] == "" {
		path = path[:len(path)-1]
	}
	// path此时只能形如 /apiRegist/login等
	if len(path) != 2 {
		log.Printf("网关%s的路由不正确\n", apiAddres)
		utils.WriteData(w, &utils.HttpRes{
			Status: utils.HttpUrlCheckFalse,
			Data:   nil,
		})
		return
	}
	// 获取api网关种类
	remoteType := path[1]

	// 判断是否在api网关地址到随机字符串的映射里面
	// 如果这个地址是第一次发来请求
	ApiToRandomStringRWMutex.RLock()
	if _, ok := ApiToRandomString[apiAddres]; !ok {
		ApiToRandomStringRWMutex.RUnlock()
		log.Printf("这是网关%s的第一次注册请求\n", apiAddres)
		// 生成随机字符串
		randomString := utils.GetRandomString(utils.LenOfKey)
		// 随机字符串加密
		encodedString, err := utils.Encrypt(randomString, ApiToPublicKey[remoteType])
		if err != nil {
			utils.WriteData(w, &utils.HttpRes{
				Status: utils.HttpApiRegistFalse,
				Data:   nil,
			})
			return
		}
		// 保存映射关系
		ApiToRandomStringRWMutex.Lock()
		ApiToRandomString[apiAddres] = randomString
		ApiToRandomStringRWMutex.Unlock()

		RandomStringToApiRWMutex.Lock()
		RandomStringToApi[randomString] = apiAddres
		RandomStringToApiRWMutex.Unlock()

		// 把密文响应给该地址
		utils.WriteData(w, &utils.HttpRes{
			Status: utils.HttpSucceed,
			Data: ApiRegistRes{
				Token: string(encodedString),
			},
		})
		log.Printf("已将密文响应给网关%s\n", apiAddres)
		return
	}
	ApiToRandomStringRWMutex.RUnlock()
	log.Printf("这不是网关%s的第一次注册请求\n", apiAddres)
	// 如果这个地址不是第一次发来请求
	// 检查字符串原文是否匹配

	// 获取api网关发来的解密后的字符串
	// 如果不匹配，响应错误信息
	plainRandomString := r.Header.Get("token")
	if plainRandomString == "" {
		log.Printf("网关%s缺少解密后的字符串\n", apiAddres)
		utils.WriteData(w, &utils.HttpRes{
			Status: utils.HttpTokenCheckFalse,
			Data:   nil,
		})
		return
	}

	ApiToRandomStringRWMutex.RLock()
	if plainRandomString != ApiToRandomString[apiAddres] {
		ApiToRandomStringRWMutex.RUnlock()
		log.Printf("网关%s解密后的字符串错误\n", apiAddres)
		utils.WriteData(w, &utils.HttpRes{
			Status: utils.HttpTokenCheckFalse,
			Data:   nil,
		})
		return
	}
	ApiToRandomStringRWMutex.RUnlock()

	RandomStringToApiRWMutex.RLock()
	if _, ok := RandomStringToApi[plainRandomString]; !ok {
		RandomStringToApiRWMutex.RUnlock()
		log.Printf("找不到网关%s解密后的字符串对应的IP地址\n", apiAddres)
		utils.WriteData(w, &utils.HttpRes{
			Status: utils.HttpTokenCheckFalse,
			Data:   nil,
		})
		return
	}
	if RandomStringToApi[plainRandomString] != apiAddres {
		RandomStringToApiRWMutex.RUnlock()
		log.Printf("网关%s解密后的字符串与记录不匹配\n", apiAddres)
		utils.WriteData(w, &utils.HttpRes{
			Status: utils.HttpTokenCheckFalse,
			Data:   nil,
		})
		return
	}
	RandomStringToApiRWMutex.RUnlock()

	log.Printf("网关%s返回的字符串验证通过, 开始注册\n", apiAddres)
	//判断ApiMap是否已存在，如果存在直接改状态
	ApiMapRWMutex.Lock()
	for i := 0; i < len(ApiMap[remoteType]); i++ {
		gate := ApiMap[remoteType][i]
		// 如果该类型下ip和端口都相同则认为是同一个api网关
		if remoteIp == gate.Address && remoteIp == gate.Port {
			// 将状态重新改回可用
			ApiMap[remoteType][i].Status = 0
			ApiMapRWMutex.Unlock()
			log.Printf("网关%s注册成功\n", apiAddres)
			utils.WriteData(w, &utils.HttpRes{
				Status: utils.HttpSucceed,
				Data:   nil,
			})
			return
		}
	}
	// 生成ApiGate信息并注册到ApiMap中
	apiGate := &ApiGate{
		Address: remoteIp,
		Port:    remotePort,
		Type:    remoteType,
		Index:   len(ApiMap[remoteType]),
	}
	ApiMap[remoteType] = append(ApiMap[remoteType], apiGate)
	ApiMapRWMutex.Unlock()

	log.Printf("网关%s注册成功\n", apiAddres)
	utils.WriteData(w, &utils.HttpRes{
		Status: utils.HttpSucceed,
		Data:   nil,
	})
}

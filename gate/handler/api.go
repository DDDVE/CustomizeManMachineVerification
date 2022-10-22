package handler

import (
	"bufio"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"gate/utils"
	"gate/utils/log"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

/**
关于api网关注册
*/
const (
	//api
	OUTPUT_MODULE_ID   = 0
	LOGIN_MODULE_ID    = 1
	EDIT_MODULE_ID     = 2
	AUDIT_MODULE_ID    = 3
	FEEDBACK_MODULE_ID = 4
	MODULE_COUNT       = 5

	//api注册
	RETRY_TIMES                 = 3
	REGIST_REQUEST_TIMEOUT      = 3
	DIGITAL_SIGNATURE_CONNECTOR = "@==@"
	API_PUBLIC_KEY_FILE         = "./conf/apipublic.pem"

	//api持久化
	API_DATA_FILE_1 = "./conf/apidata1.log"
	API_DATA_FILE_2 = "./conf/apidata2.log"

	//api定时任务
	REQUEST_TIMEOUT       = 3
	API_GATE_CHECK_PERIOD = 10
	REQUEST_URL_PREFIX    = "http://"
	REQUEST_URL_SUFFIX    = "/ping"
)

var (
	Button             chan struct{} //控制定时任务的开关
	ApiRWLock          sync.RWMutex  //
	ApiPersistenceLock sync.Mutex
	ApiDataFiles       = [2]string{API_DATA_FILE_1, API_DATA_FILE_2}
	apiPublicKeyFile   = API_PUBLIC_KEY_FILE
)

//暂时赋值，最好直接在文件读取
var ApiData = [MODULE_COUNT]Module{
	{ModuleID: 0, ApiCount: 0, ModuleName: "output", ApiAddrs: []string{}, modulePublicKeyPath: apiPublicKeyFile},
	{ModuleID: 1, ApiCount: 0, ModuleName: "login", ApiAddrs: []string{}, modulePublicKeyPath: apiPublicKeyFile},
	{ModuleID: 2, ApiCount: 0, ModuleName: "edit", ApiAddrs: []string{}, modulePublicKeyPath: apiPublicKeyFile},
	{ModuleID: 3, ApiCount: 0, ModuleName: "audit", ApiAddrs: []string{}, modulePublicKeyPath: apiPublicKeyFile},
	{ModuleID: 4, ApiCount: 0, ModuleName: "feedback", ApiAddrs: []string{}, modulePublicKeyPath: apiPublicKeyFile},
}

type Module struct {
	ModuleID            int
	ApiCount            int
	ModuleName          string
	ApiAddrs            []string
	modulePublicKeyPath string
}

//初始化时加载绝对路径
// func init() {
// 	PWD, err := os.Getwd()
// 	if err != nil {
// 		log.Error("获取工作目录报错: ", err)
// 	} else {
// 		ApiDataFiles[0] = filepath.Join(PWD, ApiDataFiles[0])
// 		ApiDataFiles[1] = filepath.Join(PWD, ApiDataFiles[1])
// 		apiPublicKeyFile = filepath.Join(PWD, apiPublicKeyFile)
// 	}
// }

//初始化api
func InitApiGate() {
	//=============从持久化文件加载api网关信息=============
	f, err := os.Open(ApiDataFiles[0])
	if err != nil {
		f, err = os.Open(ApiDataFiles[1])
		if err != nil {
			return
		}
	}
	defer f.Close()

	rd := bufio.NewReader(f)
	data, err := rd.ReadBytes('\n')
	if err != nil && err != io.EOF {
		return
	} else {
		if err = json.Unmarshal(data, &ApiData); err != nil {
			return
		}
	}

}

/*api网关注册，get请求入参：
address:注册API的IP+port
moduleID:模块ID
timestamp:时间戳
ciphertext:密文，通过ECC加密上面三个字段的hash
*/
func ApiRegist(w http.ResponseWriter, r *http.Request) {

	//解析url参数
	if err := r.ParseForm(); err != nil {
		w.Write([]byte("请求URL错误: " + r.RequestURI))
		return
	}

	// 参数校验
	id, err := strconv.Atoi(r.FormValue("moduleID"))
	if err != nil || id < 0 || id > MODULE_COUNT {
		w.Write([]byte("请求参数moduleID错误"))
		return
	}
	addr := r.FormValue("address")
	ip := strings.Split(r.RemoteAddr, ":")[0]
	if ip != strings.Split(addr, ":")[0] {
		w.Write([]byte("请求参数address与实际地址不一致"))
		return
	}

	ApiRWLock.RLock()
	for i := 0; i < ApiData[id].ApiCount; i++ {
		if addr == ApiData[id].ApiAddrs[i] {
			w.Write([]byte("此api网关已注册"))
			return
		}
	}
	ApiRWLock.RUnlock()

	stamp, err := strconv.ParseInt(r.FormValue("timestamp"), 10, 64)
	if err != nil {
		w.Write([]byte("请求参数timestamp错误"))
		return
	}
	if time.Since(time.Unix(0, stamp)) > REGIST_REQUEST_TIMEOUT*time.Second {
		w.Write([]byte("请求超时"))
		return
	}
	ciphertext := r.FormValue("ciphertext")
	if len(strings.Split(ciphertext, DIGITAL_SIGNATURE_CONNECTOR)) != 2 {
		w.Write([]byte("请求参数ciphertext格式错误"))
		return
	}

	//验证数字签名
	msg := fmt.Sprintf("%v%v%v", id, addr, stamp)
	hash := sha256.Sum256([]byte(msg)) //生成摘要
	signs := strings.Split(ciphertext, DIGITAL_SIGNATURE_CONNECTOR)
	if err := utils.VerifySignECC(hash[:], []byte(signs[0]), []byte(signs[1]), ApiData[id].modulePublicKeyPath); err != nil {
		w.Write([]byte("请求参数ciphertext错误"))
		return
	}

	//通过校验，开始注册
	ApiRWLock.Lock()
	if ApiData[id].ApiCount == len(ApiData[id].ApiAddrs) {
		ApiData[id].ApiAddrs = append(ApiData[id].ApiAddrs, addr)
	} else {
		ApiData[id].ApiAddrs[ApiData[id].ApiCount] = addr
	}
	ApiData[id].ApiCount++
	ApiRWLock.Unlock()

	w.Write([]byte("ok"))

	//注册成功开启定时任务
	select {
	case Button <- struct{}{}:
	default:
	}

	//注册成功后进行持久化
	PersistenceApi()
}

//api网关数据持久化
func PersistenceApi() {
	fileName := ""
	if _, err := os.Stat(ApiDataFiles[0]); err != nil {
		fileName = ApiDataFiles[0]
	}
	if _, err := os.Stat(ApiDataFiles[1]); err != nil {
		fileName = ApiDataFiles[1]
	}
	if fileName == "" {
		os.Remove(ApiDataFiles[0])
		fileName = ApiDataFiles[0]
	}

	ApiPersistenceLock.Lock()
	defer ApiPersistenceLock.Unlock()

	os.Create(fileName)
	f, err := os.OpenFile(fileName, os.O_RDWR, 0666)
	if err != nil {
		log.Warn("打开api持久化文件错误")
	}
	ApiRWLock.RLock()
	api, err := json.Marshal(ApiData)
	ApiRWLock.RUnlock()
	if err != nil {
		log.Warn("json序列化错误,api网关持久化失败")
	}
	api = append(api, "\n"...)
	f.Write(api)
	log.Debug("api持久化成功！")
	if fileName != ApiDataFiles[0] {
		os.Remove(ApiDataFiles[0])
	} else {
		os.Remove(ApiDataFiles[1])
	}
}

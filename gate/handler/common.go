package handler

import (
	"crypto/md5"
	"encoding/binary"
	"gate/utils/log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

/**
此文件用于汇总handler包下的公共代码，避免冗余
*/

//转发请求的客户端
var client = &http.Client{
	Timeout: 3 * time.Second,
}

// 转发请求
func forwardReq(w http.ResponseWriter, r *http.Request, addr string) error {
	// 组装URL
	log.Debug("addr//////////////////", addr)
	remote, err := url.Parse("https://" + addr)
	log.Debugf("URL:////////- %v -////////", remote.String())
	if err != nil {
		log.Println("设置请求URL出错: ", err)
		return err
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ServeHTTP(w, r)
	return nil
	// r.URL = remote
	// r.Host = addr
	// // r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// req, err := client.Do(r)
	// if err != nil {
	// 	log.Println("转发请求报错: ", err)
	// 	return nil, err
	// }
	// defer req.Body.Close()
	// if req.StatusCode/100 != 2 {
	// 	log.Println("转发请求错误码: ", req.StatusCode)
	// 	return nil, errors.New("转发请求失败，错误码: " + req.Status)
	// }
	// body, err := ioutil.ReadAll(req.Body)
	// if err != nil {
	// 	log.Println("读取转发响应报错: ", err)
	// 	return nil, err
	// }
	// return body, nil
}

// 解析请求的模块名
func getModuleID(path string) int {
	ps := strings.Split(path, "/")
	p := ps[1:]
	if p[len(p)-1] == "" {
		p = p[:len(p)-1]
	}
	moduleID := -1
	for i := 0; i < MODULE_COUNT; i++ {
		if ApiData[i].ModuleName == p[0] {
			moduleID = i
			break
		}
	}
	return moduleID
}

//返回哈希取余后对应的api网关地址
func loadBalance(addr string, moduleID int) int {
	h := md5.Sum([]byte(addr))
	hash := [8]byte{}
	for i := 0; i < 8; i++ {
		hash[i] = h[i] ^ h[len(h)-1-i]
	}
	hashNum := binary.LittleEndian.Uint64(hash[:])
	id := hashNum % uint64(ApiData[moduleID].ApiCount)
	return int(id)
}

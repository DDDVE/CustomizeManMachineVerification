package handler

import (
	"crypto/md5"
	"encoding/binary"
	"errors"
	"gate/utils/log"
	"io/ioutil"
	"net/http"
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
func forwardReq(req *http.Request, addr string) (resp []byte, err error) {
	// 组装URL
	u, err := url.Parse("http://" + addr + req.URL.String())
	log.Debugf("URL:///-%v-///", u.String())
	if err != nil {
		log.Println("设置请求URL出错: ", err)
		return nil, err
	}
	req.URL = u
	req.Host = addr
	// r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r, err := client.Do(req)
	if err != nil {
		log.Println("转发请求报错: ", err)
		return nil, err
	}
	defer r.Body.Close()
	if r.StatusCode/100 != 2 {
		log.Println("转发请求错误码: ", r.StatusCode)
		return nil, errors.New("转发请求失败，错误码: " + r.Status)
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("读取转发响应报错: ", err)
		return nil, err
	}
	return body, nil
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

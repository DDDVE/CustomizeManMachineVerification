package handler

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"gate/utils"
	"gate/utils/log"
	"gate/utils/token"
	"io/ioutil"
	"net/http"
	"unsafe"
)

type (
	employeeReq struct {
		MobileNum string `json:"mobileNum"`
		MsgCode   string `json:"msgCode"`
	}

	apiEmployeeResp struct {
		Status int         `json:"status"`
		Msg    string      `json:"msg"`
		Data   interface{} `json:"data"`
	}
)

// employee请求需要知道响应结果从而签发token,所以不用统一的反向代理方式转发
func Employee(w http.ResponseWriter, r *http.Request, addr string) {
	buf := make([]byte, 1024)
	n, err := r.Body.Read(buf)
	if err != nil && n == 0 {
		utils.RespUnknownErr(w, err)
		return
	}
	r.Body.Close()

	reader := bytes.NewReader(buf[:n])
	url := "https://" + addr + r.URL.String()
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		utils.RespUnknownErr(w, err)
		return
	}
	defer request.Body.Close() //程序在使用完回复后必须关闭回复的主体

	//必须设定该参数,POST参数才能正常提交，意思是以json串提交数据
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	request.Header.Set("Authorization", r.Header.Get("Authorization"))

	//发送请求
	resp, err := client.Do(request)
	if err != nil {
		utils.RespUnknownErr(w, err)
		return
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.RespUnknownErr(w, err)
		return
	}
	//byte数组直接转成string，优化内存
	log.Debug("转发employee请求的响应结果///", *(*string)(unsafe.Pointer(&respBytes)))

	ar := &apiEmployeeResp{}
	if err := json.Unmarshal(respBytes, ar); err != nil {
		utils.RespUnknownErr(w, err)
		return
	}

	if ar.Status == 0 {
		er := &employeeReq{}
		if err := json.Unmarshal(buf[:n], er); err != nil {
			utils.RespUnknownErr(w, err)
			return
		}
		log.Debugf("employeeReq//////////////%+v", er)

		mobile, err := base64.StdEncoding.DecodeString(er.MobileNum)
		if err != nil {
			utils.RespUnknownErr(w, err)
			return
		}
		globalToken := token.GenerateToken(string(mobile))
		log.Debug("生成的globalToken:///", globalToken)

		w.Header().Add("Authorization", globalToken)

		utils.RespFormat(w, utils.SUCCESS, nil)
	} else {
		log.Debug("ar.Msg:///", ar.Msg)
		utils.RespUnknownErr(w, errors.New(ar.Msg))
		return
	}
}

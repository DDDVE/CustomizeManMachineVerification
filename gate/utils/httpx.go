package utils

import (
	"context"
	"encoding/json"
	"errors"
	"gate/utils/log"
	"io/ioutil"
	"net/http"
)

type respBody struct {
	Status int         `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data,omitempty"`
}

// http返回数据
func RespFormat(w http.ResponseWriter, status int, data interface{}) {
	var body respBody
	body.Status = status
	if status == 0 {
		body.Msg = "success"
		body.Data = data
	} else {
		msg := GlobalError[status]
		body.Msg = msg
	}
	log.Debugf("%+v:///// ", body)
	b, err := json.Marshal(body)
	if err != nil {
		log.Println("将响应体转换为json字符串失败: ", err)
		// 服务器内部错误
		w.WriteHeader(http.StatusInternalServerError)
	}

	// 将数据写到响应体
	w.Write([]byte(b))
	log.Debug("b://// ", string(b))
	return
}

// http返回未知错误
func RespUnknownErr(w http.ResponseWriter, err error) {
	var body respBody
	body.Status = UNKNOWN_ERROR
	body.Msg = err.Error()

	b, err := json.Marshal(body)
	if err != nil {
		log.Println("将响应体转换为json字符串失败: ", err)
		// 服务器内部错误
		w.WriteHeader(http.StatusInternalServerError)
	}

	// 将数据写到响应体
	w.Write([]byte(b))
	return
}

// 作为客户端发送get请求
func SendHttpGet(ctx context.Context, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		log.Println("发送get请求报错: ", err)
		return err
	}
	defer resp.Body.Close()
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if string(d) != "pong" {
		return errors.New("响应错误")
	}
	select {
	case <-ctx.Done():
		return errors.New("响应超时")
	default:
		return nil
	}
}

// 转发POST请求
// func RetransmissionPost(req *http.Request) ([]byte, error) {
// 	client := &http.Client{
// 		Timeout: 3 * time.Second,
// 	}
// 	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		log.Println("转发POST请求报错: ", err)
// 		return nil, err
// 	}
// 	defer resp.Body.Close()
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Println("读取转发响应报错: ", err)
// 		return nil, err
// 	}
// 	return body, nil
// }

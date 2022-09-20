package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

type HttpRes struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

// http返回数据
func WriteData(w http.ResponseWriter, res *HttpRes) {
	r, err := json.Marshal(*res)
	if err != nil {
		log.Println("将响应体转换为json字符串失败: ", err)
		// 服务器内部错误
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// 将数据写到响应体
	w.Write([]byte(r))
}

// 作为客户端发送get请求
func SendHttpGet(url string) (string, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	res, err := client.Get(url)
	if err != nil {
		log.Println("发送get请求报错: ", err)
		return "", err
	}
	defer res.Body.Close()
	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := res.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			log.Println("读取响应体报错: ", err)
			return "", err
		}
	}
	return result.String(), nil
}

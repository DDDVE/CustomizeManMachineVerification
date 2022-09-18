package utils

import (
	"encoding/json"
	"log"
	"net/http"
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

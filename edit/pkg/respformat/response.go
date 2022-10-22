package respformat

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type Body struct {
	Status int         `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data,omitempty"`
}

func Response(w http.ResponseWriter, resp interface{}, err error) {
	var body Body
	if err != nil {
		body.Status = -1
		body.Msg = err.Error()
	} else {
		body.Status = 0
		body.Msg = "success"
		body.Data = resp
	}
	httpx.OkJson(w, body)
}

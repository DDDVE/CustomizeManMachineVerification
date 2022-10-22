package middleware

import (
	"context"
	"gate/utils"
	"gate/utils/limiter"
	"gate/utils/log"
	"net/http"
	"strconv"

	"time"
)

/**
流量控制相关函数
*/

var Limiter *limiter.Limiter

const (
	//本次部署时机器序号
	MACHINE_NO = "000001"

	// 网关最大连接数
	BUCKET_SIZE = 500
	RATE        = 200

	LINK_TIME_OUT = 3
)

func InitLimit() {
	Limiter = limiter.NewLimiter(BUCKET_SIZE, RATE)
	log.Printf("Limiter:%+v", Limiter)
}

func Limit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 限流
		ctx, cancel := context.WithTimeout(context.Background(), LINK_TIME_OUT*time.Second)
		i := Limiter.GetToken(ctx)
		cancel()
		if i == -1 {
			utils.RespFormat(w, utils.RUFUSE_LEGAL_REQUEST, nil)
			return
		}

		//生成全局唯一RequestID
		requestID := strconv.FormatInt(time.Now().UnixNano(), 10) + MACHINE_NO + strconv.Itoa(i)
		r.Header.Add("requestID", requestID)

		next.ServeHTTP(w, r)
	}
}

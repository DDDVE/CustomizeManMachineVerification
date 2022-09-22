package utils

import (
	"time"

	"golang.org/x/time/rate"
)

/**
流量控制相关函数
*/

var Limiter *rate.Limiter

const (
	// 网关最大连接数
	MaxRequestNum = 10000
)

func InitLimit() {
	// 每0.1毫秒一块令牌
	Limiter = rate.NewLimiter(rate.Every(time.Millisecond), MaxRequestNum)
}

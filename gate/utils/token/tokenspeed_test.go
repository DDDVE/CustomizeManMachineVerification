package token

import (
	"strconv"
	"testing"
	"time"
)

func TestToken(t *testing.T) {
	GetRandomKey()
	// // jwt 验证token的速度
	// now := time.Now()
	// for i := 0; i < 1; i++ {
	// 	token, err := GenerateTokenString(strconv.Itoa(i))
	// 	if err != nil {
	// 		t.Log("err:", err)
	// 	}
	// 	for j := 0; j < 1000; j++ {
	// 		if CheckTokenString(token) == "" {
	// 			t.Log("搞错了///")
	// 		}
	// 	}
	// }
	// t.Log("总计耗时：", time.Since(now))

	now0 := time.Now()
	for i := 0; i < 1; i++ {
		token := GenerateToken(strconv.Itoa(i))
		for j := 0; j < 1000; j++ {
			if compareToken(token, []byte(strconv.Itoa(i)), 0) {
				t.Log("你也搞错了///")
			}
		}
	}

	t.Log("总计耗时：", time.Since(now0))
}

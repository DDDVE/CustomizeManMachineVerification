package test

/**
测试限流相关
*/
// import (
// 	"bytes"
// 	"gate/utils/log"
// 	"math/rand"
// 	"net"
// 	"net/http"
// 	"strconv"
// 	"testing"
// 	"time"

// 	"golang.org/x/time/rate"
// )

// // 测试golang原生桶令牌
// func TestLimiterServerByLimiter(t *testing.T) {
// 	//表示令牌桶容量为10，每隔100ms放一个令牌到桶里面
// 	limiter := rate.NewLimiter(rate.Every(time.Second/10), 10)
// 	l, _ := net.Listen("tcp", "127.0.0.1:8001")
// 	defer l.Close()
// 	http.Handle("/test", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
// 		reservation := limiter.Reserve()
// 		waitTime := reservation.Delay().Milliseconds()
// 		log.Println("等待毫秒", waitTime, limiter.Limit(), limiter.Burst(), "   ", request.URL)
// 		if waitTime > 1 { //如果有等待时间，则放弃处理当前请求
// 			reservation.Cancel()
// 			writer.WriteHeader(http.StatusGatewayTimeout)
// 			writer.Write([]byte("Error Logic:" + strconv.Itoa(int(waitTime)) + "     -> " + request.URL.String()))
// 			return
// 		}
// 		//模拟业务处理
// 		rand.Seed(time.Now().UnixNano())
// 		sleep := time.Duration(rand.Int()%15+1) * 100
// 		time.Sleep(time.Millisecond * sleep)
// 		writer.Write([]byte("hello:" + strconv.Itoa(int(sleep)) + "     -> " + request.URL.String()))
// 	}))
// 	http.Serve(l, nil)
// }

// // 模拟客户端发送请求
// func TestLimiterClient(t *testing.T) {
// 	for i := 0; i < 300; i++ {
// 		tag := i
// 		if i%20 == 0 {
// 			time.Sleep(time.Second)
// 		}
// 		go func() {
// 			client := &http.Client{Timeout: 100 * time.Second}
// 			resp, err := client.Get("http://127.0.0.1:8001/test?" + strconv.Itoa(tag))
// 			if err != nil {
// 				log.Println("Err:", tag, err)
// 			} else {
// 				var buffer [512]byte
// 				result := bytes.NewBuffer(nil)
// 				n, _ := resp.Body.Read(buffer[0:])
// 				result.Write(buffer[0:n])
// 				log.Println("请求成功", result)
// 			}
// 			defer client.CloseIdleConnections()
// 		}()
// 	}
// 	time.Sleep(time.Second * 200)
// }

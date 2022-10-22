package middleware

import (
	"bufio"
	"gate/utils"
	"gate/utils/log"
	"gate/utils/syncmap"
	"net/http"
	"os"
	"strings"
)

const (
	BLACK_FILE_PATH        = "./conf/blackip.txt"
	BLACK_PERIOD_TASK      = 300
	FIRST_STATISTIC_TIMES  = 10
	SECOND_STATISTIC_TIMES = 60
)

var (
	//一级统计IP,如果请求大于一定次数，将加入二级统计IP并清空一级统计IP
	FirstStatistic = syncmap.NewSyncMap()
	//二级统计IP,如果一段时间内请求大于一定次数，将加入黑名单
	SecondStatistic = syncmap.NewSyncMap()
	//黑名单
	BlackIpMap = syncmap.NewSyncMap()
)

// 初始化黑名单
func InitBlackIp() {
	f, err := os.Open(BLACK_FILE_PATH)
	if err != nil {
		os.Create(BLACK_FILE_PATH)
		log.Println("已创建黑名单")
		return
	}
	rd := bufio.NewReader(f)
	for {
		ip, err := rd.ReadString('\n')
		log.Debugf("err:%v,ip:%v", err, ip)
		if err != nil {
			if len(ip) >= 7 && len(ip) <= 15 {
				BlackIpMap.Set(ip[:len(ip)-1], true)
			}
			break
		}
		if len(ip) >= 7 && len(ip) <= 15 {
			BlackIpMap.Set(ip[:len(ip)-1], true)
		}
	}
	log.Println("读取黑名单完成")
}

//持久化黑名单
func PersistenceBlack(ip string) {
	info, err := os.Stat(BLACK_FILE_PATH)
	if err != nil {
		f, _ := os.Create(BLACK_FILE_PATH)
		k, _ := BlackIpMap.GetAll()
		s := strings.Join(k, "\n")
		f.Write([]byte(s))
		return
	}
	f, err := os.OpenFile(BLACK_FILE_PATH, os.O_RDWR, 0666)
	if err != nil {
		f, _ := os.Create(BLACK_FILE_PATH)
		k, _ := BlackIpMap.GetAll()
		s := strings.Join(k, "\n")
		f.Write([]byte(s))
		return
	}
	b := append([]byte(ip), '\n')
	f.WriteAt(b, info.Size())
}

//黑名单中间件
func BlackIP(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//判断是否在黑名单
		ip := strings.Split(r.RemoteAddr, ":")[0]
		if BlackIpMap.Check(ip) {
			utils.RespFormat(w, utils.BLACK_USER, nil)
			return
		}

		//设置UserIP
		r.Header.Add("userIP", ip)

		//一级统计 大于10次加入二级统计并清空一级统计
		FirstStatistic.GetAndSetUint8(ip, func(oldval uint8) (newval uint8) {
			if oldval > FIRST_STATISTIC_TIMES {
				log.Debug("此ip已加入二级统计,ip:", ip)
				if v, ok := SecondStatistic.Get(ip).(uint16); !ok || v == 0 {
					SecondStatistic.Set(ip, uint16(oldval))
				}
				//置空一级统计map
				FirstStatistic = syncmap.NewSyncMap()
				newval = 0
			} else {
				newval = oldval + 1
			}
			return
		})

		//二级统计 5分钟60次请求加入黑名单
		if SecondStatistic.Check(ip) {
			SecondStatistic.GetAndSetUint16(ip, func(oldval uint16) uint16 {
				if oldval > SECOND_STATISTIC_TIMES {
					log.Debug("此ip已加入黑名单,ip:", ip)
					BlackIpMap.Set(ip, true)
					PersistenceBlack(ip)
				}
				log.Debug("val:", oldval)
				return oldval + 1
			})
		}
		//放行
		next.ServeHTTP(w, r)
	}
}

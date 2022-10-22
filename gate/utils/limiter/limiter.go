package limiter

import (
	"context"
	"gate/utils/log"
	"sync"
	"sync/atomic"
	"time"
)

var (
	flag     bool       //是否正在填充令牌
	flagLock sync.Mutex //只在改变flag时加的锁

	ticker *time.Ticker                        //定时任务
	button chan struct{} = make(chan struct{}) //定时任务开关

	fillTokenPeriod      int = 500 //填充token周期,单位ms
	loopWaitTime         int = 500 //每次循环等待时间,单位ms
	percentOfStartTicker int = 90  //桶在什么时候开始填充桶，单位%，即令牌数在桶90%之下开始填充桶
)

type Limiter struct {
	bucketSize  int    //令牌桶大小
	rate        int    //每秒放入令牌速度
	numOfPer    uint32 //每次循环填充的令牌数量
	warningLine uint32 //剩余令牌数量到警戒线之下则不能取令牌,暂定rate/5
	remainToken uint32 //剩余令牌
}

func NewLimiter(bucketSize int, rate int) *Limiter {
	pernum := rate / (1000 / fillTokenPeriod)
	if pernum <= 0 {
		pernum = 1
	}
	ler := &Limiter{
		bucketSize:  bucketSize,
		rate:        rate,
		numOfPer:    uint32(pernum),
		warningLine: uint32(rate) / 5,
		remainToken: uint32(bucketSize),
	}
	go ler.newScheduleTask()
	return ler
}

//返回int为了生成全局唯一requestID
func (l *Limiter) GetToken(ctx context.Context) int {
	if l.bucketSize <= 0 || l.rate <= 0 {
		return -1
	}
	go l.startFillToken()
	for {
		if l.remainToken > l.warningLine {
			remain := atomic.LoadUint32(&l.remainToken)
			result := atomic.CompareAndSwapUint32(&l.remainToken, remain, remain-1)
			if result == true {
				return int(remain)
			}
		} else {
			time.Sleep(time.Duration(loopWaitTime) * time.Millisecond)
		}
		select {
		case <-ctx.Done():
			return -1
		default:
		}
	}
}

//只获取一次令牌，未取到直接拒绝请求
func (l *Limiter) GetTokenAtOnce() bool {
	if l.bucketSize <= 0 || l.rate <= 0 {
		return false
	}
	go l.startFillToken()
	remain := atomic.LoadUint32(&l.remainToken)
	if remain == 0 {
		return false
	}
	if !atomic.CompareAndSwapUint32(&l.remainToken, remain, remain-1) {
		return false
	}
	return true
}

//通知填充桶的定时任务开始工作
func (l *Limiter) startFillToken() {
	//如果正在填充令牌或桶令牌还比较满，则直接返回
	if flag || l.remainToken >= uint32(l.bucketSize*percentOfStartTicker/100) {
		return
	}
	flagLock.Lock()
	if flag {
		flagLock.Unlock()
		return
	} else {
		flag = true
		flagLock.Unlock()
	}
	select {
	case button <- struct{}{}:
	default:
	}
}

//创建定时任务
func (l *Limiter) newScheduleTask() {
	ticker = time.NewTicker(time.Duration(fillTokenPeriod) * time.Millisecond)
	defer func() {
		ticker.Stop()
	}()
	//开启填充定时任务
	log.Info("开启填充定时任务///")
	for range ticker.C {
		if l.remainToken >= uint32(l.bucketSize) {
			flag = false
			log.Debug("等待填充桶任务通知///")
			<-button
			log.Debug("接到填充桶任务通知///")
		}
		for {
			remain := atomic.LoadUint32(&l.remainToken)
			result := atomic.CompareAndSwapUint32(&l.remainToken, remain, remain+l.numOfPer)
			if result {
				log.Debug("退出循环填充///", remain)
				break
			}
		}
	}
}

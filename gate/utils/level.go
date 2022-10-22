package utils

import (
	"encoding/json"
	"gate/utils/log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

/**
关于请求级别和系统当前级别
*/

// 系统工作路径
var PWD string

// 记录url和处理级别的映射
// 如果某个映射在map
var URLLevelMap map[string]int

// 系统请求级别
var OsLevel int

// var OsLevelRWMutex = new(sync.RWMutex)

const (
	URLLevel1 = 1
	URLLevel2 = 2
	URLLevel3 = 3
	URLLevel4 = 4

	// CPU使用率界限
	// 60以下，URLLevel4；60-80，URLLevel3；80-95：URLLevel2；95以上：URLLevel1
	CPULevelHigh   = 95
	CPULevelMiddle = 80
	CPULevelLow    = 60

	// CPU检查周期为5秒
	CPUCheckPeriod = 5
	// 根据几次取平均值判断CPU利用率
	CPUCheckAverageNum = 5
)

// 初始化服务降级相关
func InitLevel() {
	var err error
	PWD, err = os.Getwd()
	if err != nil {
		log.Panicln("获取工作目录报错: ", err, " 终止程序")
	}
	// 从配置文件中读取各个url的处理级别
	context, err := ReadFile(PWD + "\\urlLevel.txt")
	if err != nil {
		log.Panicln("读取url处理级别文件报错: ", err, " 终止程序")
	}
	err = json.Unmarshal([]byte(context), &URLLevelMap)
	if err != nil {
		log.Panicln("解析url处理级别文件报错: ", err, " 终止程序")
	}
	// 初始化系统级别为最低
	OsLevel = URLLevel4

	// 起一个协程负责监控cpu使用量和调节系统处理请求的级别
	go CheckOsLevel()
}

func CheckOsLevel() {
	log.Println("开始周期性检查CPU使用率...")
	ticker := time.NewTicker(CPUCheckPeriod * time.Second)
	defer func() {
		log.Println("周期性检查CPU使用率退出")
		ticker.Stop()
	}()
	i := 0
	osLevelRecord := make([]int, CPUCheckAverageNum)
	for range ticker.C {
		// 获取cpu使用率
		cmd := exec.Command("wmic", "cpu", "get", "loadpercentage")
		r, err := cmd.Output()
		if err != nil {
			log.Printf("%+v 获取cpu利用率报错: %+v\n", time.Now(), err)
			continue
		}
		r1 := strings.Split(string(r), "\n")[1][:4]
		for i := len(r1) - 1; i >= 0; i-- {
			if r1[i] == ' ' {
				r1 = string([]byte(r1)[:len(r1)-1])
			} else {
				break
			}
		}
		rate, err := strconv.Atoi(r1)
		if err != nil {
			log.Printf("%+v 解析cpu利用率报错: %+v\n", time.Now(), err)
			continue
		}
		log.Printf("当前CPU使用率: %d\n", rate)
		osLevelRecord[i] = rate
		i = (i + 1) % CPUCheckAverageNum
		// 求五次平均值再判断是否要改变级别
		rate = Sum(osLevelRecord) / CPUCheckAverageNum
		if rate >= CPULevelLow && rate < CPULevelMiddle {
			// OsLevelRWMutex.Lock()
			log.Printf("CPU平均使用率为: %d, 修改系统请求级别为: %d\n", rate, URLLevel3)
			OsLevel = URLLevel3
			// OsLevelRWMutex.Unlock()
		} else if rate >= CPULevelMiddle && rate < CPULevelHigh {
			// OsLevelRWMutex.Lock()
			log.Printf("CPU平均使用率为: %d, 修改系统请求级别为: %d\n", rate, URLLevel2)
			OsLevel = URLLevel2
			// OsLevelRWMutex.Unlock()
		} else if rate >= CPULevelHigh {
			// OsLevelRWMutex.Lock()
			log.Printf("CPU平均使用率为: %d, 修改系统请求级别为: %d\n", rate, URLLevel1)
			OsLevel = URLLevel1
			// OsLevelRWMutex.Unlock()
		} else {
			// OsLevelRWMutex.Lock()
			log.Printf("CPU平均使用率为: %d, 修改系统请求级别为: %d\n", rate, URLLevel4)
			OsLevel = URLLevel4
			// OsLevelRWMutex.Unlock()
		}
	}

}

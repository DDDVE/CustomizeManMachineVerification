package log

import (
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

/*
封装log
*/

const (
	//日志级别
	LEVEL_DEBUG = 0
	LEVEL_INFO  = 1
	LEVEL_WARN  = 2
	LEVEL_ERROR = 3
	//日志输出文件名后缀
	FILE_NAME_SUFFIX = ".log"
)

type LogWithPublicString struct {
	strs []string
	sep  string //分割符
}

var (
	logLock     sync.Mutex
	outFilePath string //输出路径_ 默认输出控制台
	logLevel    int    //默认debug

	shortName = map[int]string{ //日志级别对应文件名
		0: "debug",
		1: "info",
		2: "warn",
		3: "error",
	}
	//默认打印模板
	printTemplate func(level int, format string, v ...interface{}) = func(level int, format string, v ...interface{}) {
		logAddr := "nil"
		caller := "nil"
		pc2, _, lineNo, ok := runtime.Caller(2)
		if ok {
			name := runtime.FuncForPC(pc2).Name()
			logAddr = name + "_" + strconv.Itoa(lineNo)
		}
		pc3, _, lineNo, ok := runtime.Caller(3)
		if ok {
			name := runtime.FuncForPC(pc3).Name()
			caller = name + "_" + strconv.Itoa(lineNo)
		}
		if outFilePath != "" {
			logLock.Lock()
			defer logLock.Unlock()
			setOutFile(level)
		}
		if format == "" {
			v = append(v, " logAddr=="+logAddr+" caller=="+caller)
			log.Println(v...)
		} else {
			v = append(v, logAddr, caller)
			log.Printf(format+" logAddr==%v caller==%v", v...)
		}
	}
)

//新建一个有公共字符串的log，每次打印日志都会自动打印公共字符串
func NewLogWithString(publicStrings ...string) *LogWithPublicString {
	return &LogWithPublicString{
		strs: publicStrings,
		sep:  " ",
	}
}

//设置分隔符
func (l *LogWithPublicString) SetSeg(sep string) {
	l.sep = sep
}

//添加publicStrings
func (l *LogWithPublicString) Add(publicStrings ...string) {
	l.strs = append(l.strs, publicStrings...)
}

func (l *LogWithPublicString) Printf(format string, v ...interface{}) {
	public := strings.Join(l.strs, l.sep)
	Printf(public+l.sep+format, v...)
}

func (l *LogWithPublicString) Print(format string, v ...interface{}) {
	public := strings.Join(l.strs, l.sep)
	v = append(v, public)
	Println(v...)
}

func (l *LogWithPublicString) Errorf(format string, v ...interface{}) {
	public := strings.Join(l.strs, l.sep)
	Errorf(public+l.sep+format, v...)
}

func (l *LogWithPublicString) Error(format string, v ...interface{}) {
	public := strings.Join(l.strs, l.sep)
	v = append(v, public)
	Error(v...)
}

//Print默认info日志级别
func Printf(format string, v ...interface{}) {
	if logLevel > LEVEL_INFO {
		return
	}
	printTemplate(LEVEL_INFO, format, v...)
}
func Println(v ...interface{}) {
	if logLevel > LEVEL_INFO {
		return
	}
	printTemplate(LEVEL_INFO, "", v...)
}

func Debugf(format string, v ...interface{}) {
	if logLevel > LEVEL_DEBUG {
		return
	}
	printTemplate(LEVEL_DEBUG, format, v...)
}

//默认换行
func Debug(v ...interface{}) {
	if logLevel > LEVEL_DEBUG {
		return
	}
	printTemplate(LEVEL_DEBUG, "", v...)
}

func Infof(format string, v ...interface{}) {
	if logLevel > LEVEL_INFO {
		return
	}
	printTemplate(LEVEL_INFO, format, v...)
}

//默认换行
func Info(v ...interface{}) {
	if logLevel > LEVEL_INFO {
		return
	}
	printTemplate(LEVEL_INFO, "", v...)
}

func Warnf(format string, v ...interface{}) {
	if logLevel > LEVEL_WARN {
		return
	}
	printTemplate(LEVEL_WARN, format, v...)
}

//默认换行
func Warn(v ...interface{}) {
	if logLevel > LEVEL_WARN {
		return
	}
	printTemplate(LEVEL_WARN, "", v...)
}

func Errorf(format string, v ...interface{}) {
	printTemplate(LEVEL_ERROR, format, v...)
}

//默认换行
func Error(v ...interface{}) {
	printTemplate(LEVEL_ERROR, "", v...)
}

func Panic(v ...interface{}) {
	log.Panic(v...)
}
func Panicf(format string, v ...interface{}) {
	log.Panicf(format, v...)
}
func Panicln(v ...interface{}) {
	log.Panicln(v...)
}

//设置日志输出路径
func SetOutPath(path string) {
	s, err := os.Stat(path)
	if err == nil && s.IsDir() {
		outFilePath = path
	}
}

//设置日志级别
func SetLevel(level int) {
	if level > LEVEL_DEBUG && level <= LEVEL_ERROR {
		logLevel = level
	}
}

//设置打印模板
func SetTemplate(template func(level int, format string, v ...interface{})) {
	printTemplate = template
}

//设置日志输出文件
func setOutFile(level int) bool {
	longFileName := outFilePath + shortName[level] + FILE_NAME_SUFFIX
	// info, err := os.Stat(longFileName)
	// if err != nil {
	// 	f, err := os.Create(longFileName)
	// 	if err != nil {
	// 		outFilePath = ""
	// 		return false
	// 	}
	// 	defer f.Close()
	// 	log.SetOutput(f)
	// } else {
	f, err := os.OpenFile(longFileName, os.O_RDWR, 0666)
	if err != nil {
		outFilePath = ""
		return false
	}
	defer f.Close()
	f.Seek(0, os.SEEK_END)
	log.SetOutput(f)
	// }

	return true
}

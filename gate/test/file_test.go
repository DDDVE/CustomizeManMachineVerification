package test

import (
	"gate/utils/log"
	"os"
	"path/filepath"
	"testing"
)

var ApiData = [5]string{"11", "22", "33", "44", "55"}

func TestFile(t *testing.T) {
	PWD, err := os.Getwd()
	if err != nil {
		log.Error("获取工作目录报错: ", err)
	}
	filepath.Join(PWD, "../conf/test.log")
	// f, err := os.OpenFile("./test.txt", os.O_RDWR, 0666)
	// defer f.Close()
	// if err != nil {
	// 	t.Log(err)
	// }
	// reader := bufio.NewReader(f)
	// pos := int64(0)
	// flag := false
	// id := 1
	// for {
	// 	//读取每一行内容
	// 	line, err := reader.ReadString('\n')
	// 	t.Log("/////////////line:", line, "////err:", err)
	// 	if err != nil && err != io.EOF {
	// 		return
	// 	}
	// 	l := len(line)
	// 	if id == 5 {
	// 		return
	// 	}
	// 	//根据关键词覆盖当前行
	// 	if strings.Contains(line, ApiData[id]) {
	// 		flag = true
	// 	}
	// 	if flag {
	// 		bytes := []byte(ApiData[id])
	// 		bytes = append(bytes, "\n"...)
	// 		f.WriteAt(bytes, pos)
	// 		l = len(bytes)
	// 		id++
	// 	}
	// 	//读到末尾
	// 	if err == io.EOF {
	// 		break
	// 	}
	// 	// 每一行读取完后记录位置
	// 	pos += int64(l)
	// }
}

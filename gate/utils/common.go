package utils

import (
	"errors"
	"io"
	"math"
	"os"
)

func MaxOfMany(a ...int) int {
	ans := math.MinInt32
	for i := 0; i < len(a); i++ {
		if a[i] > ans {
			ans = a[i]
		}
	}
	return ans
}

func ReadFile(dir string) (string, error) {
	file, err := os.Open(dir)
	if err != nil {
		return "", err
	}
	if file == nil {
		return "", errors.New("文件为空")
	}
	defer file.Close()
	// 读取文件内容缓存
	buf := make([]byte, 512)
	context := []byte{}
	for {
		count, err := file.Read(buf)
		// 判断是否读到文件尾部
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		curByte := buf[:count]
		context = append(context, curByte...)
	}
	return string(context), nil
}

package test

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"testing"
)

// 测试执行外部方法获取cpu利用率
func TestExecCommand(t *testing.T) {
	cmd := exec.Command("wmic", "cpu", "get", "loadpercentage")
	r, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	rate := strings.Split(string(r), "\n")[1][:4]
	fmt.Print(rate)
	for i := len(rate) - 1; i >= 0; i-- {
		if rate[i] == ' ' {
			rate = string([]byte(rate)[:len(rate)-1])
		} else {
			break
		}
	}
	fmt.Println(rate)
	r1, err := strconv.Atoi(rate)
	if err != nil {
		fmt.Printf("报错: %+v\n", err)
		return
	}
	fmt.Println(r1)
}

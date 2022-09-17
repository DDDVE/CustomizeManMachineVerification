package test

import (
	"gate/schedule"
	"gate/utils"
	"testing"
)

// 测试生成随机字符串
func TestGetRandomString(t *testing.T) {
	t.Log(utils.GetRandomString(50))
}

// 测试开启定时任务是否成功
func TestInit(t *testing.T) {
	schedule.InitGenerateKeyScheTask()
}

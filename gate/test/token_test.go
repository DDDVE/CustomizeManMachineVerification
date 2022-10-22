package test

import (
	"gate/initialize/schedule"
	"testing"
)

// 测试开启定时任务是否成功
func TestInit(t *testing.T) {
	schedule.InitGenerateKeyScheTask()
}

func TestIfTheSame(t *testing.T) {
	t.Log("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJtb2JpbGVfbnVtIjoiMTg4MTEwMzUyMTMiLCJpc3MiOiJkdmUiLCJzdWIiOiJjdXN0b20gbWFuLW1hY2hpbmUgdmVyaWZ5IHBsYXQifQ.J6J-LOV0P5tP6_JGpZy5pyhD-fYgZRQM0khPRUDBVjc" == "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJtb2JpbGVfbnVtIjoiMTg4MTEwMzUyMTMiLCJpc3MiOiJkdmUiLCJzdWIiOiJjdXN0b20gbWFuLW1hY2hpbmUgdmVyaWZ5IHBsYXQifQ.LOsxn0OfN-89hGQ-_jtGTYNnk8j1KkQ0jsLDk-Vmn_M")
}

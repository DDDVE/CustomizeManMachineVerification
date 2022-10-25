package middleware

import (
	"testing"
)

var m = make(map[int]int, 100)

func TestMap(t *testing.T) {
}
func FuzzMap(f *testing.F) {
	for i := 0; i < 2; i++ {
		f.Add(i) //种子参数
	}
	f.Fuzz(func(t *testing.T, rand int) {
		//执行函数入参为i
	})
}

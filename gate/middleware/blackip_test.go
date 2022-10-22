package middleware

import (
	"testing"
	"time"
)

var m = make(map[int]int, 100)

func TestMap(t *testing.T) {
	go read()
	s := 0
	for i := 1000; i < 2000; i++ {
		s += write(i)
	}
	time.Sleep(time.Second)
}

func read() {
	s := 0
	for i := 0; i < 1000; i++ {
		s += len(m)
	}
}

func write(i int) int {
	m[i] = i
	return i
}

func FuzzMap(f *testing.F) {
	for i := 0; i < 2; i++ {
		f.Add(i) //种子参数
	}
	f.Fuzz(func(t *testing.T, rand int) {
		//执行函数入参为i
	})
}

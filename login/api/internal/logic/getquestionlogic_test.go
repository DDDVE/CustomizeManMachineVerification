package logic

import (
	"fmt"
	"testing"
)

func TestRandom(t *testing.T) {
	fmt.Println(string(RandomAnswer("中国", "的啊")))
}

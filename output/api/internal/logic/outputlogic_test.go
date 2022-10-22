package logic

import (
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	stamp := time.Now().UnixNano()
	stamp -= 7000
	t.Log("stamp:", stamp)
	// timestamp, err := strconv.ParseInt(stamp, 10, 64)
	// if err != nil {
	// 	return nil, err
	// }
	x := time.Since(time.Unix(0, stamp))
	if x < 4800 {
		t.Log("xiao")
	}
	if x > 4900 {
		t.Log("zhong")
	}
	if x > 6000 {
		t.Log("da")
	}
}

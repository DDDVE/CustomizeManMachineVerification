package utils

import "math"

func MaxOfMany(a ...int) int {
	ans := math.MinInt32
	for i := 0; i < len(a); i++ {
		if a[i] > ans {
			ans = a[i]
		}
	}
	return ans
}

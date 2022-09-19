package utils

import (
	"hash/crc32"
)

// 采用负载均衡返回某个api网关的序号
// 第一版先用普通的哈希取余
func FindApiGateToRedirect(apiType string, ip string) int {
	// 根据hash值和某种api网关的个数取余
	hash := crc32.ChecksumIEEE([]byte(ip))
	return int(hash)
}

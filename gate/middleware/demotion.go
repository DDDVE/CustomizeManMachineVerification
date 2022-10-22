package middleware

import "net/http"

func Demotion(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// // 先获取url然后判断url级别和当前系统级别
		// path := r.URL.Path
		// reqLevel, ok := utils.URLLevelMap[path]
		// // 如果此路径不存在
		// if !ok {
		// 	// log.Printf("主机%s请求的路径%s不存在\n", ip, path)
		// 	utils.WriteData(w, &utils.HttpRes{
		// 		Status: utils.HttpUrlCheckFalse,
		// 		Data:   nil,
		// 	})
		// 	return
		// }
		// // 如果此时url级别低于系统级别(数字更大)
		// // utils.OsLevelRWMutex.RLock()
		// if reqLevel > utils.OsLevel {
		// 	// utils.OsLevelRWMutex.RUnlock()
		// 	// log.Printf("主机%s请求的路径%s级别过低: %d\n", ip, path, reqLevel)
		// 	utils.WriteData(w, &utils.HttpRes{
		// 		Status: utils.HttpRefuse,
		// 		Data:   nil,
		// 	})
		// 	return
		// }
		// utils.OsLevelRWMutex.RUnlock()
		// 级别高于等于系统级别就放行
		next.ServeHTTP(w, r)
	}
}

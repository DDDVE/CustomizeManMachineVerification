package middleware

import (
	"gate/utils"
	"gate/utils/log"
	"gate/utils/token"
	"net/http"
	"strings"
)

func TokenCheck(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userToken := r.Header.Get("Authorization")
		log.Debug("userToken:///", userToken)
		id, mobileNum := token.CheckToken(userToken)
		log.Debug("id:///", id, "mobileNum:///", mobileNum)
		if id == -1 {
			utils.RespFormat(w, utils.TOKEN_CHECK_ERROR, nil)
			return
		}
		if id != 0 && mobileNum != "" {
			token := token.GenerateToken(mobileNum)
			w.Header().Add("Access-Control-Expose-Headers", "Authorization, Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
			w.Header().Add("Authorization", token)
		}
		//检查token，用于跳过登录页面
		if strings.Contains(r.URL.Path, "checktoken") {
			utils.RespFormat(w, utils.SUCCESS, nil)
			return
		}

		next.ServeHTTP(w, r)
	}
}

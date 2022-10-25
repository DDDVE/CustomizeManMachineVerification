package handler

import (
	"gate/utils"
	"gate/utils/log"
	"net/http"
)

func DistributeReq(w http.ResponseWriter, r *http.Request) {
	//不合法的URL直接在降级中间件过滤掉，所以这里暂不做细致校验
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		utils.RespFormat(w, utils.REQUEST_METHOD_ERROR, nil)
	}
	log.Debug("URL/////////", r.URL.String())
	//解析请求的模块名
	moduleID := getModuleID(r.URL.Path)
	if moduleID == -1 {
		utils.RespFormat(w, utils.REQUEST_PATH_ERROR, nil)
	}
	log.Debug("moduleID/////////", moduleID)
	//通过对userAddr哈希取余方式负载均衡
	id := loadBalance(r.RemoteAddr, moduleID)
	log.Debug("id/////////", id)

	//转发请求
	err := forwardReq(w, r, ApiData[moduleID].ApiAddrs[id])
	if err != nil {
		utils.RespUnknownErr(w, err)
	}

	// w.Write(resp)
}

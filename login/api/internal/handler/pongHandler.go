package handler

import (
	"net/http"

	"api/internal/svc"
)

func pongHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
		return
	}
}

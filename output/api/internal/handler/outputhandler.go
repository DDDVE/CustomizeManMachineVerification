package handler

import (
	"net/http"
	"pkg/respformat"

	"api/internal/logic"
	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func outputHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OutputReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewOutputLogic(r.Context(), svcCtx)
		if err := l.CheckReq(&req); err != nil {
			httpx.Error(w, err)
			return
		}
		resp, err := l.Output(&req)
		respformat.Response(w, resp, err)
	}
}

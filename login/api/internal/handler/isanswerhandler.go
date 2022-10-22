package handler

import (
	"net/http"
	"pkg/respformat"

	"api/internal/logic"
	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func isAnswerHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.IsAnswerReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewIsAnswerLogic(r.Context(), svcCtx)
		resp, err := l.IsAnswer(&req)
		respformat.Response(w, resp, err)
	}
}

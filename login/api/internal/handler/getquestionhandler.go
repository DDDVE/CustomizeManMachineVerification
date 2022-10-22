package handler

import (
	"net/http"
	"pkg/respformat"

	"api/internal/logic"
	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func getQuestionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.QuestionReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}
		l := logic.NewGetQuestionLogic(r.Context(), svcCtx)
		resp, err := l.GetQuestion(&req)
		respformat.Response(w, resp, err)
	}
}

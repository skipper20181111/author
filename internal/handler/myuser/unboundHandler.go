package myuser

import (
	"net/http"

	"author/internal/logic/myuser"
	"author/internal/svc"
	"author/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UnboundHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginRes
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := myuser.NewUnboundLogic(r.Context(), svcCtx)
		resp, err := l.Unbound(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}

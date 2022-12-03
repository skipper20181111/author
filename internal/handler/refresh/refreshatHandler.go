package refresh

import (
	"net/http"

	"author/internal/logic/refresh"
	"author/internal/svc"
	"author/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func RefreshatHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RefreshRes
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := refresh.NewRefreshatLogic(r.Context(), svcCtx)
		resp, err := l.Refreshat(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}

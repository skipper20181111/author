package myuser

import (
	"net/http"

	"author/internal/logic/myuser"
	"author/internal/svc"
	"author/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetinfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetUserInfoRes
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := myuser.NewGetinfoLogic(r.Context(), svcCtx)
		resp, err := l.Getinfo(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}

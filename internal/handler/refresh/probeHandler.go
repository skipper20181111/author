package refresh

import (
	"net/http"

	"author/internal/logic/refresh"
	"author/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ProbeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := refresh.NewProbeLogic(r.Context(), svcCtx)
		err := l.Probe()
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.Ok(w)
		}
	}
}

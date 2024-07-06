package handler

import (
	"net/http"

	"github.com/atomicals-go/atomicals-api/internal/logic"
	"github.com/atomicals-go/atomicals-api/internal/svc"
	"github.com/atomicals-go/atomicals-api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func getAssetByLocationIDHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ReqAssetByLocationID
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewGetassetByLocationIDLogic(r.Context(), svcCtx)
		resp, err := l.GetAssetByLocationID(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

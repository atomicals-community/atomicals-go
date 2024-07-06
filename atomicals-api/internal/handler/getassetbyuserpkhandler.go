package handler

import (
	"net/http"

	"github.com/atomicals-go/atomicals-api/internal/logic"
	"github.com/atomicals-go/atomicals-api/internal/svc"
	"github.com/atomicals-go/atomicals-api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func getAssetByUserPkHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ReqAssetByUserPK
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewGetAssetByUserPkLogic(r.Context(), svcCtx)
		resp, err := l.GetAssetByUserPk(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

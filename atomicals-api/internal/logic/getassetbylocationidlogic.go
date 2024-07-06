package logic

import (
	"context"

	"github.com/atomicals-go/atomicals-api/internal/svc"
	"github.com/atomicals-go/atomicals-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetassetByLocationIDLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetassetByLocationIDLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetassetByLocationIDLogic {
	return &GetassetByLocationIDLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetassetByLocationIDLogic) GetAssetByLocationID(req *types.ReqAssetByLocationID) (resp *types.RespAssetByLocationID, err error) {
	entities, err := l.svcCtx.NftUTXOsByLocationID(req.LocationID)
	if err != nil {
		l.Errorf("[GetAssetByLocationID] NftUTXOsByLocationID err:%v", err)
		return
	}
	for _, v := range entities {
		resp.Assets = append(resp.Assets, &types.UTXONftInfo{
			UserPk:                     v.UserPk,
			AtomicalsID:                v.AtomicalsID,
			LocationID:                 v.LocationID,
			RealmName:                  v.RealmName,
			SubRealmName:               v.SubRealmName,
			ParentRealmAtomicalsID:     v.ParentRealmAtomicalsID,
			ContainerName:              v.ContainerName,
			Dmitem:                     v.Dmitem,
			ParentContainerAtomicalsID: v.ParentContainerAtomicalsID,
			Nonce:                      v.Nonce,
			Time:                       v.Time,
			Bitworkc:                   v.Bitworkc,
			Bitworkr:                   v.Bitworkr,
		})
	}
	return
}

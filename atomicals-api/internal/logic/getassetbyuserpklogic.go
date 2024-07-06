package logic

import (
	"context"

	"github.com/atomicals-go/atomicals-api/internal/svc"
	"github.com/atomicals-go/atomicals-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAssetByUserPkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAssetByUserPkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAssetByUserPkLogic {
	return &GetAssetByUserPkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAssetByUserPkLogic) GetAssetByUserPk(req *types.ReqAssetByUserPK) (resp *types.RespAssetByUserPK, err error) {
	nfts, err := l.svcCtx.NftUTXOsByUserPK(req.UserPK)
	if err != nil {
		l.Errorf("[GetAssetByUserPk] NftUTXOsByUserPK err:%v", err)
		return
	}
	for _, v := range nfts {
		resp.NftAssets = append(resp.NftAssets, &types.UTXONftInfo{
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
	fts, err := l.svcCtx.FtUTXOsByUserPK(req.UserPK)
	if err != nil {
		l.Errorf("[GetAssetByUserPk] FtUTXOsByUserPK err:%v", err)
		return
	}
	for _, v := range fts {
		resp.FtAssets = append(resp.FtAssets, &types.UTXOFtInfo{
			UserPk:          v.UserPk,
			AtomicalsID:     v.AtomicalsID,
			LocationID:      v.LocationID,
			Bitworkc:        v.Bitworkc,
			Bitworkr:        v.Bitworkr,
			MintTicker:      v.MintTicker,
			Nonce:           v.Nonce,
			Time:            v.Time,
			MintBitworkVec:  v.MintBitworkVec,
			MintBitworkcInc: v.MintBitworkcInc,
			MintBitworkrInc: v.MintBitworkrInc,
			Amount:          v.Amount,
			Type:            v.Type,
			Subtype:         v.Subtype,
			TickerName:      v.TickerName,
			MaxSupply:       v.MaxSupply,
			MintAmount:      v.MintAmount,
			MintHeight:      v.MintHeight,
			MaxMints:        v.MaxMints,
		})
	}
	return
}

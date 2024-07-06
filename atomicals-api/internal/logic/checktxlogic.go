package logic

import (
	"context"
	"encoding/json"

	"github.com/atomicals-go/atomicals-api/internal/svc"
	"github.com/atomicals-go/atomicals-api/internal/types"
	"github.com/atomicals-go/repo/postsql"
	"github.com/atomicals-go/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckTxLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckTxLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckTxLogic {
	return &CheckTxLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckTxLogic) CheckTx(req *types.ReqCheckTx) (resp *types.RespCheckTx, err error) {
	txID, _ := utils.SplitAtomicalsID(req.LocationID)
	tx, height, err := l.svcCtx.GetTxByTxID(txID)
	if err != nil {
		l.Errorf("[CheckTx] GetTxByTxID err:%v", err)
		return
	}
	if height <= l.svcCtx.SyncHeight {
		resp.Status = "confirmed"
		var nftAssets []*postsql.UTXONftInfo
		nftAssets, err = l.svcCtx.NftUTXOsByLocationID(req.LocationID)
		if err != nil {
			l.Errorf("[CheckTx] NftUTXOsByLocationID err:%v", err)
			return
		}
		if len(nftAssets) != 0 {
			var res []byte
			res, err = json.Marshal(nftAssets)
			if err != nil {
				l.Errorf("[CheckTx] Marshal err:%v", err)
				return
			}
			resp.Description += string(res)
		}
		var ftAssets []*postsql.UTXOFtInfo
		ftAssets, err = l.svcCtx.FtUTXOsByLocationID(req.LocationID)
		if err != nil {
			l.Errorf("[CheckTx] NftUTXOsByLocationID err:%v", err)
			return
		}
		if len(ftAssets) != 0 {
			var res []byte
			res, err = json.Marshal(ftAssets)
			if err != nil {
				l.Errorf("[CheckTx] Marshal err:%v", err)
				return
			}
			resp.Description += string(res)
		}
	} else if l.svcCtx.SyncHeight < height && height < l.svcCtx.MaxBlockHeight {
		resp.Status = "until confirmation depth"
		resp.Description += l.svcCtx.PendingAtomicalsAsset.CheckAssetByLocationID(req.LocationID)

	} else if l.svcCtx.MaxBlockHeight < height {
		resp.Status = "in mempool"
		p := l.svcCtx.SyncMempoolAtomicalsAsset(*tx, height)
		resp.Description += p.CheckAssetByLocationID(req.LocationID)
	}
	return
}

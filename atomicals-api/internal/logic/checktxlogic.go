package logic

import (
	"context"
	"errors"

	"github.com/atomicals-go/atomicals-api/internal/svc"
	"github.com/atomicals-go/atomicals-api/internal/types"
	"github.com/atomicals-go/repo/postsql"

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
	tx, height, err := l.svcCtx.GetTxByTxID(req.Txid)
	if err != nil {
		l.Errorf("[CheckTx] GetTxByTxID err:%v", err)
		return
	}
	if 0 <= height && height <= l.svcCtx.SyncHeight {
		resp.Status = "confirmed"
		txRecord := &postsql.AtomicalsTx{}
		txRecord, err = l.svcCtx.AtomicalsTx(req.Txid)
		if err != nil {
			l.Errorf("[CheckTx] AtomicalsTx err:%v", err)
			return
		}
		resp.Operation = txRecord.Operation
		resp.Description = txRecord.Description
	} else if l.svcCtx.SyncHeight < height && height < l.svcCtx.MaxBlockHeight {
		resp.Status = "until confirmation depth"
		pendingAssets, ok := l.svcCtx.PendingAtomicalsAssetMap[req.Txid]
		if !ok {
			l.Errorf("[CheckTx] AtomicalsTx err:%v", errors.New("atomicals operation not found"))
			return
		}
		resp.Description = pendingAssets.CheckAsset()
		resp.Operation = pendingAssets.Operation
	} else if height < 0 {
		resp.Status = "in mempool"
		pendingAssets := l.svcCtx.SyncMempoolAtomicalsAsset(*tx, height)
		resp.Description = pendingAssets.CheckAsset()
		resp.Operation = pendingAssets.Operation
	}
	return
}

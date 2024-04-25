package atomicals

import (
	"time"

	"github.com/atomicals-core/pkg/log"

	"github.com/atomicals-core/atomicals/common"
	"github.com/atomicals-core/atomicals/witness"
	"github.com/btcsuite/btcd/btcjson"
)

func (m *Atomicals) TraceBlock() {
	startTime := time.Now()
	height, err := m.CurrentHeitht()
	if err != nil {
		log.Log.Panicf("CurrentHeitht err:%v", err)
	}
	height++
	blockInfo, err := m.GetBlockByHeight(height)
	if err != nil {
		log.Log.Panicf("GetBlockByHeight err:%v height:%v", err, height)
	}
	getBlockByHeightTime := time.Since(startTime)
	startTime = time.Now()
	for _, tx := range blockInfo.Tx {
		// skip this tx, it's from miner
		if tx.Vin[0].Txid == "" {
			continue
		}
		m.TraceTx(tx, blockInfo.Height)
	}
	if err := m.UpdateCurrentHeightAndExecAllSql(blockInfo.Height); err != nil {
		log.Log.Panicf("UpdateCurrentHeight err:%v", err)
	}
	log.Log.Infof("height:%v, getBlockByHeight take time:%v, TraceTx take time:%v", blockInfo.Height, getBlockByHeightTime, time.Since(startTime))
}

func (m *Atomicals) TraceTx(tx btcjson.TxRawResult, height int64) error {
	operation := witness.ParseWitness(tx, height)
	// step 1: transfer nft, transfer ft
	// m.transferNft(operation, tx)
	if height < common.AtOMICALS_FT_PARTIAL_SPLITING_HEIGHT {
		m.transferFt(operation, tx)
	} else {
		m.transferFtPartialSpliting(operation, tx)
	}

	// step 2: process operation
	userPk := tx.Vout[common.VOUT_EXPECT_OUTPUT_INDEX].ScriptPubKey.Address
	if operation.Op == "dmt" {
		m.mintDistributedFt(operation, tx.Vout, userPk)
	} else {
		switch operation.Op {
		case "dft":
			m.deployDistributedFt(operation, userPk)
		case "ft":
			m.mintDirectFt(operation, tx.Vout, userPk)
		case "nft":
			// m.mintNft(operation, userPk, height)
		case "mod":
		case "evt":
		case "dat":
		case "sl":
		default:
		}
	}

	// step 3 check payment
	return nil
}

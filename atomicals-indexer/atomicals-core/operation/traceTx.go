package atomicals

import (
	"time"

	"github.com/atomicals-go/atomicals-indexer/atomicals-core/witness"
	"github.com/atomicals-go/pkg/log"
	"github.com/atomicals-go/repo"
	"github.com/atomicals-go/repo/postsql"
	"github.com/atomicals-go/utils"
	"github.com/btcsuite/btcd/btcjson"
)

func (m *Atomicals) Run() {
	startTime := time.Now()
	if m.location.BlockHeight+utils.SafeBlockHeightInterupt > m.maxBlockHeight {
		time.Sleep(10 * time.Minute)
		var err error
		m.maxBlockHeight, err = m.GetBlockCount()
		if err != nil {
			log.Log.Panicf("GetBlockCount err:%v", err)
		}
	}
	block, err := m.GetBlockByHeight(m.location.BlockHeight)
	if err != nil {
		log.Log.Panicf("GetBlockByHeight err:%v", err)
	}
	m.location.TxIndex++
	if m.location.TxIndex >= int64(len(block.Tx)) {
		m.location.BlockHeight++
		m.location.TxIndex = 0
		block, err = m.GetBlockByHeight(m.location.BlockHeight)
		if err != nil {
			log.Log.Panicf("GetBlockByHeight err:%v", err)
		}
	}
	for txIndex, tx := range block.Tx {
		if int64(txIndex) < m.location.TxIndex {
			continue
		}
		m.location.TxIndex = int64(txIndex)
		data := m.TraceTx(tx, block.Height)
		if (data != nil && data.Op != "") || (block.Height%10 == 0 && m.location.TxIndex == 0) {
			err = m.UpdateDB(block.Height, m.location.TxIndex, tx.Txid, data)
			if err != nil {
				log.Log.Panicf("UpdateDB err:%v", err)
			}
		}
	}

	log.Log.Infof("maxBlockHeight:%v, currentHeight:%v, time:%v", m.maxBlockHeight, m.location.BlockHeight, time.Since(startTime))
}

func (m *Atomicals) TraceTx(tx btcjson.TxRawResult, height int64) *repo.AtomicaslData {
	operation := witness.ParseWitness(tx, height)
	if operation.Payload != nil && !(operation.Payload.Args.MintTicker == "atom" || operation.Payload.Args.RequestTicker == "atom") {
		return nil
	}

	// step 1: insert mod
	// if operation.Op == "mod" {
	// 	mod = m.operationMod(operation, tx)
	// }

	// step 2: transfer nft, transfer ft
	deleteFts, newFts, _ := m.transferFt(operation, tx)

	// updateNfts, _ = m.transferNft(operation, tx)

	// step 3: process operation
	var newUTXOFtInfo *postsql.UTXOFtInfo
	var updateDistributedFt *postsql.GlobalDistributedFt
	var newGlobalDistributedFt *postsql.GlobalDistributedFt
	userPk := tx.Vout[utils.VOUT_EXPECT_OUTPUT_INDEX].ScriptPubKey.Address
	if operation.Op == "dmt" {
		newUTXOFtInfo, updateDistributedFt, _ = m.mintDistributedFt(operation, tx.Vout, userPk)
	} else {
		switch operation.Op {
		case "dft":
			newGlobalDistributedFt, _ = m.deployDistributedFt(operation, userPk)
		case "ft":
			// newGlobalDirectFt, _ = m.mintDirectFt(operation, tx.Vout, userPk)
		case "nft":
			// newUTXONftInfo, _ = m.mintNft(operation, userPk)
		case "evt":
		case "dat":
		case "sl":
		default:
		}
	}

	// TODO: step 4 check payment
	data := &repo.AtomicaslData{
		// Mod:                    mod,
		DeleteFts: deleteFts,
		NewFts:    newFts,
		// UpdateNfts:             updateNfts,
		NewUTXOFtInfo:          newUTXOFtInfo,
		UpdateDistributedFt:    updateDistributedFt,
		NewGlobalDistributedFt: newGlobalDistributedFt,
		// NewGlobalDirectFt:      newGlobalDirectFt,
		// NewUTXONftInfo:         newUTXONftInfo,
	}
	data.ParseOperation()

	return data
}

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
		panic(err)
	}
	blockInfo, err := m.btcClient.GetBlockByHeight(height + 1)
	if err != nil {
		log.Log.Panicf("GetBlockByHeight err:%v height:%v", err, height)
		panic(err)
	}
	log.Log.Warnf("time.Since(startTime):%v, height:%v", time.Since(startTime), blockInfo.Height)

	startTime = time.Now()
	for _, tx := range blockInfo.Tx {
		// skip this tx, it's from miner
		if tx.Vin[0].Txid == "" {
			continue
		}
		// log.Log.Warnf("height:%v,txIndex:%v,txHash:%v", blockInfo.Height, index, tx.Hash)
		if err := m.UpdateLocation(blockInfo.Height, tx.Txid); err != nil {
			log.Log.Warnf("UpdateLocation err:%v", err)
		}
		m.TraceTx(tx, blockInfo.Height)
	}
	log.Log.Warnf("time.Since(startTime):%v, height:%v", time.Since(startTime), blockInfo.Height)
}

func (m *Atomicals) TraceTx(tx btcjson.TxRawResult, height int64) error {
	operation := witness.ParseWitness(tx, height)
	// TODO:
	// get_if_parent_spent_in_same_tx
	//
	// step 1: transfer nft, transfer ft
	if err := m.transferNft(operation, tx); err != nil {
		log.Log.Warnf("transferNft err:%+v", err)
	}
	if err := m.transferFt(operation, tx); err != nil {
		log.Log.Warnf("transferFt err:%+v", err)
	}

	// step 2: process operation
	for _, vin := range tx.Vin {
		if operation.Op != "" {
			var err error
			operation.CommitHeight, err = m.btcClient.GetCommitHeight(operation.CommitTxID)
			if err != nil {
				log.Log.Warnf("GetCommitHeight err:%+v", err)
				// todo: retry,ensure success
			}
		}
		userPk := tx.Vout[common.VOUT_EXPECT_OUTPUT_INDEX].ScriptPubKey.Address
		if operation.Op == "dmt" {
			if err := m.mintDistributedFt(operation, vin, tx.Vout, userPk); err != nil {
				log.Log.Warnf("mintDistributedFt err:%+v", err)
			}
		} else {
			switch operation.Op {
			case "dft":
				if err := m.deployDistributedFt(operation, vin, tx.Vout, userPk); err != nil {
					log.Log.Warnf("deployDistributedFt err:%+v", err)
				}
			case "ft":
				if err := m.mintDirectFt(operation, vin, tx.Vout, userPk); err != nil {
					log.Log.Warnf("mintFt err:%+v", err)
				}
			case "nft":
				if err := m.mintNft(operation, vin, tx.Vout, userPk); err != nil {
					log.Log.Warnf("mintNft err:%+v", err)
				}
			case "mod":
			case "evt":
			case "dat":
			case "sl":
			default:
			}
		}
	}
	// step 3 check payment
	return nil
}

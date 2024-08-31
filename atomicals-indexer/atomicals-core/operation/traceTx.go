package atomicals

import (
	"time"

	"github.com/atomicals-go/atomicals-indexer/atomicals-core/witness"
	"github.com/atomicals-go/pkg/log"
	"github.com/atomicals-go/repo"
	"github.com/atomicals-go/utils"
	"github.com/btcsuite/btcd/btcjson"
)

func (m *Atomicals) Run() {
	var startTime1 time.Duration
	var startTime2 time.Duration

	startTime := time.Now()
	if m.location.BlockHeight+utils.SafeBlockHeightInterupt > m.maxBlockHeight {
		time.Sleep(10 * time.Minute)
		var err error
		m.maxBlockHeight, err = m.GetBlockCount()
		if err != nil {
			log.Log.Panicf("GetBlockCount err:%v", err)
		}
		return
	}
	block, err := m.GetBlockByHeightSync(m.location.BlockHeight)
	if err != nil {
		log.Log.Panicf("GetBlockByHeightSync err:%v", err)
	}
	m.location.TxIndex++
	if m.location.TxIndex >= int64(len(block.Tx)) {
		m.location.BlockHeight++
		m.location.TxIndex = 0
		block, err = m.GetBlockByHeightSync(m.location.BlockHeight)
		if err != nil {
			log.Log.Panicf("GetBlockByHeightSync err:%v", err)
		}
	}
	for txIndex, tx := range block.Tx {
		if int64(txIndex) < m.location.TxIndex {
			continue
		}
		startTime := time.Now()
		m.location.TxIndex = int64(txIndex)
		m.location.Txid = tx.Txid
		data := m.TraceTx(tx, block.Height)
		startTime1 = startTime1 + time.Since(startTime)

		startTime = time.Now()
		err = m.UpdateDB(m.location, data)
		if err != nil {
			log.Log.Panicf("UpdateDB err:%v", err)
		}
		startTime2 = startTime2 + time.Since(startTime)
		// log.Log.Infof("maxBlockHeight:%v, currentHeight:%v, %v lenTx:%v time:%v %v %v", m.maxBlockHeight, m.location.BlockHeight, len(block.Tx), txIndex, time.Since(startTime), startTime1, startTime2)
	}
	log.Log.Infof("maxBlockHeight:%v, currentHeight:%v,lenTx:%v time:%v %v %v", m.maxBlockHeight, m.location.BlockHeight, len(block.Tx), time.Since(startTime), startTime1, startTime2)
}

func (m *Atomicals) TraceTx(tx btcjson.TxRawResult, height int64) *repo.AtomicaslData {
	operation := witness.ParseWitness(tx, height)

	data := &repo.AtomicaslData{}

	// step 1: insert mod
	if operation.Op == "mod" {
		data.Mod = m.operationMod(operation, tx)
	}

	// step 2: transfer nft first, then transfer ft
	data.UpdateNfts, _ = m.transferNft(operation, tx)

	data.DeleteFts, data.NewFts, _ = m.transferFt(operation, tx)

	// step 3: process operation
	userPk := tx.Vout[utils.VOUT_EXPECT_OUTPUT_INDEX].ScriptPubKey.Address
	if operation.Op == "dmt" {
		data.NewUTXOFtInfo, data.UpdateDistributedFt, _ = m.mintDistributedFt(operation, tx.Vout, userPk)
	} else {
		switch operation.Op {
		case "dft":
			data.NewGlobalDistributedFt, _ = m.deployDistributedFt(operation, userPk)
		case "ft":
			data.NewGlobalDirectFt, _ = m.mintDirectFt(operation, tx.Vout, userPk)
		case "nft":
			if operation.Payload.Args.RequestContainer != "" || operation.Payload.Args.RequestDmitem != "" {
				data.NewUTXONftInfo, data.DeleteUTXONfts, _ = m.mintNft(operation, userPk)
			}
		case "evt":
			panic(operation.Payload)
		case "dat":
			data.Dat = m.operationDat(operation, tx)
		case "sl":
			panic(operation.Payload)
		default:
		}
	}

	// TODO: step 4 check payment

	data.ParseOperation(operation.Op)
	return data
}

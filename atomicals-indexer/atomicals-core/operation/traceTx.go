package atomicals

import (
	"strings"
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
		data := m.TraceTx(tx, block.Height)
		startTime1 = startTime1 + time.Since(startTime)

		startTime = time.Now()
		err = m.UpdateDB(block.Height, m.location.TxIndex, tx.Txid, data)
		if err != nil {
			log.Log.Panicf("UpdateDB err:%v", err)
		}
		startTime2 = startTime2 + time.Since(startTime)
	}

	log.Log.Infof("maxBlockHeight:%v, currentHeight:%v,lenTx:%v time:%v %v %v", m.maxBlockHeight, m.location.BlockHeight, len(block.Tx), time.Since(startTime), startTime1, startTime2)
}

func (m *Atomicals) TraceTx(tx btcjson.TxRawResult, height int64) *repo.AtomicaslData {
	operation := witness.ParseWitness(tx, height)
	data := &repo.AtomicaslData{}
	if operation.Op == "x" {
		data.Op = "splat|"
	} else if operation.Op == "y" {
		data.Op = "split|"
	}

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
		if data.NewUTXOFtInfo == nil {
			data.Op = "mint-dft-failed|"
		} else {
			data.Op = "mint-dft|"
		}
	} else {
		switch operation.Op {
		case "dft":
			data.NewGlobalDistributedFt, _ = m.deployDistributedFt(operation, userPk)
		case "ft":
			data.NewGlobalDirectFt, _ = m.mintDirectFt(operation, tx.Vout, userPk)
		case "nft":
			data.NewUTXONftInfo, data.DeleteUTXONfts, _ = m.mintNft(operation, userPk)
		case "evt":
			panic("")
		case "dat":
			panic("")
		case "sl":
			panic("")
		default:
		}
	}

	// TODO: step 4 check payment
	data.Op += operation.Op
	data.ParseOperation()
	if data.Op != "" && operation.Op != "dmt" {
		wizzop, err := fetchTxFromWizz(operation.RevealLocationTxID)
		if err != nil {
			log.Log.Infof("RevealLocationTxID:%v", operation.RevealLocationTxID)
			panic(err)
		}
		if wizzop != "" {
			log.Log.Infof("txid  %v %v, operation.Op: %v,wizz op: %v, data op: %v", operation.RevealLocationHeight, operation.RevealLocationTxID, operation.Op, wizzop, data.Op)
			if !strings.Contains(data.Op, wizzop) {
				panic("")
			}
			if tx.Txid == "054cc18a8162887917a1e6e5c60389bb4b6647167e6936d231466d7b2710f413" {
				panic("")
			}
		}
	}
	return data
}

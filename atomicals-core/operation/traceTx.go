package atomicals

import (
	"encoding/json"
	"time"

	"github.com/atomicals-go/atomicals-core/witness"
	"github.com/atomicals-go/pkg/log"
	"github.com/atomicals-go/repo/postsql"
	"github.com/atomicals-go/utils"
	"github.com/btcsuite/btcd/btcjson"
)

func (m *Atomicals) Run() {
	location, err := m.CurrentLocation()
	if err != nil {
		log.Log.Panicf("CurrentLocation err:%v", err)
	}
	tx := &btcjson.TxRawResult{}
	location.BlockHeight, location.TxIndex, tx = m.GetTxByHeightAndIndex(location.BlockHeight, location.TxIndex+1)
	m.traceTx(tx, location)

	if location.TxIndex == 0 {
		log.Log.Infof("height:%v, time:%v", location.BlockHeight, time.Since(m.startTime))
		m.startTime = time.Now()
	}
}

func (m *Atomicals) traceTx(tx *btcjson.TxRawResult, location *postsql.Location) error {
	// step 1: parse witness
	operation := witness.ParseWitness(tx, location.BlockHeight)

	// step 2: insert mod data
	if operation.Op == "mod" && len(tx.Vin) != 0 {
		vin := tx.Vin[0]
		preNftLocationID := utils.AtomicalsID(vin.Txid, int64(vin.Vout))
		preNfts, err := m.NftUTXOsByLocationID(preNftLocationID)
		if err != nil {
			log.Log.Panicf("NftUTXOsByLocationID err:%v", err)
		}
		if len(preNfts) != 0 {
			r, _ := json.Marshal(operation.Payload.Dmint)
			m.InsertMod(&postsql.ModInfo{
				Height:      location.BlockHeight,
				AtomicalsID: preNfts[0].AtomicalsID,
				LocationID:  preNfts[0].LocationID,
				Mod:         string(r),
				ModStr:      operation.PayloadStr,
			})
		}
	}

	// step 3: transfer nft, transfer ft
	if location.BlockHeight < utils.AtOMICALS_FT_PARTIAL_SPLITING_HEIGHT {
		m.transferFt(operation, tx)
	} else {
		m.transferFtPartialColour(operation, tx)
	}
	m.transferNft(operation, tx)

	// step 4: process operation
	userPk := tx.Vout[utils.VOUT_EXPECT_OUTPUT_INDEX].ScriptPubKey.Address
	if operation.Op == "dmt" {
		m.mintDistributedFt(operation, tx.Vout, userPk)
	} else {
		switch operation.Op {
		case "dft":
			m.deployDistributedFt(operation, userPk)
		case "ft":
			m.mintDirectFt(operation, tx.Vout, userPk)
		case "nft":
			m.mintNft(operation, userPk, location.BlockHeight)
		case "evt":
		case "dat":
		case "sl":
		default:
		}
	}

	// step 5 check payment

	// step 6: exec sql
	if err := m.ExecAllSql(&postsql.Location{
		BlockHeight: location.BlockHeight,
		TxIndex:     location.TxIndex,
		Txid:        tx.Txid,
	}); err != nil {
		log.Log.Panicf("ExecAllSql err:%v", err)
	}
	return nil
}

func (m *Atomicals) TraceSpecificTx() {
	// height := int64(812547)
	// index := int64(5)
	// for {
	// 	tx := &btcjson.TxRawResult{}
	// 	height, index, tx = m.GetTxByHeightAndIndex(height, index)
	// 	log.Log.Infof("height: %v txIndex: %v Txid: %v", height, index, tx.Txid)
	// 	if tx.Txid == "4211d0c9b069f1c9624b9616c6ea0c0c548d8beceede393c938d09eb4e971a47" {
	// 		witness.ParseWitness(tx, height)
	// 		panic("")
	// 	}
	// 	index += 1
	// }
}

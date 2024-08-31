package atomicals

import (
	"encoding/json"

	"github.com/atomicals-go/atomicals-indexer/atomicals-core/witness"
	"github.com/atomicals-go/pkg/log"
	"github.com/atomicals-go/repo/postsql"
	"github.com/atomicals-go/utils"
	"github.com/btcsuite/btcd/btcjson"
)

func (m *Atomicals) operationMod(operation *witness.WitnessAtomicalsOperation, tx btcjson.TxRawResult) *postsql.ModInfo {
	if operation.Op != "mod" || len(tx.Vin) == 0 {
		return nil
	}
	var preNfts []*postsql.UTXONftInfo
	var err error
	for _, vin := range tx.Vin {
		preNftLocationID := utils.AtomicalsID(vin.Txid, int64(vin.Vout))
		preNfts, err = m.NftUTXOsByLocationID(preNftLocationID)
		if err != nil {
			log.Log.Panicf("NftUTXOsByLocationID err:%v", err)
		}
	}
	if len(preNfts) == 0 {
		return nil
	}
	r, err := json.Marshal(operation.Payload.Dmint)
	if err != nil {
		log.Log.Panicf("Marshal err:%v", err)
	}
	return &postsql.ModInfo{
		Height:      operation.RevealLocationHeight,
		AtomicalsID: preNfts[0].AtomicalsID,
		LocationID:  preNfts[0].LocationID,
		Mod:         string(r),
		// ModStr:      operation.PayloadStr,
	}
}

package atomicals

import (
	"github.com/atomicals-core/pkg/log"

	"github.com/atomicals-core/atomicals/common"
	"github.com/atomicals-core/atomicals/witness"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

func (m *Atomicals) TraceBlock(blockInfo *btcjson.GetBlockVerboseTxResult) {
	for index, tx := range blockInfo.Tx {
		log.Log.Warnf("height:%v,txIndex:%v,txHash:%v", blockInfo.Height, index, tx.Hash)
		m.TraceTx(tx, blockInfo.Height)
	}
	m.Height++
}

func (m *Atomicals) TraceTx(tx btcjson.TxRawResult, height int64) error {
	operation := witness.ParseWitness(tx, height)
	// step 1: mint or deploy
	// TODO:
	// get_if_parent_spent_in_same_tx
	//
	if err := m.getCommitHeight(operation); err != nil {
		log.Log.Warnf("getCommitHeight err:%+v", err)
		// todo: retry,ensure success
	}
	for _, vin := range tx.Vin {
		userPk := tx.Vout[common.VOUT_EXPECT_OUTPUT_INDEX].ScriptPubKey.Address
		switch operation.Op {
		case "dft":
			if err := m.deployDistributedFt(operation, vin, tx.Vout, userPk); err != nil {
				log.Log.Warnf("deployDistributedFt err:%+v", err)
			}
		case "dmt":
			if err := m.mintDistributedFt(operation, vin, tx.Vout, userPk); err != nil {
				log.Log.Warnf("mintDistributedFt err:%+v", err)
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
	// step 2: transfer nft
	if err := m.transferNft(operation, tx); err != nil {
		log.Log.Warnf("transferNft err:%+v", err)
	}
	// step 3: transfer ft
	if err := m.transferFt(operation, tx); err != nil {
		log.Log.Warnf("transferFt err:%+v", err)
	}
	return nil
}

func (m *Atomicals) getCommitHeight(operation *witness.WitnessAtomicalsOperation) error {
	t, err := m.btcClient.GetTransaction(operation.CommitTxID)
	if err != nil {
		log.Log.Warnf("GetTransaction err:%v", err)
	}
	blockHash, err := chainhash.NewHashFromStr(t.BlockHash)
	if err != nil {
		return err
	}
	blockInfo, err := m.btcClient.GetBlockVerboseTx(blockHash)
	if err != nil {
		return err
	}
	operation.CommitHeight = blockInfo.Height
	return nil
}

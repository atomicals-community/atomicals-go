package btcsync

import "encoding/hex"

func TxID2txHash(txID string) (string, error) {
	txidBytes, err := hex.DecodeString(txID)
	if err != nil {
		return "", err
	}
	txhash := hex.EncodeToString(txidBytes)
	return txhash, nil
}

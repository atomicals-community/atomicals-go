package btcsync

import (
	"encoding/hex"
	"fmt"
)

func (m *BtcSync) GetPublicKeyFromAddress(address string) (string, error) {
	// Get address info
	addrInfo, err := m.GetAddressInfo(address)
	if err != nil {
		return "", fmt.Errorf("error getting address info: %v", err)
	}
	// Decode public key
	pubKeyBytes, err := hex.DecodeString(addrInfo.ScriptPubKey)
	if err != nil {
		return "", fmt.Errorf("error decoding public key: %v", err)
	}
	return hex.EncodeToString(pubKeyBytes), nil
}

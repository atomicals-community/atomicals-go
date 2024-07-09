package utils

import (
	"bytes"
	"encoding/hex"
)

// is_unspendable_legacy
func IsUnspendableLegacy(script []byte) bool {
	// OP_FALSE OP_RETURN or OP_RETURN
	return bytes.Equal(script[:2], []byte{0x00, 0x6a}) || (len(script) > 0 && script[0] == 0x6a)
}

// is_unspendable_genesis
func IsUnspendableGenesis(script []byte) bool {
	// OP_FALSE OP_RETURN
	return bytes.Equal(script[:2], []byte{0x00, 0x6a})
}

func Is_op_return_subrealm_payment_marker_atomical_id(script []byte) string {
	if len(script) < 1+5+2+1+36 { // 6a04<atom><01>p<atomical_id>
		return ""
	}
	// Ensure it is an OP_RETURN
	firstByte := script[0]
	secondBytes := script[:2]
	if secondBytes[0] != 0x00 && firstByte != 0x6a {
		return ""
	}
	startIndex := 1
	if secondBytes[0] == 0x00 {
		startIndex = 2
	}
	// Check for the envelope format
	if hex.EncodeToString(script[startIndex:startIndex+5]) != ATOMICALS_ENVELOPE_MARKER_BYTES {
		return ""
	}
	// Check the next op code matches 'p' for payment
	if hex.EncodeToString(script[startIndex+5:startIndex+5+2]) != "0170" {
		return ""
	}
	// Check there is a 36 byte push data
	if hex.EncodeToString(script[startIndex+5+2:startIndex+5+2+1]) != "24" {
		return ""
	}
	// Extract and return the atomical ID
	atomicalID := script[startIndex+5+2+1 : startIndex+5+2+1+36]
	return hex.EncodeToString(atomicalID)
}
func Is_op_return_dmitem_payment_marker_atomical_id(script []byte) string {
	if len(script) < 1+5+2+1+36 { // 6a04<atom><01>p<atomical_id>
		return ""
	}
	// Ensure it is an OP_RETURN
	firstByte := script[0]
	secondBytes := script[:2]
	if secondBytes[0] != 0x00 && firstByte != 0x6a {
		return ""
	}
	startIndex := 1
	if secondBytes[0] == 0x00 {
		startIndex = 2
	}
	// Check for the envelope format
	if hex.EncodeToString(script[startIndex:startIndex+5]) != ATOMICALS_ENVELOPE_MARKER_BYTES {
		return ""
	}
	// Check the next op code matches 'p' for payment
	if hex.EncodeToString(script[startIndex+5:startIndex+5+2]) != "0164" {
		return ""
	}
	// Check there is a 36 byte push data
	if hex.EncodeToString(script[startIndex+5+2:startIndex+5+2+1]) != "24" {
		return ""
	}
	// Extract and return the atomical ID
	atomicalID := script[startIndex+5+2+1 : startIndex+5+2+1+36]
	return hex.EncodeToString(atomicalID)
}

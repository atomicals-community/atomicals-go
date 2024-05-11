package utils

import "bytes"

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

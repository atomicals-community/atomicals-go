package common

import (
	"encoding/hex"
)

func IsHexString(value string) bool {
	// Try to decode the string as hexadecimal
	_, err := hex.DecodeString(value)
	if err != nil {
		return false
	}
	return true
}

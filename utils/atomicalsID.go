package utils

import (
	"encoding/hex"
	"fmt"
	"strings"
)

func AtomicalsID(txID string, voutIndex int64) string {
	return fmt.Sprintf("%vi%v", txID, voutIndex)
}

func SplitAtomicalsID(atomicalsID string) (string, string) {
	parts := strings.SplitN(atomicalsID, "i", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return "", ""
}

// is_compact_atomical_id
func IsCompactAtomicalID(value string) bool {
	// Check if the value is an integer
	// Go doesn't have a built-in 'isinstance' function like Python, so we check if we can convert it to an integer
	var intValue int
	_, err := fmt.Sscan(value, &intValue)
	if err == nil {
		return false
	}

	// Check if the value is empty or None
	if value == "" {
		return false
	}

	// Check if the length is at least 64 characters and the 64th character is 'i'
	if len(value) < 64 || value[63] != 'i' {
		return false
	}

	// Extract the raw hash part and convert it to bytes
	rawHashHex := value[:64]
	rawHash, err := hex.DecodeString(rawHashHex)
	if err != nil {
		return false
	}

	// Check if the raw hash has a length of 32 bytes
	return len(rawHash) == 32
}

package utils

import (
	"encoding/hex"
	"errors"
	"fmt"
)

// Convert the compact string form to the expanded 36 byte sequence
func compact_to_location_id_bytes(value string) ([]byte, error) {
	if value == "" {
		return nil, errors.New("value in compactToLocationIDBytes is not set")
	}

	indexOfI := len(value) - 1
	if indexOfI != 64 || value[indexOfI] != 'i' {
		return nil, fmt.Errorf("%s should be 32 bytes hex followed by i<number>", value)
	}

	rawHash, err := hex.DecodeString(value[:64])
	if err != nil {
		return nil, err
	}

	if len(rawHash) != 32 {
		return nil, fmt.Errorf("%s should be 32 bytes hex followed by i<number>", value)
	}

	num, err := parseLEUint32(value[65:])
	if err != nil {
		return nil, err
	}

	if num < 0 || num > 100000 {
		return nil, fmt.Errorf("%s index output number was parsed to be less than 0 or greater than 100000", value)
	}
	return append(rawHash, packLEUint32(num)...), nil
}

func is_op_return_subrealm_payment_marker_atomical_id(script []byte) string {
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
func is_op_return_dmitem_payment_marker_atomical_id(script []byte) string {
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

package common

import (
	"crypto/hmac"
	sha256p "crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"strings"
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

// reverseHex reverses a hex
func reverseHex(input string) []byte {
	hexBytes, _ := hex.DecodeString(input)
	for i, j := 0, len(hexBytes)-1; i < j; i, j = i+1, j-1 {
		hexBytes[i], hexBytes[j] = hexBytes[j], hexBytes[i]
	}
	return hexBytes
}

func sha256(data []byte) []byte {
	hash := sha256p.New()
	hash.Write(data)
	return hash.Sum(nil)
}

func hmacDigest(key, data []byte) []byte {
	hash := hmac.New(sha256p.New, key)
	hash.Write(data)
	return hash.Sum(nil)
}

func hexToBytes(hexString string) ([]byte, error) {
	return hex.DecodeString(hexString)
}

func header_hash(header []byte) []byte {
	// ”'Given a header return hash”'
	return double_sha256(header)
}

func double_sha256(x []byte) []byte {
	// '''SHA-256 of SHA-256, as used extensively in bitcoin.'''
	return sha256(sha256(x))
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

// func create_or_delete_subname_payment_output_if_valid(atomical_id_for_payment string, payment_marker_idx int, blueprint_builder *AtomicalsTransferBlueprintBuilder, tx *btcjson.TxRawResult, height int64, db_prefix string) string {
// if blueprint_builder.IsSplitOperation(){
// 	return nil
// }
// matched_price_point, parent_id, request_subname, subname_type := get_expected_subname_payment_info(atomical_id_for_payment, height)
// // # An expected payment amount might ! be set if there is no valid subrealm minting rules, or something invalid was found
// if ! matched_price_point{
// 	return nil
// }
// regex := matched_price_point["matched_rule"]["p"]
// if ! is_valid_regex(regex){
// 	return nil
// }
// // # The pattern should have already matched, sanity check
// valid_pattern = re.compile(rf"{regex}")
// if ! valid_pattern.match(request_subname){
// 	panic("err")
// }
// if ! blueprint_builder.are_payments_satisfied(matched_price_point["matched_rule"].get("o")):
// 	return nil
// // # Delete or create the record based on whether we are reorg rollback or creating new
// payment_outpoint = tx_hash + pack_le_uint32(payment_marker_idx)
// not_initated_by_parent = "00" # Used to indicate it was minted according to rules payment match
// self.put_pay_record(atomical_id_for_payment, tx_num, payment_outpoint + not_initated_by_parent, db_prefix)
// return tx_hash
// 	return ""
// }

func get_adjusted_sats_needed_by_exponent(value float64, exponent int64) float64 {
	return (float64(value) / math.Pow(10, float64(exponent)))
}

func get_nominal_token_value(value float64, exponent int64) float64 {
	if value < 0 || exponent < 0 {
		panic("Value and exponent must be non-negative")
	}
	return value / math.Pow10(int(exponent))
}

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

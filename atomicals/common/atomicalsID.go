package common

import (
	"encoding/hex"
	"errors"
	"fmt"
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

package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"math"
)

// reverseHex reverses a hex
func reverseHex(input string) []byte {
	hexBytes, _ := hex.DecodeString(input)
	for i, j := 0, len(hexBytes)-1; i < j; i, j = i+1, j-1 {
		hexBytes[i], hexBytes[j] = hexBytes[j], hexBytes[i]
	}
	return hexBytes
}

func hmacDigest(key, data []byte) []byte {
	hash := hmac.New(sha256.New, key)
	hash.Write(data)
	return hash.Sum(nil)
}

func hexToBytes(hexString string) ([]byte, error) {
	return hex.DecodeString(hexString)
}

func header_hash(header []byte) []byte {
	// ”'Given a header return hash”'
	return DoubleSha256(header)
}

func Sha256(data []byte) []byte {
	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)
}
func DoubleSha256(x []byte) []byte {
	// '''SHA-256 of SHA-256, as used extensively in bitcoin.'''
	return Sha256(Sha256(x))
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

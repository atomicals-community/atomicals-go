package common

import (
	"crypto/hmac"
	sha256p "crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
)

func Get_next_bitwork_full_str(bitworkVec string, currentPrefixLen int) string {
	baseBitworkPadded := fmt.Sprintf("%-32s", bitworkVec)
	if currentPrefixLen >= 31 {
		return baseBitworkPadded
	}
	return baseBitworkPadded[:currentPrefixLen+1]
}

func Is_mint_pow_valid(txid, mintPowCommit string) bool {
	bitworkCommitParts := ParseBitwork(mintPowCommit)
	if bitworkCommitParts == nil {
		return false
	}
	mintBitworkPrefix := bitworkCommitParts.Prefix
	mintBitworkExt := bitworkCommitParts.Ext
	if IsProofOfWorkPrefixMatch(txid, mintBitworkPrefix, mintBitworkExt) {
		return true
	}
	return false
}

func isIntInRange(value, min, max int) bool {
	return value >= min && value <= max
}

func Calculate_expected_bitwork(bitwork_vec string, actual_mints, max_mints, target_increment, starting_target int64) string {
	if starting_target < 64 || starting_target > 256 {
		panic("err")
	}
	if max_mints < 1 || max_mints > 100000 {
		panic("err")
	}
	if target_increment < 1 || target_increment > 64 {
		panic("err")
	}
	target_steps := (actual_mints) / (max_mints)
	current_target := starting_target + (target_steps * target_increment)
	return derive_bitwork_prefix_from_target(bitwork_vec, current_target)
}

func derive_bitwork_prefix_from_target(baseBitworkPrefix string, target int64) string {
	if target < 16 {
		panic(fmt.Sprintf("increments must be at least 16. Provided: %d", target))
	}

	baseBitworkPadded := fmt.Sprintf("%-32s", baseBitworkPrefix)
	multiples := float64(target) / 16
	fullAmount := int(math.Floor(multiples))
	modulo := target % 16

	bitworkPrefix := baseBitworkPadded
	if fullAmount < 32 {
		bitworkPrefix = baseBitworkPadded[:fullAmount]
	}

	if modulo > 0 {
		return bitworkPrefix + "." + fmt.Sprint(modulo)
	}

	return bitworkPrefix
}

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

func parseLEUint32(s string) (uint32, error) {
	bytes, err := hex.DecodeString(s)
	if err != nil {
		return 0, err
	}
	if len(bytes) != 4 {
		return 0, errors.New("invalid length for LE uint32")
	}
	return binary.LittleEndian.Uint32(bytes), nil
}

func packLEUint32(num uint32) []byte {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, num)
	return bytes
}

// Base atomical commit to reveal delay allowed
func is_within_acceptable_blocks_for_general_reveal(commitHeight, revealLocationHeight int64) bool {
	return commitHeight >= revealLocationHeight-MINT_GENERAL_COMMIT_REVEAL_DELAY_BLOCKS
}

// A realm, ticker, or container reveal is valid as long as it is within MINT_REALM_CONTAINER_TICKER_COMMIT_REVEAL_DELAY_BLOCKS of the reveal and commit
func is_within_acceptable_blocks_for_name_reveal(commitHeight, revealLocationHeight int64) bool {
	return commitHeight >= revealLocationHeight-MINT_REALM_CONTAINER_TICKER_COMMIT_REVEAL_DELAY_BLOCKS
}

// A payment for a subrealm is acceptable as long as it is within MINT_SUBNAME_COMMIT_PAYMENT_DELAY_BLOCKS of the commitHeight
func isWithinAcceptableBlocksForSubItemPayment(commitHeight, currentHeight int64) bool {
	return currentHeight <= commitHeight+MINT_SUBNAME_COMMIT_PAYMENT_DELAY_BLOCKS
}

// Encoder struct
type Encoder struct {
	data []byte
}

// NewEncoder creates a new Encoder
func NewEncoder() *Encoder {
	return &Encoder{}
}

// Update updates the encoder with the given data
func (e *Encoder) Update(data []byte) {
	e.data = append(e.data, data...)
}

// Finalize finalizes the encoding and returns the result as a hex string
func (e *Encoder) Finalize() string {
	return hex.EncodeToString(e.data)
}

// reverseHex reverses a hex
func reverseHex(input string) []byte {
	hexBytes, _ := hex.DecodeString(input)
	for i, j := 0, len(hexBytes)-1; i < j; i, j = i+1, j-1 {
		hexBytes[i], hexBytes[j] = hexBytes[j], hexBytes[i]
	}
	return hexBytes
}

func is_density_activated(height int64) bool {
	if height >= ATOMICALS_ACTIVATION_HEIGHT_DENSITY {
		return true
	}
	return false
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

func readfile() string {
	filePath := "./witnessscript.txt"
	// Read the content of the file
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return ""
	}

	// Convert content to a string
	longString := string(content)
	return longString
}

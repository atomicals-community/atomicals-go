package witness

import "encoding/hex"

// parse_operation_from_script
func parseAtomicalsOperation(scriptBytes []byte, startIndex int64) (string, int64) {
	one_letter_op_len := int64(2)
	two_letter_op_len := int64(3)
	three_letter_op_len := int64(4)
	operationType := ""

	// # check the 3 letter protocol operations
	atomOp := scriptBytes[startIndex : startIndex+three_letter_op_len]
	switch hex.EncodeToString(atomOp) {
	case "036e6674":
		operationType = "nft" // nft - Mint non-fungible token
	case "03646674":
		operationType = "dft" // dft - Deploy distributed mint fungible token starting point
	case "036d6f64":
		operationType = "mod" // mod - Modify general state
	case "03657674":
		operationType = "evt" // evt - Message response/reply
	case "03646d74":
		operationType = "dmt" // dmt - Mint tokens of distributed mint type (dft)
	case "03646174":
		operationType = "dat" // dat - Store data on a transaction (dat)
	}
	if operationType != "" {
		return operationType, startIndex + three_letter_op_len
	}
	// # check the 2 letter protocol operations
	atomOp = scriptBytes[startIndex : startIndex+two_letter_op_len]
	switch hex.EncodeToString(atomOp) {
	case "026674":
		operationType = "ft" //# ft - Mint fungible token with direct fixed supply
	case "02736c":
		operationType = "sl" //# sl - Seal an NFT and lock it from further changes forever
	}
	if operationType != "" {
		return operationType, startIndex + two_letter_op_len
	}
	// # check the 1 letter protocol operations
	atomOp = scriptBytes[startIndex : startIndex+one_letter_op_len]
	switch hex.EncodeToString(atomOp) {
	case "0178":
		operationType = "x" //# extract - move atomical to 0'th output
	case "0179":
		operationType = "y" //# split -
	}
	if operationType != "" {
		return operationType, startIndex + one_letter_op_len
	}
	return operationType, -1
}

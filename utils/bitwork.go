package utils

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

func ParseMintBitwork(commitTxID, mintBitworkc, mintBitworkr string) (*Bitwork, *Bitwork, error) {
	bitworkc := ParseBitwork(mintBitworkc)
	bitworkr := ParseBitwork(mintBitworkr)
	return bitworkc, bitworkr, nil
}

// is_proof_of_work_prefix_match
func IsProofOfWorkPrefixMatch(txID string, powPrefix string, powPrefixExt int) bool {
	if powPrefixExt < 0 {
		return strings.HasPrefix(txID, powPrefix)
	}
	if powPrefixExt < 0 || powPrefixExt > 15 {
		return false
	}

	// Check that the main prefix matches
	initialTestMatchesMainPrefix := strings.HasPrefix(txID, powPrefix)
	if !initialTestMatchesMainPrefix {
		return false
	}

	// If there is an extended powPrefixExt, then we require it to validate the POW
	if powPrefixExt > 0 {
		// Check that the next digit is within the range of powPrefixExt
		nextChar := txID[len(powPrefix)]
		charMap := map[rune]int{
			'0': 0, '1': 1, '2': 2, '3': 3, '4': 4,
			'5': 5, '6': 6, '7': 7, '8': 8, '9': 9,
			'a': 10, 'b': 11, 'c': 12, 'd': 13, 'e': 14, 'f': 15,
		}
		getNumericValue := charMap[rune(nextChar)]

		// powPrefixExt == 0 is functionally equivalent to not having a powPrefixExt
		// powPrefixExt == 15 is functionally equivalent to extending the powPrefix by 1
		return getNumericValue >= powPrefixExt
	}

	// There is no extended powPrefixExt, and we just apply the main prefix
	return true
}

// # Parse a bitwork stirng such as '123af.15'
type Bitwork struct {
	Prefix string
	Ext    int
}

// is_valid_bitwork_string
func ParseBitwork(bitwork string) *Bitwork {
	if bitwork == "" {
		return nil
	}
	if strings.Count(bitwork, ".") > 1 {
		return nil
	}
	splitted := strings.Split(bitwork, ".")
	prefix := splitted[0]
	ext := -1
	if len(splitted) > 1 {
		extStr := splitted[1]
		extInt, err := strconv.Atoi(extStr)
		if err != nil {
			return nil
		}
		ext = extInt
	}
	if prefix == "" {
		return nil
	}
	if !regexp.MustCompile("^[a-f0-9]{1,64}$").MatchString(prefix) {
		return nil
	}
	if ext > 15 {
		return nil
	}
	return &Bitwork{
		Prefix: prefix,
		Ext:    ext,
	}
}

func GetNextBitworkFullStr(bitworkVec string, currentPrefixLen int) string {
	baseBitworkPadded := fmt.Sprintf("%-32s", bitworkVec)
	if currentPrefixLen >= 31 {
		return baseBitworkPadded
	}
	return baseBitworkPadded[:currentPrefixLen+1]
}

func IsMintPowValid(txid, mintPowCommit string) bool {
	bitworkCommitParts := ParseBitwork(mintPowCommit)
	if bitworkCommitParts == nil {
		return false
	}
	mintBitworkPrefix := bitworkCommitParts.Prefix
	mintBitworkExt := bitworkCommitParts.Ext
	return IsProofOfWorkPrefixMatch(txid, mintBitworkPrefix, mintBitworkExt)
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

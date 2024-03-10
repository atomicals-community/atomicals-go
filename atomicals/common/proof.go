package common

import (
	"regexp"
	"strconv"
	"strings"
)

// is_proof_of_work_prefix_match
func IsProofOfWorkPrefixMatch(txID string, powPrefix string, powPrefixExt int) bool {
	if powPrefixExt > 0 && (powPrefixExt < 1 || powPrefixExt > 15) {
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
	if ext < 0 || ext > 15 {
		return nil
	}
	return &Bitwork{
		Prefix: prefix,
		Ext:    ext,
	}
}

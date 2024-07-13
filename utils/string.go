package utils

import (
	"encoding/hex"
	"regexp"
	"strings"
)

func IsHexString(value string) bool {
	// Try to decode the string as hexadecimal
	_, err := hex.DecodeString(value)
	if err != nil {
		return false
	}
	return true
}

func IsValidRegex(regex string) bool {
	if regex == "" {
		return false
	}

	if strings.ContainsAny(regex, "()") {
		return false
	}

	_, err := regexp.Compile(regex)
	if err != nil {
		return false
	}

	return true
}

func CompileRegex(pattern string) (*regexp.Regexp, error) {
	// Compile the regex pattern
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return regex, nil
}

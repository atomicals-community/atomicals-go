package utils

import (
	"regexp"
	"strings"
)

// is_valid_ticker_string
func IsValidTicker(ticker string) bool {
	if ticker == "" {
		return false
	}
	m := regexp.MustCompile(`^[a-z0-9]{1,21}$`)
	return m.MatchString(ticker)
}

// is_valid_realm_string_name
func IsValidRealm(realmName string) bool {
	if realmName == "" {
		return false
	}
	if strings.HasPrefix(realmName, "-") {
		return false
	}
	if strings.HasSuffix(realmName, "-") {
		return false
	}
	if len(realmName) > 64 || len(realmName) <= 0 {
		return false
	}
	// # Realm names must start with an alphabetical character
	m := regexp.MustCompile(`^[a-z][a-z0-9\-]{0,63}$`)
	return m.MatchString(realmName)
}

// is_valid_subrealm_string_name
func IsValidSubRealm(realmName string) bool {
	if realmName == "" {
		return false
	}
	if strings.HasPrefix(realmName, "-") {
		return false
	}
	if strings.HasSuffix(realmName, "-") {
		return false
	}
	if len(realmName) > 64 || len(realmName) <= 0 {
		return false
	}
	// # Realm names must start with an alphabetical character
	m := regexp.MustCompile(`^[a-z0-9][a-z0-9\-]{0,63}$`)
	return m.MatchString(realmName)
}

// is_valid_container_string_name
func IsValidContainer(containerName string) bool {
	if containerName == "" {
		return false
	}
	if strings.HasPrefix(containerName, "-") {
		return false
	}
	if strings.HasSuffix(containerName, "-") {
		return false
	}
	if len(containerName) > 64 || len(containerName) <= 0 {
		return false
	}
	// # Realm names must start with an alphabetical character
	m := regexp.MustCompile(`^[a-z0-9][a-z0-9\-]{0,63}$`)
	return m.MatchString(containerName)
}

// is_valid_container_dmitem_string_name
func IsValidDmitem(dmitemName string) bool {
	if dmitemName == "" {
		return false
	}
	if strings.HasPrefix(dmitemName, "-") {
		return false
	}
	if strings.HasSuffix(dmitemName, "-") {
		return false
	}
	if len(dmitemName) > 64 || len(dmitemName) <= 0 {
		return false
	}
	m := regexp.MustCompile(`^[a-z0-9][a-z0-9\-]{0,63}$`)
	return m.MatchString(dmitemName)
}

// is_hex_string_regex
func IsHexStringRegex(value string) bool {
	m := regexp.MustCompile(`^[a-z0-9]+$`)
	return m.MatchString(value)

}

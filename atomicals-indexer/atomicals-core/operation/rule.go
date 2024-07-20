package atomicals

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/atomicals-go/atomicals-indexer/atomicals-core/witness"
	"github.com/atomicals-go/utils"
)

func (m *Atomicals) checkRule(rule *witness.RuleInfo, bitworkc_actual, bitworkr_actual string) bool {
	if rule == nil {
		return false
	}
	bitworkc := rule.Bitworkc
	bitworkr := rule.Bitworkr

	if bitworkc == "any" {
		return true
	} else {
		if bitworkc_actual != bitworkc {
			return false
		}
	}
	if bitworkr == "any" {
		return true
	} else {
		if bitworkr_actual != bitworkr {
			return false
		}
	}
	if rule.O != nil {
		return true
	}
	if bitworkc != "" || bitworkr != "" {
		return true
	}
	return false
}
func (m *Atomicals) get_applicable_rule_by_height(parent_atomical_id string, proposed_subnameid string, height int64) *witness.RuleInfo {
	latest_state, err := m.getModHistory(parent_atomical_id, height)
	if err != nil {
		panic(err)
	}
	regex_price_point_list := validateRulesData(latest_state.Rules)
	for _, regex_price_point := range regex_price_point_list {
		valid_pattern := regexp.MustCompile(regex_price_point.P)
		if !valid_pattern.MatchString(proposed_subnameid) {
			continue
		}
		return regex_price_point
	}
	return nil
}

// get_mod_history
func (m *Atomicals) getModHistory(parentContainerAtomicalsID string, height int64) (*witness.Mod, error) {
	mods, err := m.ModHistory(parentContainerAtomicalsID, height)
	if err != nil {
		return nil, err
	}
	if mods == nil {
		return nil, errors.New("invalid mod")
	}
	dmints := make([]*witness.Mod, 0)
	for _, mod := range mods {
		dmint := &witness.Mod{}
		if err := json.Unmarshal([]byte(mod.Mod), dmint); err != nil {
			return nil, err
		}
		dmints = append(dmints, dmint)
	}
	// calculate_latest_state_from_mod_history
	// Ensure it is sorted in ascending order
	// sort.Slice(mod_history, func(i, j int) bool {
	// 	return mod_history[i].ID < mod_history[j].ID
	// })
	current_object_state := &witness.Mod{}
	for _, element := range dmints {
		if element.A == 1 {
			current_object_state = nil
		} else {
			current_object_state = element
		}
	}
	if current_object_state == nil {
		return nil, errors.New("invalid mod")
	}
	if validateRulesData(current_object_state.Rules) == nil {
		return nil, errors.New("invalid mod")
	}
	if current_object_state.MintHeight < 0 {
		return nil, errors.New("invalid mod")
	}
	if current_object_state.V != "1" {
		return nil, errors.New("invalid mod")
	}
	if len(current_object_state.Merkle) != 64 {
		return nil, errors.New("invalid mod")
	}
	return current_object_state, nil
}

func validateRulesData(rules []*witness.RuleInfo) []*witness.RuleInfo {
	if len(rules) <= 0 || len(rules) > utils.MAX_SUBNAME_RULE_ENTRIES {
		return nil
	}
	validated_rules_list := []*witness.RuleInfo{}
	for _, rule := range rules {
		regex_pattern := rule.P
		if len(regex_pattern) > utils.MAX_SUBNAME_RULE_SIZE_LEN || len(regex_pattern) < 1 {
			return nil
		}
		if strings.ContainsAny(regex_pattern, "()") {
			return nil
		}
		_, err := regexp.Compile(regex_pattern)
		if err != nil {
			fmt.Println("Regex compile error:", err)
			return nil
		}
		bitworkc := rule.Bitworkc
		bitworkr := rule.Bitworkr
		if regex_pattern == "" {
			return nil
		}
		if strings.Contains(regex_pattern, "(") || strings.Contains(regex_pattern, ")") {
			return nil
		}
		price_point := &witness.RuleInfo{
			P: regex_pattern,
		}
		if bitworkc != "" {
			res := utils.ParseBitwork(bitworkc)
			if res != nil {
				price_point.Bitworkc = bitworkc
			} else if bitworkc == "any" {
				price_point.Bitworkc = bitworkc
			} else {
				return nil
			}
		}
		if bitworkr != "" {
			res := utils.ParseBitwork(bitworkr)
			if res != nil {
				price_point.Bitworkr = bitworkr
			} else if bitworkr == "any" {
				price_point.Bitworkr = bitworkr
			} else {
				return nil
			}
		}
		if len(rule.O) > 0 {
			if !validate_subrealm_rules_outputs_format(rule.O) {
				return nil
			}
			price_point.O = rule.O
			validated_rules_list = append(validated_rules_list, price_point)
		} else if bitworkc != "" || bitworkr != "" {
			validated_rules_list = append(validated_rules_list, price_point)
		} else {
			return nil
		}
		if rule.O == nil && bitworkc == "" && bitworkr == "" {
			return nil
		}
	}
	if validated_rules_list == nil {
		return nil
	}
	if len(validated_rules_list) == 0 {
		return nil
	}
	return validated_rules_list
}

func validate_subrealm_rules_outputs_format(outputs map[string]*witness.Output) bool {
	for expected_output_script, expected_output_value := range outputs {
		expected_output_id := expected_output_value.ID
		expected_output_qty := expected_output_value.V
		if expected_output_qty < utils.SUBNAME_MIN_PAYMENT_DUST_LIMIT {
			return false // # Reject if one of the entries expects less than the minimum payment amount
		}
		// # If there is a type restriction on the payment type then ensure it is a valid atomical id
		if expected_output_id != "" {
			if utils.IsCompactAtomicalID(expected_output_id) {
				return false
			}
		}
		// # script must be paid to mint a subrealm
		if !utils.IsHexString(expected_output_script) {
			return false // # Reject if one of the payment output script is not a valid hex
		}
	}
	return true
}

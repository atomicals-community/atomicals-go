package atomicals

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/atomicals-go/atomicals-core/witness"
	"github.com/atomicals-go/pkg/merkle"
	"github.com/atomicals-go/utils"
)

// get_dmitem_parent_container_info
func (m *Atomicals) verifyRuleAndMerkle(operation *witness.WitnessAtomicalsOperation, height int64) bool {
	dmint_validated_status := m.get_container_dmint_status_for_atomical_id(operation.Payload.Args.ParentContainer)
	if dmint_validated_status == nil {
		return false
	}
	if operation.CommitHeight < dmint_validated_status.MintHeight || height < dmint_validated_status.MintHeight {
		return false
	}
	is_proof_valid := validate_dmitem_mint_args_with_container_dmint(operation.Payload, dmint_validated_status)
	if !is_proof_valid {
		return false
	}
	matched_price_point := m.get_applicable_rule_by_height(operation.Payload.Args.ParentContainer, operation.Payload.Args.RequestDmitem)
	bitworkc := matched_price_point.Bitworkc
	bitworkr := matched_price_point.Bitworkr
	bitworkc_actual := operation.Payload.Args.Bitworkc
	bitworkr_actual := operation.Payload.Args.Bitworkr
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
	if matched_price_point.O != nil {
		return true
	}
	if bitworkc != "" || bitworkr != "" {
		return true
	}
	// todo: put_pay_record
	return false
}

func (m *Atomicals) get_applicable_rule_by_height(parent_atomical_id string, proposed_subnameid string) *witness.RuleInfo {
	rule_mint_mod_history, err := m.get_mod_history(parent_atomical_id)
	if err != nil {
		panic(err)
	}
	latest_state := calculate_latest_state_from_mod_history(rule_mint_mod_history)
	regex_price_point_list := validate_rules_data(latest_state)
	for _, regex_price_point := range regex_price_point_list {
		valid_pattern := regexp.MustCompile(regex_price_point.P)
		if !valid_pattern.MatchString(proposed_subnameid) {
			continue
		}
		return regex_price_point
	}
	return nil
}

func (m *Atomicals) get_container_dmint_status_for_atomical_id(atomical_id string) *witness.Dmint {
	rule_mint_mod_history, err := m.get_mod_history(atomical_id)
	if err != nil {
		panic(err)
	}
	latest_state := calculate_latest_state_from_mod_history(rule_mint_mod_history)
	if validate_rules_data(latest_state) == nil {
		return nil
	}
	if latest_state.MintHeight < 0 {
		return nil
	}
	if latest_state.V != "1" {
		return nil
	}
	if len(latest_state.Merkle) != 64 {
		return nil
	}
	// todo: get_general_data_with_cache
	return latest_state
}

func validate_rules_data(namespace_data *witness.Dmint) []*witness.RuleInfo {
	if len(namespace_data.Rules) <= 0 || len(namespace_data.Rules) > utils.MAX_SUBNAME_RULE_ENTRIES {
		return nil
	}
	validated_rules_list := []*witness.RuleInfo{}
	for _, rule := range namespace_data.Rules {
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
	if validated_rules_list == nil || len(validated_rules_list) == 0 {
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

func validate_dmitem_mint_args_with_container_dmint(payload *witness.PayLoad, dmint *witness.Dmint) bool {
	for _, proof_item := range payload.Args.Proof {
		if len(proof_item.D) != 64 {
			return false
		}
	}
	image := payload.Args.DynamicFields[payload.Args.Main]
	is_proof_valid, err := validate_merkle_proof_dmint(dmint.Merkle, payload.Args.RequestDmitem, payload.Args.Bitworkc, payload.Args.Bitworkr, payload.Args.Main, utils.DoubleSha256(image), payload.Args.Proof)
	if err != nil {
		return false
	}
	return is_proof_valid
}

func validate_merkle_proof_dmint(merkleStr string, item_name string, possible_bitworkc, possible_bitworkr, main string, main_hash []byte, proof []witness.Proof) (bool, error) {
	expected_root_hash, err := hex.DecodeString(merkleStr)
	if err != nil {
		return false, err
	}
	// # Case 1: any/any
	concat_str1 := item_name + ":any" + ":any:" + main + ":" + hex.EncodeToString(main_hash)
	target_hash := utils.Sha256([]byte(concat_str1))
	// log.Log.Panicf("UpdateCurrentHeight err:%v", expected_root_hash1)

	if merkle.CheckValidateProof(expected_root_hash, target_hash, proof) {
		return true, nil
	}
	// # Case 2: specific_bitworkc/any
	if possible_bitworkc != "" {
		concat_str2 := item_name + ":" + possible_bitworkc + ":any:" + main + ":" + hex.EncodeToString(main_hash)
		target_hash := utils.Sha256([]byte(concat_str2))
		if merkle.CheckValidateProof(expected_root_hash, target_hash, proof) {
			return true, nil
		}
	}
	// # Case 3: any/specific_bitworkr
	if possible_bitworkr != "" {
		concat_str3 := item_name + ":any" + ":" + possible_bitworkr + ":" + main + ":" + hex.EncodeToString(main_hash)
		target_hash := utils.Sha256([]byte(concat_str3))
		if merkle.CheckValidateProof(expected_root_hash, target_hash, proof) {
			return true, nil
		}
	}
	// # Case 4:
	if possible_bitworkc != "" && possible_bitworkr != "" {
		concat_str4 := item_name + ":" + possible_bitworkc + ":" + possible_bitworkr + ":" + main + ":" + hex.EncodeToString(main_hash)
		target_hash := utils.Sha256([]byte(concat_str4))
		if merkle.CheckValidateProof(expected_root_hash, target_hash, proof) {
			return true, nil
		}
	}
	return false, nil
}

func calculate_latest_state_from_mod_history(mod_history []*witness.Dmint) *witness.Dmint {
	// Ensure it is sorted in ascending order
	// sort.Slice(mod_history, func(i, j int) bool {
	// 	return mod_history[i].ID < mod_history[j].ID
	// })
	current_object_state := &witness.Dmint{}
	for _, element := range mod_history {
		if element.A == 1 {
			current_object_state = nil
		} else {
			current_object_state = element
		}
	}
	return current_object_state
}

type regex_price_point struct {
	o        map[string]*witness.Output
	Bitworkc string
	Bitworkr string
	p        string
}

func (m *Atomicals) get_mod_history(parentContainerAtomicalsID string) ([]*witness.Dmint, error) {
	mod, err := m.Mod(parentContainerAtomicalsID)
	if err != nil {
		return nil, err
	}
	if mod == nil {
		return nil, nil
	}
	dmint := &witness.Dmint{}
	if err := json.Unmarshal([]byte(mod.Mod), dmint); err != nil {
		return nil, err
	}
	dmints := make([]*witness.Dmint, 0)
	dmints = append(dmints, dmint)
	return dmints, nil
}

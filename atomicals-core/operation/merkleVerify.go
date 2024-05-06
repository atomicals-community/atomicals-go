package atomicals

import (
	"encoding/hex"
	"regexp"
	"sort"
	"strings"

	"github.com/atomicals-go/atomicals-core/common"
	"github.com/atomicals-go/atomicals-core/witness"
	"github.com/atomicals-go/pkg/merkle"
)

// get_dmitem_parent_container_info
func verifyRuleAndMerkle(operation *witness.WitnessAtomicalsOperation, height int64) bool {
	matched_price_point := get_applicable_rule_by_height(operation.Payload.Args.ParentContainer, operation.Payload.Args.RequestDmitem, operation.CommitHeight-common.MINT_SUBNAME_RULES_BECOME_EFFECTIVE_IN_BLOCKS)
	dmint_validated_status := make_container_dmint_status_by_atomical_id_at_height(operation.Payload.Args.ParentContainer, height)
	if dmint_validated_status.status != "valid" {
		return false
	}
	mint_height := dmint_validated_status.mint_height
	expected_payment_height := operation.CommitHeight
	if expected_payment_height < mint_height {
		return false
	}
	if height < mint_height {
		return false
	}
	is_proof_valid := validate_dmitem_mint_args_with_container_dmint(operation.Payload.Args, operation.Payload, dmint_validated_status)
	if !is_proof_valid {
		return false
	}
	matched_rule := matched_price_point
	bitworkc := matched_rule.bitworkc
	bitworkr := matched_rule.bitworkr
	bitworkc_actual := operation.Payload.Args.Bitworkc
	bitworkr_actual := operation.Payload.Args.Bitworkr
	if bitworkc == "any" {
		return true
	} else if bitworkc != "any" {
		if bitworkc_actual != bitworkc {
			return false
		}
	}
	if bitworkr == "any" {
		return true
	} else if bitworkr != "any" {
		if bitworkr_actual != bitworkr {
			return false
		}
	}
	if matched_price_point.o != nil {
		return true
	}
	return false
}

func get_applicable_rule_by_height(parent_atomical_id string, proposed_subnameid string, height int64) *rule {
	rule_mint_mod_history := get_mod_history(parent_atomical_id, height)
	latest_state := calculate_latest_state_from_mod_history(rule_mint_mod_history)
	regex_price_point_list := validate_rules_data(latest_state)
	for _, regex_price_point := range regex_price_point_list {
		regex_pattern := regex_price_point.p
		if strings.Contains(regex_pattern, "(") || strings.Contains(regex_pattern, ")") {
			return nil
		}
		valid_pattern := regexp.MustCompile(regex_pattern)
		if !valid_pattern.MatchString(proposed_subnameid) {
			continue
		}
		return regex_price_point
	}
	return nil
}

func validate_rules_data(namespace_data *current_object_stateInfo_dmint) []*rule {
	rules := namespace_data.rules
	if len(rules) <= 0 || len(rules) > common.MAX_SUBNAME_RULE_ENTRIES {
		return nil
	}
	validated_rules_list := make([]*rule, 0)
	for _, rule_set_entry := range rules {
		regex_pattern := rule_set_entry.p
		if len(regex_pattern) > common.MAX_SUBNAME_RULE_SIZE_LEN || len(regex_pattern) < 1 {
			return nil
		}
		outputs := rule_set_entry.o
		bitworkc := rule_set_entry.bitworkc
		bitworkr := rule_set_entry.bitworkr
		if regex_pattern == "" {
			return nil
		}
		if strings.Contains(regex_pattern, "(") || strings.Contains(regex_pattern, ")") {
			return nil
		}
		price_point := &rule{}
		price_point.p = regex_pattern
		if outputs == nil && bitworkc == "" && bitworkr == "" {
			return nil
		}
		if bitworkc != "" {
			res := common.ParseBitwork(bitworkc)
			if res != nil {
				price_point.bitworkc = bitworkc
			} else if bitworkc == "any" {
				price_point.bitworkc = bitworkc
			} else {
				return nil
			}
		}
		if bitworkr != "" {
			res := common.ParseBitwork(bitworkr)
			if res != nil {
				price_point.bitworkr = bitworkr
			} else if bitworkr == "any" {
				price_point.bitworkr = bitworkr
			} else {
				return nil
			}
		}
		if outputs != nil {
			if !validate_subrealm_rules_outputs_format(outputs) {
				return nil
			}
			price_point.o = outputs
			validated_rules_list = append(validated_rules_list, price_point)
		} else if bitworkc != "" || bitworkr != "" {
			validated_rules_list = append(validated_rules_list, price_point)
		} else {
			return nil
		}
	}
	return validated_rules_list
}

func validate_subrealm_rules_outputs_format(outputs map[string]*output) bool {
	for expected_output_script, expected_output_value := range outputs {
		expected_output_id := expected_output_value.id
		expected_output_qty := expected_output_value.v
		if expected_output_qty < common.SUBNAME_MIN_PAYMENT_DUST_LIMIT {
			return false // # Reject if one of the entries expects less than the minimum payment amount
		}
		// # If there is a type restriction on the payment type then ensure it is a valid atomical id
		if expected_output_id != "" {
			if common.IsCompactAtomicalID(expected_output_id) {
				return false
			}
		}
		// # script must be paid to mint a subrealm
		if !common.IsHexString(expected_output_script) {
			return false // # Reject if one of the payment output script is not a valid hex
		}
	}
	return true
}

func make_container_dmint_status_by_atomical_id_at_height(atomical_id string, height int64) *current_object_stateInfo_dmint {
	rule_mint_mod_history := get_mod_history(atomical_id, height)
	latest_state := calculate_latest_state_from_mod_history(rule_mint_mod_history)

	dmint_format_status := get_container_dmint_format_status(latest_state)
	items := latest_state.items
	if items != nil {
		dmint_format_status.errors = append(dmint_format_status.errors, "items cannot be set manually for dmint")
		dmint_format_status.status = "invalid"
	}

	sealed_locationID := get_general_data_with_cache("sealed" + atomical_id)
	if sealed_locationID == "" {
		dmint_format_status.errors = append(dmint_format_status.errors, "container not sealed")
		dmint_format_status.status = "invalid"
	}
	dmint_format_status.dmint = items
	return dmint_format_status
}

func get_container_dmint_format_status(dmint *current_object_stateInfo_dmint) *current_object_stateInfo_dmint {
	base_status := &current_object_stateInfo_dmint{
		errors: make([]string, 0),
	}

	rules_list := validate_rules_data(dmint)
	if rules_list == nil || len(rules_list) == 0 {
		base_status.errors = append(base_status.errors, "rules list is invalid")
	}

	mint_height := dmint.mint_height
	if mint_height < 0 {
		base_status.errors = append(base_status.errors, "mint_height is invalid")
	}

	v := dmint.v
	if v != "1" {
		base_status.errors = append(base_status.errors, "v must be str 1")
	}

	immutable := dmint.immutable
	if immutable {
		base_status.errors = append(base_status.errors, "immutable must be a bool")
	}
	merkle := dmint.merkle
	if len(merkle) != 64 {
		base_status.errors = append(base_status.errors, "merkle str must be 64 hex characters")
	}

	if len(base_status.errors) == 0 {
		dmint.status = "valid"
	} else {
		dmint.status = "invalid"
	}
	return dmint
}

func validate_dmitem_mint_args_with_container_dmint(args *witness.Args, payload *witness.PayLoad, dmint *current_object_stateInfo_dmint) bool {
	for _, proof_item := range args.Proof {
		if len(proof_item.D) != 64 {
			return false
		}
	}
	request_dmitem := args.RequestDmitem
	merkle := dmint.merkle
	main := args.Main
	main_data := payload.Main
	main_hash := common.DoubleSha256(main_data)
	bitworkc := args.Bitworkc
	bitworkr := args.Bitworkr
	is_proof_valid, err := validate_merkle_proof_dmint(merkle, request_dmitem, bitworkc, bitworkr, main, main_hash, args.Proof)
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
	concat_str1Hex, err := hex.DecodeString(concat_str1)
	if err != nil {
		return false, err
	}
	target_hash := common.Sha256(concat_str1Hex)
	if merkle.CheckValidateProof(expected_root_hash, target_hash, proof) {
		return true, nil
	}
	// # Case 2: specific_bitworkc/any
	if possible_bitworkc != "" {
		concat_str2 := item_name + ":" + possible_bitworkc + ":any:" + main + ":" + hex.EncodeToString(main_hash)
		concat_str2Hex, err := hex.DecodeString(concat_str2)
		if err != nil {
			return false, err
		}
		target_hash := common.Sha256(concat_str2Hex)
		if merkle.CheckValidateProof(expected_root_hash, target_hash, proof) {
			return true, nil
		}
	}
	// # Case 3: any/specific_bitworkr
	if possible_bitworkr != "" {
		concat_str3 := item_name + ":any" + ":" + possible_bitworkr + ":" + main + ":" + hex.EncodeToString(main_hash)
		concat_str3Hex, err := hex.DecodeString(concat_str3)
		if err != nil {
			return false, err
		}
		target_hash := common.Sha256(concat_str3Hex)
		if merkle.CheckValidateProof(expected_root_hash, target_hash, proof) {
			return true, nil
		}
	}
	if possible_bitworkc != "" && possible_bitworkr != "" {
		concat_str4 := item_name + ":" + possible_bitworkc + ":" + possible_bitworkr + ":" + main + ":" + hex.EncodeToString(main_hash)
		concat_str4Hex, err := hex.DecodeString(concat_str4)
		if err != nil {
			return false, err
		}
		target_hash := common.Sha256(concat_str4Hex)

		if merkle.CheckValidateProof(expected_root_hash, target_hash, proof) {
			return true, nil
		}
	}
	return false, nil

}

func calculate_latest_state_from_mod_history(mod_history []*modInfo) *current_object_stateInfo_dmint {
	sort.Slice(mod_history, func(i, j int) bool {
		return mod_history[i].tx_num < mod_history[j].tx_num
	})
	current_object_state := &current_object_stateInfo_dmint{}
	// for _, element := range mod_history {
	// 	has_action_prop := element.data.a
	// 	if has_action_prop == 1 { //# delete = 1
	// 		for prop, _ := range element.data.props {
	// 			if _, exist := current_object_state[prop]; exist {
	// 				delete(current_object_state, prop)
	// 			}
	// 		}
	// 	} else {
	// 		current_object_state.dmint = element.data.props["element.data.props"]
	// 	}
	// }
	return current_object_state
}

type regex_price_point struct {
	o        string
	Bitworkc string
	Bitworkr string
	p        string
}

type modInfo struct {
	tx_num int64
	height int64
	txid   string
	index  int64
	data   *modData
}

type modData struct {
	a     int64 // $a
	args  string
	props map[string]string
}
type all_entryInfo struct {
	value  string // atomical_id
	tx_num int64
	cache  bool
}

type output struct {
	id string
	v  int64
}
type rule struct {
	p        string
	o        map[string]*output
	bitworkc string
	bitworkr string
}

type itemInfo struct {
	mint_height int64
}
type current_object_stateInfo_dmint struct {
	errors []string
	status string
	// dmint
	dmint *itemInfo
	items *itemInfo

	mint_height int64
	rules       []*rule
	merkle      string
	immutable   bool
	v           string
}

func get_mod_history(ParentContainerAtomicalsID string, height int64) []*modInfo {
	res := make([]*modInfo, 0)
	return res
}

func get_general_data_with_cache(atomicalsID string) string {
	return ""

}

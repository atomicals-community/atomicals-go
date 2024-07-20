package atomicals

import (
	"encoding/hex"

	"github.com/atomicals-go/atomicals-indexer/atomicals-core/witness"
	"github.com/atomicals-go/pkg/log"
	"github.com/atomicals-go/pkg/merkle"
	"github.com/atomicals-go/utils"
)

func (m *Atomicals) verifyRuleAndMerkle(operation *witness.WitnessAtomicalsOperation) bool {
	// get_dmitem_parent_container_info
	dmintValidatedStatus, err := m.getModHistory(operation.Payload.Args.ParentContainer, operation.RevealLocationHeight)
	if err != nil {
		panic(err)
	}
	if operation.CommitHeight < dmintValidatedStatus.MintHeight || operation.RevealLocationHeight < dmintValidatedStatus.MintHeight {
		return false
	}
	parentContainer, err := m.NftUTXOByAtomicalsID(operation.Payload.Args.ParentContainer)
	if err != nil {
		log.Log.Panicf("ParentContainerHasExist err:%v", err)
	}
	latestItem, err := m.LatestItemByContainerName(parentContainer.ContainerName)
	if err != nil {
		log.Log.Panicf("LatestItemByContainerName err:%v", err)
	}
	txID, _ := utils.SplitAtomicalsID(latestItem.LocationID)
	latsetMintHeight, err := m.AtomicalsTxHeight(txID)
	if err != nil {
		log.Log.Panicf("AtomicalsTxHeight err:%v", err)
	}
	if operation.CommitHeight < latsetMintHeight {
		return false
	}
	if operation.RevealLocationHeight < latsetMintHeight {
		return false
	}
	// validate_dmitem_mint_args_with_container_dmint
	for _, proof_item := range operation.Payload.Args.Proof {
		if len(proof_item.D) != 64 {
			return false
		}
	}
	image := operation.Payload.Args.DynamicFields[operation.Payload.Args.Main]
	is_proof_valid, err := validateMerkleProofDmint(dmintValidatedStatus.Merkle,
		operation.Payload.Args.RequestDmitem, operation.Payload.Args.Bitworkc,
		operation.Payload.Args.Bitworkr, operation.Payload.Args.Main, utils.DoubleSha256(image), operation.Payload.Args.Proof)
	if err != nil {
		return false
	}
	if !is_proof_valid {
		return false
	}
	matchedPricePoint := m.get_applicable_rule_by_height(operation.Payload.Args.ParentContainer, operation.Payload.Args.RequestDmitem, operation.RevealLocationHeight)
	return m.checkRule(matchedPricePoint, operation.Payload.Args.Bitworkc, operation.Payload.Args.Bitworkr)
}

func validateMerkleProofDmint(merkleStr string, item_name string, possible_bitworkc, possible_bitworkr, main string, main_hash []byte, proof []witness.Proof) (bool, error) {
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

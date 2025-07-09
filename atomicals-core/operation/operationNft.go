package atomicals

import (
	"github.com/atomicals-go/atomicals-core/witness"
	"github.com/atomicals-go/pkg/errors"
	"github.com/atomicals-go/pkg/log"
	"github.com/atomicals-go/repo/postsql"
	"github.com/atomicals-go/utils"
)

func (m *Atomicals) mintNft(operation *witness.WitnessAtomicalsOperation, userPk string) (newUTXONftInfo *postsql.UTXONftInfo, deleteUTXONfts []*postsql.UTXONftInfo, err error) {
	if operation.RevealInputIndex != 0 {
		return nil, nil, errors.ErrInvalidRevealInputIndex
	}
	if !operation.Payload.CheckRequest() {
		return nil, nil, errors.ErrCheckRequest
	}
	operation.CommitHeight = m.getTxHeight(operation.CommitTxID)
	if operation.CommitHeight < utils.ATOMICALS_ACTIVATION_HEIGHT {
		return nil, nil, errors.ErrInvalidCommitHeight
	}
	if !operation.IsWithinAcceptableBlocksForGeneralReveal() {
		return nil, nil, errors.ErrInvalidCommitHeight
	}
	if !operation.IsWithinAcceptableBlocksForNameReveal() {
		return nil, nil, errors.ErrInvalidCommitHeight
	}
	if operation.RevealLocationHeight >= utils.ATOMICALS_ACTIVATION_HEIGHT_COMMITZ && operation.CommitVoutIndex != utils.VOUT_EXPECT_OUTPUT_INDEX {
		return nil, nil, errors.ErrInvalidVinIndex
	}
	if operation.Payload.Args.RequestRealm != "" {
		if !utils.IsValidRealm(operation.Payload.Args.RequestRealm) {
			return nil, nil, errors.ErrInvalidRealm
		}
		b, err := m.GetBlockByHeight(operation.CommitHeight)
		if err != nil {
			log.Log.Panicf("GetBlockByHeight err:%v", err)
		}
		for txIndex, tx := range b.Tx {
			if tx.Txid == operation.CommitTxID {
				operation.CommitTxIndex = int64(txIndex)
			}
		}
		if operation.Payload.IsImmutable() {
			return nil, nil, errors.ErrCannotBeImmutable
		}
		if operation.Payload.Args.Bitworkc == "" {
			return nil, nil, errors.ErrBitworkcNeeded
		}
		bitworkc, _, err := operation.IsValidBitwork()
		if err != nil {
			return nil, nil, err
		}
		if bitworkc != nil && len(bitworkc.Prefix) < 4 {
			return nil, nil, errors.ErrInvalidBitworkcPrefix
		}
		realms, err := m.NftRealmByName(operation.Payload.Args.RequestRealm)
		if err != nil {
			log.Log.Panicf("NftRealmByName err:%v", err)
		}
		if len(realms) != 0 {
			for _, v := range realms {
				if v.CommitHeight < operation.CommitHeight {
					return nil, nil, errors.ErrRealmHasExist
				} else if v.CommitHeight == operation.CommitHeight {
					if v.CommitTxIndex < operation.CommitTxIndex {
						return nil, nil, errors.ErrRealmHasExist
					} else {
						deleteUTXONfts = append(deleteUTXONfts, realms...)
					}
				} else {
					deleteUTXONfts = append(deleteUTXONfts, realms...)
				}
			}
		}
		newUTXONftInfo = &postsql.UTXONftInfo{
			UserPk:        userPk,
			RealmName:     operation.Payload.Args.RequestRealm,
			Time:          operation.Payload.Args.Time,
			Bitworkc:      operation.Payload.Args.Bitworkc,
			Bitworkr:      operation.Payload.Args.Bitworkr,
			AtomicalsID:   operation.AtomicalsID,
			CommitHeight:  operation.CommitHeight,
			CommitTxIndex: operation.CommitTxIndex,
			LocationID:    operation.LocationID,
		}

	} else if operation.Payload.Args.RequestSubRealm != "" {
		if !utils.IsValidSubRealm(operation.Payload.Args.RequestSubRealm) {
			return nil, nil, errors.ErrInvalidContainer
		}
		parentRealm, err := m.NftUTXOByAtomicalsID(operation.Payload.Args.ParentRealm)
		if err != nil {
			log.Log.Panicf("NftUTXOByAtomicalsID err:%v", err)
		}
		if parentRealm == nil {
			return nil, nil, errors.ErrParentRealmNotExist
		}
		isExist, err := m.NftSubRealmByNameHasExist(operation.Payload.Args.ParentRealm, operation.Payload.Args.RequestSubRealm)
		if err != nil {
			log.Log.Panicf("NftSubRealmByName err:%v", err)
		}
		if isExist {
			return nil, nil, errors.ErrSubRealmHasExist
		}
		if operation.Payload.IsImmutable() {
			return nil, nil, errors.ErrCannotBeImmutable
		}
		if operation.Payload.Args.ClaimType == witness.Rule {
			matched_rule := m.get_applicable_rule_by_height(operation.Payload.Args.ParentRealm,
				operation.Payload.Args.RequestSubRealm, operation.CommitHeight-utils.MINT_SUBNAME_RULES_BECOME_EFFECTIVE_IN_BLOCKS)
			if !m.checkRule(matched_rule, operation.Payload.Args.Bitworkc, operation.Payload.Args.Bitworkr) {
				return nil, nil, errors.ErrInvalidRule
			}
		}
		newUTXONftInfo = &postsql.UTXONftInfo{
			UserPk:                 userPk,
			RealmName:              parentRealm.RealmName,
			SubRealmName:           operation.Payload.Args.RequestSubRealm,
			ClaimType:              operation.Payload.Args.ClaimType,
			ParentRealmAtomicalsID: operation.Payload.Args.ParentRealm,
			Time:                   operation.Payload.Args.Time,
			Bitworkc:               operation.Payload.Args.Bitworkc,
			Bitworkr:               operation.Payload.Args.Bitworkr,
			AtomicalsID:            operation.AtomicalsID,
			LocationID:             operation.LocationID,
		}
	} else if operation.Payload.Args.RequestContainer != "" {
		if !utils.IsValidContainer(operation.Payload.Args.RequestContainer) {
			return nil, nil, errors.ErrInvalidContainer
		}
		isExist, err := m.NftContainerByNameHasExist(operation.Payload.Args.RequestContainer)
		if err != nil {
			log.Log.Panicf("NftContainerByName err:%v", err)
		}
		if isExist {
			return nil, nil, errors.ErrContainerHasExist
		}
		if operation.Payload.IsImmutable() {
			return nil, nil, errors.ErrCannotBeImmutable
		}
		if operation.Payload.Args.Bitworkc == "" {
			return nil, nil, errors.ErrBitworkcNeeded
		}
		newUTXONftInfo = &postsql.UTXONftInfo{
			UserPk:        userPk,
			ContainerName: operation.Payload.Args.RequestContainer,
			Time:          operation.Payload.Args.Time,
			Bitworkc:      operation.Payload.Args.Bitworkc,
			Bitworkr:      operation.Payload.Args.Bitworkr,
			AtomicalsID:   operation.AtomicalsID,
			LocationID:    operation.LocationID,
		}
		bitworkc, _, err := operation.IsValidBitwork()
		if err != nil {
			return nil, nil, err
		}
		if bitworkc != nil && len(bitworkc.Prefix) < 4 {
			return nil, nil, errors.ErrInvalidBitworkcPrefix
		}
	} else if operation.Payload.Args.RequestDmitem != "" {
		if !utils.IsDmintActivated(operation.RevealLocationHeight) {
			return nil, nil, errors.ErrDmintNotStart
		}
		parentContainer, err := m.NftUTXOByAtomicalsID(operation.Payload.Args.ParentContainer)
		if err != nil {
			log.Log.Panicf("NftUTXOByAtomicalsID err:%v", err)
		}
		if parentContainer == nil {
			return nil, nil, errors.ErrContainerNotExist
		}

		if !utils.IsValidDmitem(operation.Payload.Args.RequestDmitem) {
			return nil, nil, errors.ErrInvalidContainerDmitem
		}
		isExist, err := m.ContainerItemByNameHasExist(parentContainer.ContainerName, operation.Payload.Args.RequestDmitem)
		if err != nil {
			log.Log.Panicf("ContainerItemByName err:%v", err)
		}
		if isExist {
			return nil, nil, errors.ErrSubRealmHasExist
		}
		if !m.verifyRuleAndMerkle(operation) {
			return nil, nil, errors.ErrInvalidMerkleVerify
		}
		newUTXONftInfo = &postsql.UTXONftInfo{
			UserPk:                     userPk,
			ContainerName:              parentContainer.ContainerName,
			Dmitem:                     operation.Payload.Args.RequestDmitem,
			ParentContainerAtomicalsID: operation.Payload.Args.ParentContainer,
			Time:                       operation.Payload.Args.Time,
			Bitworkc:                   operation.Payload.Args.Bitworkc,
			Bitworkr:                   operation.Payload.Args.Bitworkr,
			AtomicalsID:                operation.AtomicalsID,
			LocationID:                 operation.LocationID,
		}
	}
	return newUTXONftInfo, deleteUTXONfts, nil
}

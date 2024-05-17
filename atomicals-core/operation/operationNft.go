package atomicals

import (
	"github.com/atomicals-go/atomicals-core/witness"
	"github.com/atomicals-go/pkg/errors"
	"github.com/atomicals-go/pkg/log"
	"github.com/atomicals-go/repo/postsql"
	"github.com/atomicals-go/utils"
)

func (m *Atomicals) mintNft(operation *witness.WitnessAtomicalsOperation, userPk string, height int64) error {
	// if operation.Payload.Args.RequestContainer == "" && operation.Payload.Args.RequestDmitem == "" {
	// 	return nil
	// }
	if operation.RevealInputIndex != 0 {
		return errors.ErrInvalidRevealInputIndex
	}
	if !operation.Payload.CheckRequest() {
		return errors.ErrCheckRequest
	}
	bitworkc, _, err := operation.IsValidBitwork()
	if err != nil {
		return err
	}
	operation.CommitHeight, err = m.BtcTxHeight(operation.CommitTxID)
	if err != nil {
		panic(err)
	}
	if operation.CommitHeight < utils.ATOMICALS_ACTIVATION_HEIGHT {
		return errors.ErrInvalidCommitHeight
	}
	if !operation.IsWithinAcceptableBlocksForGeneralReveal() {
		return errors.ErrInvalidCommitHeight
	}
	if !operation.IsWithinAcceptableBlocksForNameReveal() {
		return errors.ErrInvalidCommitHeight
	}
	if operation.RevealLocationHeight >= utils.ATOMICALS_ACTIVATION_HEIGHT_COMMITZ && operation.CommitVoutIndex != utils.VOUT_EXPECT_OUTPUT_INDEX {
		return errors.ErrInvalidVinIndex
	}
	if operation.Payload.Args.RequestRealm != "" {
		if !utils.IsValidRealm(operation.Payload.Args.RequestRealm) {
			return errors.ErrInvalidRealm
		}
		isExist, err := m.NftRealmByNameHasExist(operation.Payload.Args.RequestRealm)
		if err != nil {
			log.Log.Panicf("NftRealmByNameHasExist err:%v", err)
		}
		if isExist {
			return errors.ErrRealmHasExist
		}
		if operation.Payload.IsImmutable() {
			return errors.ErrCannotBeImmutable
		}
		if operation.Payload.Args.Bitworkc == "" {
			return errors.ErrBitworkcNeeded
		}
		entity := &postsql.UTXONftInfo{
			UserPk:      userPk,
			RealmName:   operation.Payload.Args.RequestRealm,
			Nonce:       operation.Payload.Args.Nonce,
			Time:        operation.Payload.Args.Time,
			Bitworkc:    operation.Payload.Args.Bitworkc,
			Bitworkr:    operation.Payload.Args.Bitworkr,
			AtomicalsID: operation.AtomicalsID,
			LocationID:  operation.LocationID,
		}
		if bitworkc != nil && len(bitworkc.Prefix) < 4 {
			return errors.ErrInvalidBitworkcPrefix
		}
		if err := m.InsertNftUTXO(entity); err != nil {
			log.Log.Panicf("InsertNftUTXO err:%v", err)
		}
	} else if operation.Payload.Args.RequestSubRealm != "" {
		if !utils.IsValidSubRealm(operation.Payload.Args.RequestSubRealm) {
			return errors.ErrInvalidContainer
		}
		if operation.Payload.Args.ClaimType != witness.Direct && operation.Payload.Args.ClaimType != witness.Rule {
			return errors.ErrInvalidClaimType
		}
		parentRealmName, err := m.ParentRealmHasExist(operation.Payload.Args.ParentRealm)

		if err != nil {
			log.Log.Panicf("ParentRealmHasExist err:%v", err)
		}
		if parentRealmName == "" {
			return errors.ErrParentRealmNotExist
		}
		isExist, err := m.NftSubRealmByNameHasExist(parentRealmName, operation.Payload.Args.RequestSubRealm)
		if err != nil {
			log.Log.Panicf("NftSubRealmByName err:%v", err)
		}
		if isExist {
			return errors.ErrSubRealmHasExist
		}
		if operation.Payload.IsImmutable() {
			return errors.ErrCannotBeImmutable
		}
		entity := &postsql.UTXONftInfo{
			UserPk:                 userPk,
			RealmName:              parentRealmName,
			SubRealmName:           operation.Payload.Args.RequestSubRealm,
			ClaimType:              operation.Payload.Args.ClaimType,
			ParentRealmAtomicalsID: operation.Payload.Args.ParentRealm,
			Nonce:                  operation.Payload.Args.Nonce,
			Time:                   operation.Payload.Args.Time,
			Bitworkc:               operation.Payload.Args.Bitworkc,
			Bitworkr:               operation.Payload.Args.Bitworkr,
			AtomicalsID:            operation.AtomicalsID,
			LocationID:             operation.LocationID,
		}
		if err := m.InsertNftUTXO(entity); err != nil {
			log.Log.Panicf("InsertNftUTXO err:%v", err)
		}
	} else if operation.Payload.Args.RequestContainer != "" {
		if !utils.IsValidContainer(operation.Payload.Args.RequestContainer) {
			return errors.ErrInvalidContainer
		}
		isExist, err := m.NftContainerByNameHasExist(operation.Payload.Args.RequestContainer)
		if err != nil {
			log.Log.Panicf("NftContainerByName err:%v", err)
		}
		if isExist {
			return errors.ErrContainerHasExist
		}
		if operation.Payload.IsImmutable() {
			return errors.ErrCannotBeImmutable
		}
		if operation.Payload.Args.Bitworkc == "" {
			return errors.ErrBitworkcNeeded
		}
		entity := &postsql.UTXONftInfo{
			UserPk:        userPk,
			ContainerName: operation.Payload.Args.RequestContainer,
			Nonce:         operation.Payload.Args.Nonce,
			Time:          operation.Payload.Args.Time,
			Bitworkc:      operation.Payload.Args.Bitworkc,
			Bitworkr:      operation.Payload.Args.Bitworkr,
			AtomicalsID:   operation.AtomicalsID,
			LocationID:    operation.LocationID,
		}
		if bitworkc != nil && len(bitworkc.Prefix) < 4 {
			return errors.ErrInvalidBitworkcPrefix
		}
		if err := m.InsertNftUTXO(entity); err != nil {
			log.Log.Panicf("InsertNftUTXO err:%v", err)
		}
	} else if operation.Payload.Args.RequestDmitem != "" {
		if !utils.IsDmintActivated(operation.RevealLocationHeight) {
			return errors.ErrDmintNotStart
		}
		parentContainer, err := m.ParentContainerHasExist(operation.Payload.Args.ParentContainer)
		if err != nil {
			log.Log.Panicf("ParentContainerHasExist err:%v", err)
		}
		if parentContainer == nil {
			return errors.ErrContainerNotExist
		}
		if !utils.IsValidDmitem(operation.Payload.Args.RequestDmitem) {
			return errors.ErrInvalidContainerDmitem
		}
		isExist, err := m.ContainerItemByNameHasExist(parentContainer.ContainerName, operation.Payload.Args.RequestDmitem)
		if err != nil {
			log.Log.Panicf("ContainerItemByName err:%v", err)
		}
		if isExist {
			return errors.ErrSubRealmHasExist
		}
		if height >= utils.ATOMICALS_ACTIVATION_HEIGHT_DMINT {
			if !m.verifyRuleAndMerkle(operation, height) {
				return errors.ErrInvalidMerkleVerify
			}
		}
		entity := &postsql.UTXONftInfo{
			UserPk:                     userPk,
			ContainerName:              parentContainer.ContainerName,
			Dmitem:                     operation.Payload.Args.RequestDmitem,
			ParentContainerAtomicalsID: operation.Payload.Args.ParentContainer,
			Nonce:                      operation.Payload.Args.Nonce,
			Time:                       operation.Payload.Args.Time,
			Bitworkc:                   operation.Payload.Args.Bitworkc,
			Bitworkr:                   operation.Payload.Args.Bitworkr,
			AtomicalsID:                operation.AtomicalsID,
			LocationID:                 operation.LocationID,
		}
		if err := m.InsertNftUTXO(entity); err != nil {
			log.Log.Panicf("InsertNftUTXO err:%v", err)
		}
	}
	return nil
}

package atomicals

import (
	"github.com/atomicals-core/pkg/log"

	"github.com/atomicals-core/atomicals/DB/postsql"
	"github.com/atomicals-core/atomicals/common"
	"github.com/atomicals-core/atomicals/witness"
	"github.com/atomicals-core/pkg/errors"
)

func (m *Atomicals) mintNft(operation *witness.WitnessAtomicalsOperation, userPk string) error {
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
	operation.CommitHeight = m.GetCommitHeight(operation.CommitTxID)
	if operation.CommitHeight < common.ATOMICALS_ACTIVATION_HEIGHT {
		return errors.ErrInvalidCommitHeight
	}
	if !operation.IsWithinAcceptableBlocksForGeneralReveal() {
		return errors.ErrInvalidCommitHeight
	}
	if !operation.IsWithinAcceptableBlocksForNameReveal() {
		return errors.ErrInvalidCommitHeight
	}
	if operation.RevealLocationHeight >= common.ATOMICALS_ACTIVATION_HEIGHT_COMMITZ && operation.CommitVoutIndex != common.VOUT_EXPECT_OUTPUT_INDEX {
		return errors.ErrInvalidVinIndex
	}
	atomicalsID := operation.AtomicalsID
	if operation.Payload.Args.RequestRealm != "" {
		if !common.IsValidRealm(operation.Payload.Args.RequestRealm) {
			return errors.ErrInvalidRealm
		}
		isExist, err := m.NftRealmByNameHasExist(operation.Payload.Args.RequestRealm)
		if err != nil {
			log.Log.Panicf("NftRealmByNameHasExist err:%v", err)
		}
		if isExist {
			return errors.ErrRealmHasExist
		}
		if operation.IsImmutable() {
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
			AtomicalsID: atomicalsID,
			LocationID:  atomicalsID,
		}
		if bitworkc != nil && len(bitworkc.Prefix) < 4 {
			return errors.ErrInvalidBitworkcPrefix
		}
		if err := m.InsertNftUTXO(entity); err != nil {
			log.Log.Panicf("InsertNftUTXO err:%v", err)
		}
	} else if operation.Payload.Args.RequestSubRealm != "" {
		if !common.IsValidSubRealm(operation.Payload.Args.RequestSubRealm) {
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
		if operation.IsImmutable() {
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
			AtomicalsID:            atomicalsID,
			LocationID:             atomicalsID,
		}
		if err := m.InsertNftUTXO(entity); err != nil {
			log.Log.Panicf("InsertNftUTXO err:%v", err)
		}
	} else if operation.Payload.Args.RequestContainer != "" {
		if !common.IsValidContainer(operation.Payload.Args.RequestContainer) {
			return errors.ErrInvalidContainer
		}
		isExist, err := m.NftContainerByNameHasExist(operation.Payload.Args.RequestContainer)
		if err != nil {
			log.Log.Panicf("NftContainerByName err:%v", err)
		}
		if isExist {
			return errors.ErrContainerHasExist
		}
		if operation.IsImmutable() {
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
			AtomicalsID:   atomicalsID,
			LocationID:    atomicalsID,
		}
		if bitworkc != nil && len(bitworkc.Prefix) < 4 {
			return errors.ErrInvalidBitworkcPrefix
		}
		if err := m.InsertNftUTXO(entity); err != nil {
			log.Log.Panicf("InsertNftUTXO err:%v", err)
		}
	} else if operation.Payload.Args.RequestDmitem != "" {
		if !common.IsDmintActivated(operation.RevealLocationHeight) {
			return errors.ErrDmintNotStart
		}
		parentContainerName, err := m.ParentContainerHasExist(operation.Payload.Args.ParentContainer)
		if err != nil {
			log.Log.Panicf("ParentContainerHasExist err:%v", err)
		}
		if parentContainerName == "" {
			return errors.ErrParentRealmNotExist
		}
		if !common.IsValidDmitem(operation.Payload.Args.RequestDmitem) {
			return errors.ErrInvalidContainerDmitem
		}
		isExist, err := m.ContainerItemByNameHasExist(parentContainerName, operation.Payload.Args.RequestDmitem)
		if err != nil {
			log.Log.Panicf("ContainerItemByName err:%v", err)
		}
		if isExist {
			return errors.ErrSubRealmHasExist
		}
		entity := &postsql.UTXONftInfo{
			UserPk:                     userPk,
			ContainerName:              parentContainerName,
			Dmitem:                     operation.Payload.Args.RequestDmitem,
			ParentContainerAtomicalsID: operation.Payload.Args.ParentContainer,
			Nonce:                      operation.Payload.Args.Nonce,
			Time:                       operation.Payload.Args.Time,
			Bitworkc:                   operation.Payload.Args.Bitworkc,
			Bitworkr:                   operation.Payload.Args.Bitworkr,
			AtomicalsID:                atomicalsID,
			LocationID:                 atomicalsID,
		}
		if err := m.InsertNftUTXO(entity); err != nil {
			log.Log.Panicf("InsertNftUTXO err:%v", err)
		}
	}
	// skip other nft
	// else {
	// 	log.Log.Panicf("operation.Script:%+v", operation.Script)
	// 	log.Log.Panicf("operation.Payload.Args:%+v", operation.Payload.Args)
	// }
	return nil
}

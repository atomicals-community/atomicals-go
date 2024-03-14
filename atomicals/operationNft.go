package atomicals

import (
	"github.com/atomicals-core/pkg/log"

	db "github.com/atomicals-core/atomicals/DB"
	"github.com/atomicals-core/atomicals/common"
	"github.com/atomicals-core/atomicals/witness"
	"github.com/atomicals-core/pkg/errors"
	"github.com/btcsuite/btcd/btcjson"
)

func (m *Atomicals) mintNft(operation *witness.WitnessAtomicalsOperation, vin btcjson.Vin, vout []btcjson.Vout, userPk string) error {
	if !operation.IsWithinAcceptableBlocksForGeneralReveal() {
		return errors.ErrInvalidCommitHeight
	}
	if operation.CommitHeight < common.ATOMICALS_ACTIVATION_HEIGHT {
		return errors.ErrInvalidCommitHeight
	}
	if operation.CommitVoutIndex != common.VOUT_EXPECT_OUTPUT_INDEX {
		return errors.ErrInvalidVinIndex
	}
	bitworkc, bitworkr, err := operation.IsValidBitwork()
	if err != nil {
		return err
	}
	atomicalsID := operation.AtomicalsID
	if operation.Payload.Args.RequestRealm != "" {
		// seems not necessary
		if !operation.IsValidCommitVoutIndexForNameRevel() {
			return errors.ErrInvalidCommitVoutIndex
		}
		if !operation.IsWithinAcceptableBlocksForNameReveal() {
			return errors.ErrInvalidCommitHeight
		}
		if operation.IsImmutable() {
			return errors.ErrCannotBeImmutable
		}
		entity := &db.UserNftInfo{
			RealmName:   operation.Payload.Args.RequestRealm,
			Nonce:       operation.Payload.Args.Nonce,
			Time:        operation.Payload.Args.Time,
			Bitworkc:    bitworkc,
			Bitworkr:    bitworkr,
			AtomicalsID: atomicalsID,
			LocationID:  atomicalsID,
		}
		if entity.Bitworkc != nil && len(entity.Bitworkc.Prefix) < 4 {
			return errors.ErrInvalidBitworkcPrefix
		}
		if !common.IsValidRealm(entity.RealmName) {
			return errors.ErrInvalidRealm
		}
		if subRealmMap, err := m.NftRealmByName(entity.RealmName); err != nil && subRealmMap != nil {
			return errors.ErrRealmHasExist
		}
		subRealmMap, err := m.NftRealmByName(entity.RealmName)
		if err != nil {
			log.Log.Panicf("NftRealmByName err:%v", err)
		}
		if subRealmMap != nil {
			return errors.ErrSubRealmHasExist
		}
		if err := m.InsertNftUTXOByAtomicalsID(entity); err != nil {
			log.Log.Panicf("InsertNftUTXOByAtomicalsID err:%v", err)
		}
		if err := m.InsertNftUTXOByLocationID(entity); err != nil {
			log.Log.Panicf("InsertNftUTXOByLocationID err:%v", err)
		}
		if err := m.InsertRealm(entity.RealmName); err != nil {
			log.Log.Panicf("InsertRealm err:%v", err)
		}
	} else if operation.Payload.Args.RequestSubRealm != "" {
		// seems not necessary
		if !operation.IsValidCommitVoutIndexForNameRevel() {
			return errors.ErrInvalidCommitVoutIndex
		}
		if !operation.IsWithinAcceptableBlocksForNameReveal() {
			return errors.ErrInvalidCommitHeight
		}
		if operation.IsImmutable() {
			return errors.ErrCannotBeImmutable
		}
		entity := &db.UserNftInfo{
			SubRealmName:           operation.Payload.Args.RequestSubRealm,
			ClaimType:              operation.Payload.Args.ClaimType,
			ParentRealmAtomicalsID: operation.Payload.Args.ParentRealm,
			Nonce:                  operation.Payload.Args.Nonce,
			Time:                   operation.Payload.Args.Time,
			Bitworkc:               bitworkc,
			Bitworkr:               bitworkr,
			AtomicalsID:            atomicalsID,
			LocationID:             atomicalsID,
		}
		if entity.ClaimType != witness.Direct && entity.ClaimType != witness.Rule {
			return errors.ErrInvalidClaimType
		}
		if !common.IsValidSubRealm(entity.SubRealmName) {
			return errors.ErrInvalidContainer
		}
		parentRealmName, err := m.ParentRealmHasExist(entity.ParentRealmAtomicalsID)
		if err != nil {
			log.Log.Panicf("ParentRealmHasExist err:%v", err)
		}
		if parentRealmName == "" {
			return errors.ErrParentRealmNotExist
		}
		isExist, err := m.NftSubRealmByName(parentRealmName, entity.SubRealmName)
		if err != nil {
			log.Log.Panicf("NftSubRealmByName err:%v", err)
		}
		if isExist {
			return errors.ErrSubRealmHasExist
		}
		if err := m.InsertNftUTXOByAtomicalsID(entity); err != nil {
			log.Log.Panicf("InsertNftUTXOByAtomicalsID err:%v", err)
		}
		if err := m.InsertNftUTXOByLocationID(entity); err != nil {
			log.Log.Panicf("InsertNftUTXOByLocationID err:%v", err)
		}
		if err := m.InsertSubRealm(entity.RealmName, entity.SubRealmName); err != nil {
			log.Log.Panicf("InsertSubRealm err:%v", err)
		}
	} else if operation.Payload.Args.RequestContainer != "" {
		// seems not necessary
		if !operation.IsValidCommitVoutIndexForNameRevel() {
			return errors.ErrInvalidCommitVoutIndex
		}
		if !operation.IsWithinAcceptableBlocksForNameReveal() {
			return errors.ErrInvalidCommitHeight
		}
		if operation.IsImmutable() {
			return errors.ErrCannotBeImmutable
		}
		entity := &db.UserNftInfo{
			ContainerName: operation.Payload.Args.RequestContainer,
			Nonce:         operation.Payload.Args.Nonce,
			Time:          operation.Payload.Args.Time,
			Bitworkc:      bitworkc,
			Bitworkr:      bitworkr,
			AtomicalsID:   atomicalsID,
			LocationID:    atomicalsID,
		}
		if entity.Bitworkc != nil && len(entity.Bitworkc.Prefix) < 4 {
			return errors.ErrInvalidBitworkcPrefix
		}
		if !common.IsValidContainer(entity.ContainerName) {
			return errors.ErrInvalidContainer
		}
		items, err := m.NftContainerByName(entity.ContainerName)
		if err != nil {
			log.Log.Panicf("NftContainerByName err:%v", err)
		}
		if items != nil {
			return errors.ErrContainerHasExist
		}
		if err := m.InsertNftUTXOByAtomicalsID(entity); err != nil {
			log.Log.Panicf("InsertNftUTXOByAtomicalsID err:%v", err)
		}
		if err := m.InsertNftUTXOByLocationID(entity); err != nil {
			log.Log.Panicf("InsertNftUTXOByLocationID err:%v", err)
		}
		if err := m.InsertContainer(entity.ContainerName); err != nil {
			log.Log.Panicf("InsertContainer err:%v", err)
		}
	} else if operation.Payload.Args.RequestDmitem != "" {
		// seems not necessary
		if !operation.IsValidCommitVoutIndexForNameRevel() {
			return errors.ErrInvalidCommitVoutIndex
		}
		if !operation.IsWithinAcceptableBlocksForNameReveal() {
			return errors.ErrInvalidCommitHeight
		}
		entity := &db.UserNftInfo{
			Dmitem:                     operation.Payload.Args.RequestDmitem,
			ParentContainerAtomicalsID: operation.Payload.Args.ParentContainer,
			Nonce:                      operation.Payload.Args.Nonce,
			Time:                       operation.Payload.Args.Time,
			// Bitworkc:                   bitworkc,
			// Bitworkr:                   bitworkr,
			AtomicalsID: atomicalsID,
			LocationID:  atomicalsID,
		}
		if !common.IsValidDmitem(entity.Dmitem) {
			return errors.ErrInvalidContainerDmitem
		}
		parentContainerName, err := m.ParentContainerHasExist(entity.ParentContainerAtomicalsID)
		if err != nil {
			log.Log.Panicf("ParentContainerHasExist err:%v", err)
		}
		if parentContainerName == "" {
			return errors.ErrParentRealmNotExist
		}
		if err := m.InsertNftUTXOByAtomicalsID(entity); err != nil {
			log.Log.Panicf("InsertNftUTXOByAtomicalsID err:%v", err)
		}
		if err := m.InsertNftUTXOByLocationID(entity); err != nil {
			log.Log.Panicf("InsertNftUTXOByLocationID err:%v", err)
		}
		if err := m.InsertItemInContainer(parentContainerName, entity.Dmitem); err != nil {
			log.Log.Panicf("InsertItemInContainer err:%v", err)
		}
	} else {
		log.Log.Warnf("operation.Script:%+v", operation.Script)
		log.Log.Warnf("operation.Payload.Args:%+v", operation.Payload.Args)
	}
	return nil
}

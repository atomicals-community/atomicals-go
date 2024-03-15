package atomicals

import (
	"github.com/atomicals-core/pkg/log"

	db "github.com/atomicals-core/atomicals/DB"
	"github.com/atomicals-core/atomicals/common"
	"github.com/atomicals-core/atomicals/witness"
	"github.com/atomicals-core/pkg/errors"
	"github.com/btcsuite/btcd/btcjson"
)

func (m *Atomicals) mintNft(operation *witness.WitnessAtomicalsOperation, vout []btcjson.Vout, userPk string) error {
	if !operation.Payload.CheckRequest() {
		return errors.ErrCheckRequest
	}
	bitworkc, bitworkr, err := operation.IsValidBitwork()
	if err != nil {
		return err
	}
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
		subRealmMap, err := m.NftRealmByName(operation.Payload.Args.RequestRealm)
		if err != nil {
			log.Log.Panicf("NftRealmByName err:%v", err)
		}
		if subRealmMap != nil {
			return errors.ErrRealmHasExist
		}
		if operation.IsImmutable() {
			return errors.ErrCannotBeImmutable
		}
		if operation.Payload.Args.Bitworkc == "" {
			return errors.ErrBitworkcNeeded
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
		isExist, err := m.NftSubRealmByName(parentRealmName, operation.Payload.Args.RequestSubRealm)
		if err != nil {
			log.Log.Panicf("NftSubRealmByName err:%v", err)
		}
		if isExist {
			return errors.ErrSubRealmHasExist
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
		if !common.IsValidContainer(operation.Payload.Args.RequestContainer) {
			return errors.ErrInvalidContainer
		}
		container, err := m.NftContainerByName(operation.Payload.Args.RequestContainer)
		if err != nil {
			log.Log.Panicf("NftContainerByName err:%v", err)
		}
		if container != nil {
			return errors.ErrContainerHasExist
		}
		if operation.IsImmutable() {
			return errors.ErrCannotBeImmutable
		}
		if operation.Payload.Args.Bitworkc == "" {
			return errors.ErrBitworkcNeeded
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
		entity := &db.UserNftInfo{
			Dmitem:                     operation.Payload.Args.RequestDmitem,
			ParentContainerAtomicalsID: operation.Payload.Args.ParentContainer,
			Nonce:                      operation.Payload.Args.Nonce,
			Time:                       operation.Payload.Args.Time,
			Bitworkc:                   bitworkc,
			Bitworkr:                   bitworkr,
			AtomicalsID:                atomicalsID,
			LocationID:                 atomicalsID,
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

package atomicals

import (
	"github.com/atomicals-core/pkg/log"

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
	atomicalsID := atomicalsID(operation.RevealLocationTxID, operation.RevealLocationVoutIndex)
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
		entity := &UserNftInfo{
			RealmName:   operation.Payload.Args.RequestRealm,
			Nonce:       operation.Payload.Args.Nonce,
			Time:        operation.Payload.Args.Time,
			Bitworkc:    bitworkc,
			Bitworkr:    bitworkr,
			AtomicalsID: atomicalsID,
			Location:    atomicalsID,
		}
		if entity.Bitworkc != nil && len(entity.Bitworkc.Prefix) < 4 {
			return errors.ErrInvalidBitworkcPrefix
		}
		if entity.Bitworkc == nil {
			return errors.ErrNameTypeMintMastHaveBitworkc
		}
		if !common.IsValidRealm(entity.RealmName) {
			return errors.ErrInvalidRealm
		}
		if m.RealmHasExist(entity.RealmName) {
			return errors.ErrRealmHasExist
		}
		m.ensureUTXONotNil(atomicalsID)
		m.UTXOs[atomicalsID].Nft = append(m.UTXOs[atomicalsID].Nft, entity)
		m.GlobalNftRealmMap[entity.RealmName] = make(map[string]bool)
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
		entity := &UserNftInfo{
			SubRealmName:           operation.Payload.Args.RequestSubRealm,
			ClaimType:              operation.Payload.Args.ClaimType,
			ParentRealmAtomicalsID: operation.Payload.Args.ParentRealm,
			Nonce:                  operation.Payload.Args.Nonce,
			Time:                   operation.Payload.Args.Time,
			Bitworkc:               bitworkc,
			Bitworkr:               bitworkr,
			AtomicalsID:            atomicalsID,
			Location:               atomicalsID,
		}
		if entity.ClaimType != witness.Direct && entity.ClaimType != witness.Rule {
			return errors.ErrInvalidClaimType
		}
		if !common.IsValidSubRealm(entity.SubRealmName) {
			return errors.ErrInvalidContainer
		}
		parentRealmName, ok := m.ParentRealmHasExist(entity.ParentRealmAtomicalsID)
		if !ok {
			return errors.ErrParentRealmNotExist
		}
		if m.SubRealmHasExist(parentRealmName, entity.SubRealmName) {
			return errors.ErrSubRealmHasExist
		}
		m.ensureUTXONotNil(atomicalsID)
		m.UTXOs[atomicalsID].Nft = append(m.UTXOs[atomicalsID].Nft, entity)
		m.GlobalNftRealmMap[parentRealmName][entity.SubRealmName] = true
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
		entity := &UserNftInfo{
			ContainerName: operation.Payload.Args.RequestContainer,
			Nonce:         operation.Payload.Args.Nonce,
			Time:          operation.Payload.Args.Time,
			Bitworkc:      bitworkc,
			Bitworkr:      bitworkr,
			AtomicalsID:   atomicalsID,
			Location:      atomicalsID,
		}
		if entity.Bitworkc != nil && len(entity.Bitworkc.Prefix) < 4 {
			return errors.ErrInvalidBitworkcPrefix
		}
		if entity.Bitworkc == nil {
			return errors.ErrNameTypeMintMastHaveBitworkc
		}
		if !common.IsValidContainer(entity.ContainerName) {
			return errors.ErrInvalidContainer
		}
		if m.ContainerHasExist(entity.ContainerName) {
			return errors.ErrContainerHasExist
		}
		m.ensureUTXONotNil(atomicalsID)
		m.UTXOs[atomicalsID].Nft = append(m.UTXOs[atomicalsID].Nft, entity)
		m.GlobalNftContainerMap[entity.ContainerName] = make(map[string]bool, 0)
	} else if operation.Payload.Args.RequestDmitem != "" {
		// seems not necessary
		if !operation.IsValidCommitVoutIndexForNameRevel() {
			return errors.ErrInvalidCommitVoutIndex
		}
		if !operation.IsWithinAcceptableBlocksForNameReveal() {
			return errors.ErrInvalidCommitHeight
		}
		entity := &UserNftInfo{
			Dmitem:                     operation.Payload.Args.RequestDmitem,
			ParentContainerAtomicalsID: operation.Payload.Args.ParentContainer,
			Nonce:                      operation.Payload.Args.Nonce,
			Time:                       operation.Payload.Args.Time,
			Bitworkc:                   bitworkc,
			Bitworkr:                   bitworkr,
			AtomicalsID:                atomicalsID,
			Location:                   atomicalsID,
		}
		if !common.IsValidDmitem(entity.Dmitem) {
			return errors.ErrInvalidContainerDmitem
		}
		parentContainerName, ok := m.ParentContainerHasExist(entity.ParentContainerAtomicalsID)
		if !ok {
			return errors.ErrParentRealmNotExist
		}
		m.ensureUTXONotNil(atomicalsID)
		m.UTXOs[atomicalsID].Nft = append(m.UTXOs[atomicalsID].Nft, entity)
		m.GlobalNftContainerMap[parentContainerName][entity.Dmitem] = true
	} else {
		log.Log.Warnf("operation.Script:%+v", operation.Script)
		log.Log.Warnf("operation.Payload.Args:%+v", operation.Payload.Args)
	}
	return nil
}

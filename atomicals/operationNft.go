package atomicals

import (
	"log"

	"github.com/atomicals-core/atomicals/common"
	"github.com/atomicals-core/atomicals/witness"
	"github.com/atomicals-core/pkg/errors"
	"github.com/btcsuite/btcd/btcjson"
)

func (m *Atomicals) mintNft(operation *witness.WitnessAtomicalsOperation, vin btcjson.Vin, vout []btcjson.Vout, userPk string) error {
	atomicalsID := atomicalsID(operation.TxID, common.VOUT_EXPECT_OUTPUT_BYTES)
	entity := &UserNftEntity{}
	if operation.Payload.Args.RequestRealm != "" {
		entity = &UserNftEntity{
			EntityType: EntityTypeNftRealm,
			Name:       operation.Payload.Args.RequestRealm,
			Nonce:      operation.Payload.Args.Nonce,
			Time:       operation.Payload.Args.Time,
			Bitworkc:   operation.Payload.Args.Bitworkc,

			AtomicalsID: atomicalsID,
			Location:    atomicalsID,
		}
		if !common.IsValidRealm(entity.Name) {
			return errors.ErrInvalidRealm
		}
		m.UTXOs[atomicalsID] = &AtomicalsUTXO{
			AtomicalID: atomicalsID,
			Nft:        make([]*UserNftEntity, 0),
		}
		m.UTXOs[atomicalsID].Nft = append(m.UTXOs[atomicalsID].Nft, entity)
	} else if operation.Payload.Args.RequestContainer != "" {
		entity = &UserNftEntity{
			EntityType:  EntityTypeNftContainer,
			Name:        operation.Payload.Args.RequestRealm,
			Nonce:       operation.Payload.Args.Nonce,
			Time:        operation.Payload.Args.Time,
			Bitworkc:    operation.Payload.Args.Bitworkc,
			AtomicalsID: atomicalsID,
			Location:    atomicalsID,
		}
		if !common.IsValidContainer(entity.Name) {
			return errors.ErrInvalidContainer
		}
		m.UTXOs[atomicalsID] = &AtomicalsUTXO{
			AtomicalID: atomicalsID,
			Nft:        make([]*UserNftEntity, 0),
		}
		m.UTXOs[atomicalsID].Nft = append(m.UTXOs[atomicalsID].Nft, entity)
	} else if operation.Payload.Args.RequestSubRealm != "" {
		log.Printf("operation.Payload.Args:%+v", operation.Payload.Args)
	} else if operation.Payload.Args.RequestDmitem != "" {
		log.Printf("operation.Payload.Args:%+v", operation.Payload.Args)
	} else {
		log.Printf("operation.Payload.Args:%+v", operation.Payload.Args)
	}
	return nil
}

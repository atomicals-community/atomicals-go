package atomicals

import (
	"fmt"

	"github.com/atomicals-core/atomicals/witness"
)

type Atomicals struct {
	Height            int64
	TxIndex           int64
	UTXOs             map[string]*AtomicalsUTXO   // key:txID-voutIndex(atomicalsID)
	AtomicalsFtEntity map[string]*AtomicalsFtInfo // key:userPk-name
}

func NewAtomicals(height int64) *Atomicals {
	return &Atomicals{
		Height:            height,
		TxIndex:           0,
		AtomicalsFtEntity: make(map[string]*AtomicalsFtInfo, 0),
		UTXOs:             make(map[string]*AtomicalsUTXO, 0),
	}
}

func atomicalsID(txID string, voutIndex int64) string {
	return fmt.Sprintf("%v_%v", txID, voutIndex)
}

type AtomicalsUTXO struct {
	AtomicalID    string
	Nft           []*UserNftEntity // operationType: nft
	DistributedFt []*UserDmtEntity // operationType: dmt
	DirectFt      []*UserFtEntity  // operationType: ft
}

type AtomicalsFtInfo struct {
	AtomicalsID string
	Ticker      string

	MintAmount float64 // mint amount once
	MintHeight int64   // start mint height
	MaxMints   int64   // total supply
	Bitworkc   string
	Meta       *witness.Meta

	MintedAmount float64 // record minted amount
}

type UserNftEntity struct {
	UserPk      string
	AtomicalsID string // txID _ VOUT_EXPECT_OUTPUT_BYTES when be minted
	Location    string // txID_voutIndex updated after being transfered

	EntityType int // nftRealm,nftTicker,nftContainer
	Name       string

	Nonce    int64
	Time     int64
	Bitworkc string
}

const (
	EntityTypeNftRealm     = 0
	EntityTypeNftTicker    = 1
	EntityTypeNftContainer = 2
)

type UserDmtEntity struct {
	UserPk      string
	AtomicalsID string // txID _ VOUT_EXPECT_OUTPUT_BYTES when be minted
	Location    string // txID_voutIndex updated after being transfered

	Name string

	Nonce    int64
	Time     int64
	Bitworkc string
	Amount   float64
}

type UserFtEntity struct {
	UserPk      string
	AtomicalsID string // txID _ VOUT_EXPECT_OUTPUT_BYTES when be minted
	Location    string // txID_voutIndex updated after being transfered

	CommitTxID   string
	CommitIndex  int64
	CommitHeight int64

	CurrentTxID   string
	CurrentHeight int64

	Meta     *witness.Meta
	Bitworkc string

	// record
	Type          string
	Subtype       string
	RequestTicker string
	MaxSupply     int64
	MintAmount    int64
	MintHeight    int64
	MaxMints      int64
}

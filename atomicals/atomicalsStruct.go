package atomicals

import (
	"github.com/atomicals-core/atomicals/common"
	"github.com/atomicals-core/atomicals/witness"
)

type AtomicalsUTXO struct {
	AtomicalID    string
	Nft           []*UserNftInfo         // operationType: nft
	DistributedFt []*UserDistributedInfo // operationType: dmt
	DirectFt      []*UserDirectFtInfo    // operationType: ft
}

type DistributedFtInfo struct {
	AtomicalsID    string
	Ticker         string
	Type           string
	Subtype        string
	Md             string
	MintMode       string
	MaxSupply      float64
	MaxMintsGlobal float64
	MintAmount     float64 // mint amount once
	MintHeight     int64   // start mint height
	MintBitworkc   *common.Bitwork
	MintBitworkr   *common.Bitwork
	MaxMints       float64 // total supply
	Bitworkc       *common.Bitwork
	Bitworkr       *common.Bitwork
	Meta           *witness.Meta

	MintedAmount float64 // record minted amount
}

const (
	TypeNftRealm     = 0
	TypeNftSubRealm  = 1
	TypeNftContainer = 2
	TypeNftTicker    = 3
)

type UserNftInfo struct {
	UserPk      string
	AtomicalsID string // txID _ VOUT_EXPECT_OUTPUT_INDEX when be minted
	Location    string // txID_voutIndex updated after being transfered

	NftType int64

	// realm
	RealmName string

	// subRealm
	SubRealmName           string
	ClaimType              witness.NftSubrealmClaimType
	ParentRealmAtomicalsID string // ParentRealm atomicalsID

	// container
	ContainerName string

	// Dmitem
	Dmitem                     string
	ParentContainerAtomicalsID string

	Nonce    int64
	Time     int64
	Bitworkc *common.Bitwork
	Bitworkr *common.Bitwork
}

type UserDistributedInfo struct {
	UserPk      string
	AtomicalsID string // txID _ VOUT_EXPECT_OUTPUT_INDEX when be minted
	Location    string // txID_voutIndex updated after being transfered

	Name string

	Nonce    int64
	Time     int64
	Bitworkc *common.Bitwork
	Bitworkr *common.Bitwork
	Amount   float64
}

type UserDirectFtInfo struct {
	UserPk      string
	AtomicalsID string // txID _ VOUT_EXPECT_OUTPUT_INDEX when be minted
	Location    string // txID_voutIndex updated after being transfered

	Meta     *witness.Meta
	Bitworkc *common.Bitwork
	Bitworkr *common.Bitwork

	// record
	Type          string
	Subtype       string
	RequestTicker string
	MaxSupply     int64
	MintAmount    int64
	MintHeight    int64
	MaxMints      int64
}

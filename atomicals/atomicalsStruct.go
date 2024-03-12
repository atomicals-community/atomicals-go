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
	MintMode       string // emu: perpetual, fixed
	MaxMintsGlobal int64  // total mint times allowed
	MintAmount     int64  // mint amount once
	MaxMints       int64  // # In the fixed mode there is a max number of mints allowed and then no more, only used when mintMode="fixed"
	MaxSupply      int64  // total supply = MaxMintsGlobal*MintAmount
	MintHeight     int64  // start mint height
	MintedTimes    int64  // record minted times
	MintBitworkc   *common.Bitwork
	MintBitworkr   *common.Bitwork
	Bitworkc       *common.Bitwork
	Bitworkr       *common.Bitwork
	Meta           *witness.Meta
	Md             string // emu:"", "0", "1"
	Bv             string // mint_info['$mint_bitwork_vec'] = bv
	Bci            string // mint_info['$mint_bitworkc_inc'] = bci
	Bri            string // mint_info['$mint_bitworkr_inc'] = bri
	Bcs            int64  // mint_info['$mint_bitworkc_start'] = bcs
	Brs            int64  // mint_info['$mint_bitworkr_start'] = brs
	Maxg           int64
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

	Nonce           int64
	Time            int64
	Bitworkc        *common.Bitwork
	Bitworkr        *common.Bitwork
	MintBitworkVec  *common.Bitwork
	MintBitworkcInc *common.Bitwork
	MintBitworkrInc *common.Bitwork
	Amount          int64
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

package postsql

import (
	"github.com/atomicals-go/atomicals-core/witness"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

const UserNftInfoTableName = "atomicals_utxo_nft"

type NftType string

const (
	TypeNftRealm     NftType = "Realm"
	TypeNftSubRealm  NftType = "SubRealm"
	TypeNftContainer NftType = "NftContainer"
	TypeNftItem      NftType = "NftItem"
)

type UTXONftInfo struct {
	gorm.Model
	UserPk      string
	AtomicalsID string `gorm:"index"` // txID _ VOUT_EXPECT_OUTPUT_INDEX when be minted
	LocationID  string `gorm:"index"` // txID_voutIndex updated after being transfered

	// realm
	RealmName string `gorm:"index"`
	// subRealm
	SubRealmName           string `gorm:"index"`
	ClaimType              witness.NftSubrealmClaimType
	ParentRealmAtomicalsID string // ParentRealm atomicalsID

	// container
	ContainerName string `gorm:"index"`
	// Dmitem
	Dmitem                     string `gorm:"index"`
	ParentContainerAtomicalsID string

	Nonce    int64
	Time     int64
	Bitworkc string
	Bitworkr string
}

func (*UTXONftInfo) TableName() string {
	return UserNftInfoTableName
}

func (*UTXONftInfo) Init(db *gorm.DB) {
	var err error
	dmodel := newDefaultModel(UserNftInfoTableName, db)
	err = dmodel.DropTable()
	assert.Nil(nil, err)
	err = dmodel.CreateTable(&UTXONftInfo{})
	assert.Nil(nil, err)
}

func (*UTXONftInfo) AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(UTXONftInfo{})
	assert.Nil(nil, err)
}

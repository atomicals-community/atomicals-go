package postsql

import (
	"github.com/atomicals-core/atomicals/witness"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

const UserNftInfoTableName = "atomicals_user_nft_info"

type UserNftInfo struct {
	gorm.Model
	UserPk      string
	AtomicalsID string `gorm:"index"` // txID _ VOUT_EXPECT_OUTPUT_INDEX when be minted
	LocationID  string `gorm:"index"` // txID_voutIndex updated after being transfered

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

	Nonce int64
	Time  int64
}

func (*UserNftInfo) TableName() string {
	return UserNftInfoTableName
}

func (*UserNftInfo) Init(db *gorm.DB) {
	var err error
	dmodel := newDefaultModel(UserNftInfoTableName, db)
	err = dmodel.DropTable()
	assert.Nil(nil, err)
	err = dmodel.CreateTable(&UserNftInfo{})
	assert.Nil(nil, err)
}

func (*UserNftInfo) AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(UserNftInfo{})
	assert.Nil(nil, err)
}

package db

import (
	"github.com/atomicals-core/atomicals/DB/postsql"
	"gorm.io/gorm"
)

type Postgres struct {
	*gorm.DB
	SQLRaw string

	UserNftInfoCache map[string][]*UserNftInfo // key: locationID
	UserFtInfoCache  map[string][]*UserFtInfo  // key: locationID
}

func (m *Postgres) UpdateLocation(height int64, txID string) error {
	entity := &postsql.Location{}
	dbTx := m.Take(&entity)
	if dbTx.Error != nil && dbTx.Error != gorm.ErrRecordNotFound {
		return dbTx.Error
	}
	entity.Owner = "atomicals"
	entity.Height = height
	entity.TxID = txID
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Save(entity)
	})
	m.SQLRaw = m.SQLRaw + sql + ";"
	dbTx = m.Exec(m.SQLRaw)
	if dbTx.Error != nil {
		return dbTx.Error
	}
	m.SQLRaw = ""
	return nil
}

func (m *Postgres) InsertNftUTXO(UTXO *UserNftInfo) error {
	// update TXCache first
	_, ok := m.UserNftInfoCache[UTXO.LocationID]
	if !ok {
		m.UserNftInfoCache[UTXO.LocationID] = make([]*UserNftInfo, 0)
	}
	m.UserNftInfoCache[UTXO.LocationID] = append(m.UserNftInfoCache[UTXO.LocationID], UTXO)

	m.SQLRaw += m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Save(&postsql.UserNftInfo{
			UserPk:                     UTXO.UserPk,
			AtomicalsID:                UTXO.AtomicalsID,
			LocationID:                 UTXO.LocationID,
			RealmName:                  UTXO.RealmName,
			SubRealmName:               UTXO.SubRealmName,
			ClaimType:                  UTXO.ClaimType,
			ParentRealmAtomicalsID:     UTXO.ParentRealmAtomicalsID,
			ContainerName:              UTXO.ContainerName,
			Dmitem:                     UTXO.Dmitem,
			ParentContainerAtomicalsID: UTXO.ParentContainerAtomicalsID,
			Nonce:                      UTXO.Nonce,
			Time:                       UTXO.Time,
		})
	}) + ";"
	return nil
}

func (m *Postgres) InsertGlobalNft(nftType int64, name, subName string) error {
	m.SQLRaw += m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Save(&postsql.GlobalNft{
			NftType: nftType,
			Name:    name,
			SubName: subName,
		})
	}) + ";"
	return nil
}

func (m *Postgres) TransferNftUTXO(oldLocationID, newLocationID, newUserPk string) error {
	m.SQLRaw += m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Where("loaction_id = ?", oldLocationID).Updates(map[string]interface{}{"loaction_id": newLocationID, "user_pk": newUserPk})
	}) + ";"
	return nil
}

func (m *Postgres) InsertFtUTXO(UTXO *UserFtInfo) error {
	// update TXCache first
	_, ok := m.UserFtInfoCache[UTXO.LocationID]
	if !ok {
		m.UserFtInfoCache[UTXO.LocationID] = make([]*UserFtInfo, 0)
	}
	m.UserFtInfoCache[UTXO.LocationID] = append(m.UserFtInfoCache[UTXO.LocationID], UTXO)

	m.SQLRaw += m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Save(&postsql.UserFtInfo{
			UserPk:        UTXO.UserPk,
			AtomicalsID:   UTXO.AtomicalsID,
			LocationID:    UTXO.LocationID,
			MintTicker:    UTXO.MintTicker,
			Nonce:         UTXO.Nonce,
			Time:          UTXO.Time,
			Amount:        UTXO.Amount,
			Type:          UTXO.Type,
			Subtype:       UTXO.Subtype,
			RequestTicker: UTXO.RequestTicker,
			MaxSupply:     UTXO.MaxSupply,
			MintAmount:    UTXO.MintAmount,
			MintHeight:    UTXO.MintHeight,
			MaxMints:      UTXO.MaxMints,
		})
	}) + ";"
	return nil
}

func (m *Postgres) DeleteFtUTXO(locationID string) error {
	m.SQLRaw += m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Where("loaction_id = ?", locationID).Delete(&postsql.UserFtInfo{})
	}) + ";"
	return nil
}

func (m *Postgres) UpdateDistributedFtAmount(tickerName string, mintTimes int64) error {
	m.SQLRaw += m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(postsql.GlobalDistributedFt{}).Where("ticker_name = ?", tickerName).Update("minted_times", mintTimes)
	}) + ";"
	return nil
}

func (m *Postgres) InsertDirectFt(tickerName string) error {
	m.SQLRaw += m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Save(&postsql.GlobalDirectFt{
			TickerName: tickerName,
		})
	}) + ";"
	return nil
}

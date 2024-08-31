package repo

import (
	"fmt"

	"github.com/atomicals-go/repo/postsql"
	"gorm.io/gorm"
)

type Postgres struct {
	*gorm.DB
	bloomFilter map[string]*bloomFilterInfo
}
type AtomicaslData struct {
	Op                     string
	Description            string
	Mod                    *postsql.ModInfo
	Dat                    *postsql.DatInfo
	DeleteFts              []*postsql.UTXOFtInfo
	NewFts                 []*postsql.UTXOFtInfo
	UpdateNfts             []*postsql.UTXONftInfo
	NewUTXOFtInfo          *postsql.UTXOFtInfo
	UpdateDistributedFt    *postsql.GlobalDistributedFt
	NewGlobalDistributedFt *postsql.GlobalDistributedFt
	NewGlobalDirectFt      *postsql.GlobalDirectFt
	NewUTXONftInfo         *postsql.UTXONftInfo
	DeleteUTXONfts         []*postsql.UTXONftInfo
}

func (m *AtomicaslData) ParseOperation(orgOp string) {
	if m == nil {
		return
	}

	// transfer ft
	if len(m.DeleteFts) > 0 {
		for _, v := range m.DeleteFts {
			m.Description += fmt.Sprintf("delete#ticker:%v,locationID:%v,userPk:%v,amount:%v\n", v.MintTicker, v.LocationID, v.UserPk, v.Amount)
		}
		m.Op = "transfer"
		if orgOp == "x" {
			m.Op = "splat"
		} else if orgOp == "y" {
			m.Op = "split"
		}
	}
	if len(m.NewFts) > 0 {
		for _, v := range m.NewFts {
			m.Description += fmt.Sprintf("insert#ticker:%v,locationID:%v,userPk:%v,amount:%v\n", v.MintTicker, v.LocationID, v.UserPk, v.Amount)
		}
	}

	// transfer nft
	if len(m.UpdateNfts) > 0 {
		m.Op = "transfer"
		if orgOp == "x" {
			m.Op = "splat"
		} else if orgOp == "y" {
			m.Op = "split"
		}
	}

	// mod
	if m.Mod != nil {
		m.Op = "mod"
	}

	if m.Dat != nil {
		m.Op = "dat"
	}

	// mint ft
	if m.NewUTXOFtInfo != nil {
		m.Op = "mint-dft"
		m.Description += fmt.Sprintf("insert#ticker:%v,locationID:%v,userPk:%v,amount:%v\n", m.NewUTXOFtInfo.MintTicker, m.NewUTXOFtInfo.LocationID, m.NewUTXOFtInfo.UserPk, m.NewUTXOFtInfo.Amount)
	} else {
		if orgOp == "dmt" {
			m.Op = "mint-dft-failed"
		}
	}
	if m.UpdateDistributedFt != nil {
		m.Op = "dmt"
	}
	if m.NewGlobalDistributedFt != nil {
		m.Op = "dft"
	}
	if m.NewGlobalDirectFt != nil {
		m.Op = "ft"
	}

	// mint nft
	if m.NewUTXONftInfo != nil {
		if m.NewUTXONftInfo.RealmName != "" {
			m.Op = "mint-nft-realm"
		}
		if m.NewUTXONftInfo.SubRealmName != "" {
			m.Op = "mint-nft-subrealm"
		}
		if m.NewUTXONftInfo.ContainerName != "" {
			m.Op = "mint-nft-container"
		}
		if m.NewUTXONftInfo.Dmitem != "" {
			m.Op = "mint-nft"
		}
	}
}

func (m *Postgres) UpdateDB(location *postsql.Location, data *AtomicaslData) error {
	if data == nil {
		return nil
	}

	if !((data.Op != "") || (location.BlockHeight%10 == 0 && location.TxIndex == 0)) {
		return nil
	}

	err := m.Transaction(func(tx *gorm.DB) error {
		// mod
		if data.Mod != nil {
			dbErr := tx.Save(data.Mod)
			if dbErr.Error != nil {
				return dbErr.Error
			}
		}

		if data.Dat != nil {
			dbErr := tx.Save(data.Dat)
			if dbErr.Error != nil {
				return dbErr.Error
			}
		}

		// transfer ft
		if len(data.DeleteFts) > 0 {
			locationIDs := make([]string, len(data.DeleteFts))
			for i, v := range data.DeleteFts {
				locationIDs[i] = v.LocationID
			}
			if dbErr := tx.Model(&postsql.UTXOFtInfo{}).Unscoped().Where("location_id IN ?", locationIDs).Delete(&postsql.UTXOFtInfo{}).Error; dbErr != nil {
				return dbErr
			}
		}
		if len(data.NewFts) > 0 {
			for _, v := range data.NewFts {
				m.addFtLocationID(v.LocationID)
			}
			if dbErr := tx.Create(&data.NewFts).Error; dbErr != nil {
				return dbErr
			}
		}

		// transfer nft
		if len(data.UpdateNfts) > 0 {
			for _, v := range data.UpdateNfts {
				m.addNftLocationID(v.LocationID)
			}
			if dbErr := tx.Save(&data.UpdateNfts).Error; dbErr != nil {
				return dbErr
			}
		}

		// mint ft
		if data.NewUTXOFtInfo != nil {
			m.addFtLocationID(data.NewUTXOFtInfo.LocationID)
			if dbErr := tx.Save(data.NewUTXOFtInfo).Error; dbErr != nil {
				return dbErr
			}
		}
		if data.UpdateDistributedFt != nil {
			dbErr := tx.Model(postsql.GlobalDistributedFt{}).Where("ticker_name = ?", data.UpdateDistributedFt.TickerName).Updates(map[string]interface{}{"minted_times": data.UpdateDistributedFt.MintedTimes})
			if dbErr.Error != nil {
				return dbErr.Error
			}
		}
		if data.NewGlobalDistributedFt != nil {
			// m.addDistributedFt(data.NewGlobalDistributedFt.TickerName)
			m.addFtLocationID(data.NewGlobalDistributedFt.LocationID)
			dbErr := tx.Save(data.NewGlobalDistributedFt)
			if dbErr.Error != nil {
				return dbErr.Error
			}
		}
		if data.NewGlobalDirectFt != nil {
			m.addFtLocationID(data.NewGlobalDirectFt.LocationID)
			dbErr := tx.Save(data.NewGlobalDirectFt)
			if dbErr.Error != nil {
				return dbErr.Error
			}
		}

		// mint nft
		if data.NewUTXONftInfo != nil {
			if data.NewUTXONftInfo.RealmName != "" {
				m.addRealm(data.NewUTXONftInfo.RealmName)
			} else if data.NewUTXONftInfo.ContainerName != "" {
				m.addContainer(data.NewUTXONftInfo.ContainerName)
			}
			m.addNftLocationID(data.NewUTXONftInfo.LocationID)
			dbErr := tx.Save(data.NewUTXONftInfo)
			if dbErr.Error != nil {
				return dbErr.Error
			}
		}

		if len(data.DeleteUTXONfts) != 0 {
			dbErr := tx.Delete(&data.DeleteUTXONfts)
			if dbErr.Error != nil {
				return dbErr.Error
			}
		}

		// update bloom filter
		for name, v := range m.bloomFilter {
			if v.needUpdate {
				data, err := v.filter.MarshalJSON()
				if err != nil {
					return err
				}
				dbErr := tx.Model(postsql.BloomFilter{}).Where("name = ?", name).Update("data", data)
				if dbErr.Error != nil {
					return dbErr.Error
				}
				v.needUpdate = false
			}
		}

		// insert btc tx record
		if data.Op != "" {
			dbErr := tx.Save(&postsql.AtomicalsTx{
				BlockHeight: location.BlockHeight,
				TxIndex:     location.TxIndex,
				TxID:        location.Txid,
				Operation:   data.Op,
				Description: data.Description,
			})
			if dbErr.Error != nil {
				return dbErr.Error
			}
		}

		// update location
		dbErr := tx.Model(postsql.Location{}).Where("key = ?", postsql.LocationKey).Save(location)
		if dbErr.Error != nil {
			return dbErr.Error
		}

		// we don't need save all height-txid in db, delete atomicals tx until
		// if currentTxIndex == 0 {
		// 	dbErr := tx.Model(postsql.AtomicalsTx{}).Unscoped().Where("block_height = ? and operation = ?", currentHeight-utils.MINT_GENERAL_COMMIT_REVEAL_DELAY_BLOCKS, "").Delete(&postsql.AtomicalsTx{})
		// 	if dbErr.Error != nil {
		// 		return dbErr.Error
		// 	}
		// }
		return nil
	})
	return err
}

func (m *Postgres) PostgresDB() *Postgres {
	return m
}

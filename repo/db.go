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
	DeleteFts              []*postsql.UTXOFtInfo
	NewFts                 []*postsql.UTXOFtInfo
	UpdateNfts             []*postsql.UTXONftInfo
	NewUTXOFtInfo          *postsql.UTXOFtInfo
	UpdateDistributedFt    *postsql.GlobalDistributedFt
	NewGlobalDistributedFt *postsql.GlobalDistributedFt
	NewGlobalDirectFt      *postsql.GlobalDirectFt
	NewUTXONftInfo         *postsql.UTXONftInfo
}

func (m *AtomicaslData) ParseOperation() {
	// mod
	if m.Mod != nil {
		m.Op = "mod"
	}

	// transfer ft
	if len(m.DeleteFts) > 0 {
		for _, v := range m.DeleteFts {
			m.Description += fmt.Sprintf("delete#ticker:%v,locationID:%v,userPk:%v,amount:%v\n", v.MintTicker, v.LocationID, v.UserPk, v.Amount)
		}
		m.Op = "transfer_ft"
	}
	if len(m.NewFts) > 0 {
		for _, v := range m.NewFts {
			m.Description += fmt.Sprintf("insert#ticker:%v,locationID:%v,userPk:%v,amount:%v\n", v.MintTicker, v.LocationID, v.UserPk, v.Amount)
		}
		m.Op = "transfer_ft"
	}

	// transfer nft
	if len(m.UpdateNfts) > 0 {
		m.Op = "transfer_nft"
	}

	// mint ft
	if m.NewUTXOFtInfo != nil {
		m.Op = "dmt"
		m.Description += fmt.Sprintf("insert#ticker:%v,locationID:%v,userPk:%v,amount:%v\n", m.NewUTXOFtInfo.MintTicker, m.NewUTXOFtInfo.LocationID, m.NewUTXOFtInfo.UserPk, m.NewUTXOFtInfo.Amount)
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
		m.Op = "nft"
	}
}

func (m *Postgres) UpdateDB(
	currentHeight, currentTxIndex int64, txID string,
	data *AtomicaslData,
) error {
	if data == nil {
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
			m.addDistributedFt(data.NewGlobalDistributedFt.TickerName)
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
				BlockHeight: currentHeight,
				TxIndex:     currentTxIndex,
				TxID:        txID,
				Operation:   data.Op,
				Description: data.Description,
			})
			if dbErr.Error != nil {
				return dbErr.Error
			}
		}

		// update location
		dbErr := tx.Model(postsql.Location{}).Where("name = ?", "atomicals").Updates(map[string]interface{}{"block_height": currentHeight, "tx_index": currentTxIndex})
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

package repo

import (
	"fmt"

	"github.com/atomicals-go/repo/postsql"
	"github.com/atomicals-go/utils"
	"gorm.io/gorm"
)

type Postgres struct {
	*gorm.DB
	bloomFilter map[string]*bloomFilterInfo
}
type AtomicaslData struct {
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

func (m *Postgres) UpdateDB(
	currentHeight, currentTxIndex int64, txID string,
	data *AtomicaslData,
) error {
	if data == nil {
		return nil
	}
	op := ""
	description := ""
	err := m.Transaction(func(tx *gorm.DB) error {
		// mod
		if data.Mod != nil {
			dbErr := tx.Save(data.Mod)
			if dbErr.Error != nil {
				return dbErr.Error
			}
			op = "mod"
		}

		// transfer ft
		if len(data.DeleteFts) > 0 {
			locationIDs := make([]string, len(data.DeleteFts))
			for i, v := range data.DeleteFts {
				locationIDs[i] = v.LocationID
				description += fmt.Sprintf("delete#ticker:%v,locationID:%v,userPk:%v,amount:%v\n", v.MintTicker, v.LocationID, v.UserPk, v.Amount)
			}
			if dbErr := tx.Model(&postsql.UTXOFtInfo{}).Unscoped().Where("location_id IN ?", locationIDs).Delete(&postsql.UTXOFtInfo{}).Error; dbErr != nil {
				return dbErr
			}
			op = "transfer_ft"
		}
		if len(data.NewFts) > 0 {
			for _, v := range data.NewFts {
				m.addFtLocationID(v.LocationID)
				description += fmt.Sprintf("insert#ticker:%v,locationID:%v,userPk:%v,amount:%v\n", v.MintTicker, v.LocationID, v.UserPk, v.Amount)
			}
			if dbErr := tx.Create(&data.NewFts).Error; dbErr != nil {
				return dbErr
			}
			op = "transfer_ft"
		}

		// transfer nft
		if len(data.UpdateNfts) > 0 {
			for _, v := range data.UpdateNfts {
				m.addNftLocationID(v.LocationID)
			}
			if dbErr := tx.Save(&data.UpdateNfts).Error; dbErr != nil {
				return dbErr
			}
			op = "transfer_nft"
		}

		// mint ft
		if data.NewUTXOFtInfo != nil {
			m.addFtLocationID(data.NewUTXOFtInfo.LocationID)
			if dbErr := tx.Save(data.NewUTXOFtInfo).Error; dbErr != nil {
				return dbErr
			}
			op = "dmt"
			description += fmt.Sprintf("insert#ticker:%v,locationID:%v,userPk:%v,amount:%v\n", data.NewUTXOFtInfo.MintTicker, data.NewUTXOFtInfo.LocationID, data.NewUTXOFtInfo.UserPk, data.NewUTXOFtInfo.Amount)
		}
		if data.UpdateDistributedFt != nil {
			dbErr := tx.Model(postsql.GlobalDistributedFt{}).Where("ticker_name = ?", data.UpdateDistributedFt.TickerName).Updates(map[string]interface{}{"minted_times": data.UpdateDistributedFt.MintedTimes})
			if dbErr.Error != nil {
				return dbErr.Error
			}
			op = "dmt"
		}
		if data.NewGlobalDistributedFt != nil {
			m.addDistributedFt(data.NewGlobalDistributedFt.TickerName)
			m.addFtLocationID(data.NewGlobalDistributedFt.LocationID)
			dbErr := tx.Save(data.NewGlobalDistributedFt)
			if dbErr.Error != nil {
				return dbErr.Error
			}
			op = "dft"
		}
		if data.NewGlobalDirectFt != nil {
			m.addFtLocationID(data.NewGlobalDirectFt.LocationID)
			dbErr := tx.Save(data.NewGlobalDirectFt)
			if dbErr.Error != nil {
				return dbErr.Error
			}
			op = "ft"
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
			op = "nft"
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

		if op != "" {
			// insert btc tx record
			dbErr := tx.Save(&postsql.AtomicalsTx{
				BlockHeight: currentHeight,
				TxIndex:     currentTxIndex,
				TxID:        txID,
				Operation:   op,
				Description: description,
			})
			if dbErr.Error != nil {
				return dbErr.Error
			}
		}

		// update location
		// if op != "" || (currentHeight%10 == 0 && currentTxIndex == 0) {
		dbErr := tx.Model(postsql.Location{}).Where("name = ?", "atomicals").Updates(map[string]interface{}{"block_height": currentHeight, "tx_index": currentTxIndex})
		if dbErr.Error != nil {
			return dbErr.Error
		}
		// }

		// we don't need save all height-txid in db, delete atomicals tx until
		if currentTxIndex == 0 {
			dbErr := tx.Model(postsql.AtomicalsTx{}).Unscoped().Where("block_height = ? and operation = ?", currentHeight-utils.MINT_GENERAL_COMMIT_REVEAL_DELAY_BLOCKS, "").Delete(&postsql.AtomicalsTx{})
			if dbErr.Error != nil {
				return dbErr.Error
			}
		}
		return nil
	})
	return err
}

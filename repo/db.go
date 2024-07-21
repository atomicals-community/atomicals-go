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

func (m *Postgres) UpdateDB(
	currentHeight, currentTxIndex int64, txID string,
	mod *postsql.ModInfo,
	deleteFts []*postsql.UTXOFtInfo, newFts []*postsql.UTXOFtInfo,
	updateNfts []*postsql.UTXONftInfo,
	newUTXOFtInfo *postsql.UTXOFtInfo, updateDistributedFt *postsql.GlobalDistributedFt,
	newGlobalDistributedFt *postsql.GlobalDistributedFt,
	newGlobalDirectFt *postsql.GlobalDirectFt,
	newUTXONftInfo *postsql.UTXONftInfo) error {
	op := ""
	description := ""
	err := m.Transaction(func(tx *gorm.DB) error {
		// mod
		if mod != nil {
			dbErr := tx.Save(mod)
			if dbErr.Error != nil {
				return dbErr.Error
			}
			op = "mod"
		}

		// transfer ft
		if len(deleteFts) > 0 {
			locationIDs := make([]string, len(deleteFts))
			for i, v := range deleteFts {
				locationIDs[i] = v.LocationID
				description += fmt.Sprintf("delete#ticker:%v,locationID:%v,userPk:%v,amount:%v\n", v.MintTicker, v.LocationID, v.UserPk, v.Amount)
			}
			if dbErr := tx.Model(&postsql.UTXOFtInfo{}).Unscoped().Where("location_id IN ?", locationIDs).Delete(&postsql.UTXOFtInfo{}).Error; dbErr != nil {
				return dbErr
			}
			op = "transfer_ft"
		}
		if len(newFts) > 0 {
			for _, v := range newFts {
				m.addFtLocationID(v.LocationID)
				description += fmt.Sprintf("insert#ticker:%v,locationID:%v,userPk:%v,amount:%v\n", v.MintTicker, v.LocationID, v.UserPk, v.Amount)
			}
			if dbErr := tx.Create(&newFts).Error; dbErr != nil {
				return dbErr
			}
			op = "transfer_ft"
		}

		// transfer nft
		if len(updateNfts) > 0 {
			for _, v := range updateNfts {
				m.addNftLocationID(v.LocationID)
			}
			if dbErr := tx.Save(&updateNfts).Error; dbErr != nil {
				return dbErr
			}
			op = "transfer_nft"
		}

		// mint ft
		if newUTXOFtInfo != nil {
			m.addFtLocationID(newUTXOFtInfo.LocationID)
			if dbErr := tx.Save(newUTXOFtInfo).Error; dbErr != nil {
				return dbErr
			}
			op = "dmt"
			description += fmt.Sprintf("insert#ticker:%v,locationID:%v,userPk:%v,amount:%v\n", newUTXOFtInfo.MintTicker, newUTXOFtInfo.LocationID, newUTXOFtInfo.UserPk, newUTXOFtInfo.Amount)
		}
		if updateDistributedFt != nil {
			dbErr := tx.Model(postsql.GlobalDistributedFt{}).Where("ticker_name = ?", updateDistributedFt.TickerName).Updates(map[string]interface{}{"minted_times": updateDistributedFt.MintedTimes})
			if dbErr.Error != nil {
				return dbErr.Error
			}
			op = "dmt"
		}
		if newGlobalDistributedFt != nil {
			m.addDistributedFt(newGlobalDistributedFt.TickerName)
			m.addFtLocationID(newGlobalDistributedFt.LocationID)
			dbErr := tx.Save(newGlobalDistributedFt)
			if dbErr.Error != nil {
				return dbErr.Error
			}
			op = "dft"
		}
		if newGlobalDirectFt != nil {
			m.addFtLocationID(newGlobalDirectFt.LocationID)
			dbErr := tx.Save(newGlobalDirectFt)
			if dbErr.Error != nil {
				return dbErr.Error
			}
			op = "ft"
		}

		// mint nft
		if newUTXONftInfo != nil {
			if newUTXONftInfo.RealmName != "" {
				m.addRealm(newUTXONftInfo.RealmName)
			} else if newUTXONftInfo.ContainerName != "" {
				m.addContainer(newUTXONftInfo.ContainerName)
			}
			m.addNftLocationID(newUTXONftInfo.LocationID)
			dbErr := tx.Save(newUTXONftInfo)
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
			}
		}
		for _, v := range m.bloomFilter {
			v.needUpdate = false
		}

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

		// we don't need save all height-txid in db, delete atomicals tx until
		if currentTxIndex == 0 {
			dbErr = tx.Model(postsql.AtomicalsTx{}).Unscoped().Where("block_height = ? and operation = ?", currentHeight-utils.MINT_GENERAL_COMMIT_REVEAL_DELAY_BLOCKS, "").Delete(&postsql.AtomicalsTx{})
			if dbErr.Error != nil {
				return dbErr.Error
			}
		}

		// update location
		dbErr = tx.Model(postsql.Location{}).Where("name = ?", "atomicals").Updates(map[string]interface{}{"block_height": currentHeight, "tx_index": currentTxIndex})
		if dbErr.Error != nil {
			return dbErr.Error
		}

		return nil
	})
	return err
}

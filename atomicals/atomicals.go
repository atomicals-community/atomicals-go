package atomicals

import (
	"fmt"

	"github.com/atomicals-core/pkg/btcsync"
)

type Atomicals struct {
	btcClient *btcsync.BtcSync
	Height    int64
	// TxIndex           int64

	NftUTXOsByAtomicalsID map[string][]*UserNftInfo  // key: txID-voutIndex(atomicalsID, generate when be minted or deployed)
	NftUTXOsByLocationID  map[string][]*UserNftInfo  // key: txID-voutIndex(locationID, change when be transfered, we have to make two index for nft,cause )
	GlobalNftRealmMap     map[string]map[string]bool // operation: nft, key: realmName key: subRealmName
	GlobalNftContainerMap map[string]map[string]bool // operation: nft, key: containerName key: dmitem

	FtUTXOs                map[string][]*UserFtInfo      // key: txID-voutIndex(locationID, change when be transfered)
	GlobalDistributedFtMap map[string]*DistributedFtInfo // operation: dmt, key: tickerName
	GlobalDirectFtMap      map[string]bool               // operation: ft, key: tickerName

}

func NewAtomicals(btcClient *btcsync.BtcSync, height int64) *Atomicals {
	return &Atomicals{
		btcClient: btcClient,
		Height:    height,
		// TxIndex:           0,
		NftUTXOsByAtomicalsID: make(map[string][]*UserNftInfo, 0),
		NftUTXOsByLocationID:  make(map[string][]*UserNftInfo, 0),
		GlobalNftRealmMap:     make(map[string]map[string]bool, 0),
		GlobalNftContainerMap: make(map[string]map[string]bool, 0),

		FtUTXOs:                make(map[string][]*UserFtInfo, 0),
		GlobalDistributedFtMap: make(map[string]*DistributedFtInfo, 0),
		GlobalDirectFtMap:      make(map[string]bool, 0),
	}
}

func (m *Atomicals) ensureNftUTXONotNil(atomicalsID string) {
	if _, ok := m.NftUTXOsByAtomicalsID[atomicalsID]; !ok {
		m.NftUTXOsByAtomicalsID[atomicalsID] = make([]*UserNftInfo, 0)
	}
	if _, ok := m.NftUTXOsByLocationID[atomicalsID]; !ok {
		m.NftUTXOsByLocationID[atomicalsID] = make([]*UserNftInfo, 0)
	}
}

func (m *Atomicals) ensureFtUTXONotNil(atomicalsID string) {
	if _, ok := m.FtUTXOs[atomicalsID]; !ok {
		m.FtUTXOs[atomicalsID] = make([]*UserFtInfo, 0)
	}
}

func (m *Atomicals) RealmHasExist(realmName string) bool {
	if _, ok := m.GlobalNftRealmMap[realmName]; !ok {
		return false
	}
	return true
}

func (m *Atomicals) ParentRealmHasExist(parentRealmAtomicalsID string) (string, bool) {
	if _, ok := m.NftUTXOsByAtomicalsID[parentRealmAtomicalsID]; !ok {
		return "", false
	}
	// TODO: need to ensure if nft will be merged when they're transferred
	if len(m.NftUTXOsByAtomicalsID[parentRealmAtomicalsID]) == 0 {
		return "", false
	}
	parentRealmName := m.NftUTXOsByAtomicalsID[parentRealmAtomicalsID][0].RealmName
	if _, ok := m.GlobalNftRealmMap[parentRealmName]; !ok {
		panic("GlobalNftRealmMap and NftUTXOsByAtomicalsID are not match")
	}
	return parentRealmName, true
}

func (m *Atomicals) SubRealmHasExist(parentRealmName, subRealmName string) bool {
	if _, ok := m.GlobalNftRealmMap[parentRealmName][subRealmName]; !ok {
		return false
	}
	return true
}

func (m *Atomicals) ContainerHasExist(containerName string) bool {
	if _, ok := m.GlobalNftContainerMap[containerName]; ok {
		return true
	}
	return false
}

func (m *Atomicals) ParentContainerHasExist(parentContainerAtomicalsID string) (string, bool) {
	if _, ok := m.NftUTXOsByAtomicalsID[parentContainerAtomicalsID]; !ok {
		return "", false
	}
	// TODO: need to ensure if nft will be merged when they're transferred
	if len(m.NftUTXOsByAtomicalsID[parentContainerAtomicalsID]) == 0 {
		return "", false
	}
	parentContainerName := m.NftUTXOsByAtomicalsID[parentContainerAtomicalsID][0].ContainerName
	if _, ok := m.GlobalNftContainerMap[parentContainerName]; !ok {
		panic("GlobalNftContainerMap and NftUTXOsByAtomicalsID are not match")
	}
	return parentContainerName, true
}
func (m *Atomicals) DistributedFtHasExist(tickerName string) bool {
	if _, ok := m.GlobalDistributedFtMap[tickerName]; ok {
		return true
	}
	return false
}

func atomicalsID(txID string, voutIndex int64) string {
	return fmt.Sprintf("%v_%v", txID, voutIndex)
}

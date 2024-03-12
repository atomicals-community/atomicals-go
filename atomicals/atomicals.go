package atomicals

import (
	"fmt"

	"github.com/atomicals-core/pkg/btcsync"
)

type Atomicals struct {
	btcClient *btcsync.BtcSync
	Height    int64
	// TxIndex           int64
	UTXOs                  map[string]*AtomicalsUTXO     // key: txID-voutIndex(atomicalsID)
	GlobalDistributedFtMap map[string]*DistributedFtInfo // key: tickerName
	GlobalDirectFtMap      map[string]bool               // key: tickerName
	GlobalNftRealmMap      map[string]map[string]bool    // key: realmName key: subRealmName
	GlobalNftContainerMap  map[string]map[string]bool    // key: containerName key: dmitem
}

func NewAtomicals(btcClient *btcsync.BtcSync, height int64) *Atomicals {
	return &Atomicals{
		btcClient: btcClient,
		Height:    height,
		// TxIndex:           0,
		UTXOs:                  make(map[string]*AtomicalsUTXO, 0),
		GlobalDistributedFtMap: make(map[string]*DistributedFtInfo, 0),
		GlobalDirectFtMap:      make(map[string]bool, 0),
		GlobalNftRealmMap:      make(map[string]map[string]bool, 0),
		GlobalNftContainerMap:  make(map[string]map[string]bool, 0),
	}
}

func (m *Atomicals) ensureUTXONotNil(atomicalsID string) {
	if _, ok := m.UTXOs[atomicalsID]; !ok {
		m.UTXOs[atomicalsID] = &AtomicalsUTXO{
			AtomicalID: atomicalsID,
		}
	}
	if m.UTXOs[atomicalsID].Nft == nil {
		m.UTXOs[atomicalsID].Nft = make([]*UserNftInfo, 0)
	}
	if m.UTXOs[atomicalsID].DistributedFt == nil {
		m.UTXOs[atomicalsID].DistributedFt = make([]*UserDistributedInfo, 0)
	}
	if m.UTXOs[atomicalsID].DirectFt == nil {
		m.UTXOs[atomicalsID].DirectFt = make([]*UserDirectFtInfo, 0)
	}
}

func (m *Atomicals) RealmHasExist(realmName string) bool {
	if _, ok := m.GlobalNftRealmMap[realmName]; !ok {
		return false
	}
	return true
}

func (m *Atomicals) ParentRealmHasExist(parentRealmAtomicalsID string) (string, bool) {
	if _, ok := m.UTXOs[parentRealmAtomicalsID]; !ok {
		return "", false
	}
	// TODO: need to ensure if nft will be merged when they're transferred
	if len(m.UTXOs[parentRealmAtomicalsID].Nft) == 0 {
		return "", false
	}
	parentRealmName := m.UTXOs[parentRealmAtomicalsID].Nft[0].RealmName
	if _, ok := m.GlobalNftRealmMap[parentRealmName]; !ok {
		panic("GlobalNftRealmMap and UTXOs are not match")
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
	if _, ok := m.UTXOs[parentContainerAtomicalsID]; !ok {
		return "", false
	}
	// TODO: need to ensure if nft will be merged when they're transferred
	if len(m.UTXOs[parentContainerAtomicalsID].Nft) == 0 {
		return "", false
	}
	parentContainerName := m.UTXOs[parentContainerAtomicalsID].Nft[0].ContainerName
	if _, ok := m.GlobalNftContainerMap[parentContainerName]; !ok {
		panic("GlobalNftContainerMap and UTXOs are not match")
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

package db

type Memory struct {
	height int64
	txID   string

	nftUTXOsByAtomicalsID map[string][]*UserNftInfo  // key: txID-voutIndex(atomicalsID, generate when be minted or deployed)
	nftUTXOsByLocationID  map[string][]*UserNftInfo  // key: txID-voutIndex(locationID, change when be transfered, we have to make two index for nft,cause )
	globalNftRealmMap     map[string]map[string]bool // operation: nft, key: realmName key: subRealmName
	globalNftContainerMap map[string]map[string]bool // operation: nft, key: containerName key: dmitem

	ftUTXOs                map[string][]*UserFtInfo      // key: txID-voutIndex(locationID, change when be transfered)
	globalDistributedFtMap map[string]*DistributedFtInfo // operation: dft, key: tickerName
	globalDirectFtMap      map[string]bool               // operation: ft, key: tickerName
}

func (m *Memory) CurrentHeitht() (int64, error) {
	return m.height, nil
}

func (m *Memory) CurrentLocation() (int64, string, error) {
	return m.height, m.txID, nil
}

func (m *Memory) UpdateLocation(height int64, txID string) error {
	m.height = height
	m.txID = txID
	return nil
}

func (m *Memory) NftUTXOsByAtomicalsID(atomicalsID string) ([]*UserNftInfo, error) {
	nfts, ok := m.nftUTXOsByAtomicalsID[atomicalsID]
	if !ok {
		return nil, nil
	}
	return nfts, nil
}

func (m *Memory) InsertNftUTXOByAtomicalsID(UTXO *UserNftInfo) error {
	if _, ok := m.nftUTXOsByAtomicalsID[UTXO.AtomicalsID]; !ok {
		m.nftUTXOsByAtomicalsID[UTXO.AtomicalsID] = make([]*UserNftInfo, 0)
	}
	m.nftUTXOsByAtomicalsID[UTXO.AtomicalsID] = append(m.nftUTXOsByAtomicalsID[UTXO.AtomicalsID], UTXO)
	return nil
}

func (m *Memory) NftUTXOsByLocationID(locationID string) ([]*UserNftInfo, error) {
	nfts, ok := m.nftUTXOsByLocationID[locationID]
	if !ok {
		return nil, nil
	}
	return nfts, nil
}

func (m *Memory) InsertNftUTXOByLocationID(UTXO *UserNftInfo) error {
	if _, ok := m.nftUTXOsByLocationID[UTXO.LocationID]; !ok {
		m.nftUTXOsByLocationID[UTXO.LocationID] = make([]*UserNftInfo, 0)
	}
	m.nftUTXOsByLocationID[UTXO.LocationID] = append(m.nftUTXOsByLocationID[UTXO.LocationID], UTXO)
	return nil
}
func (m *Memory) TransferNftUTXO(oldLocationID, newLocationID, newUserPk string) error {
	return nil
}

func (m *Memory) DeleteNftUTXOByLocationID(locationID string) error {
	m.nftUTXOsByLocationID[locationID] = nil
	return nil
}

func (m *Memory) ParentRealmHasExist(parentRealmAtomicalsID string) (string, error) {
	if _, ok := m.nftUTXOsByAtomicalsID[parentRealmAtomicalsID]; !ok {
		return "", nil
	}
	// TODO: need to ensure if nft will be merged when they're transferred
	if len(m.nftUTXOsByAtomicalsID[parentRealmAtomicalsID]) == 0 {
		return "", nil
	}
	parentRealmName := m.nftUTXOsByAtomicalsID[parentRealmAtomicalsID][0].RealmName
	if _, ok := m.globalNftRealmMap[parentRealmName]; !ok {
		panic("globalNftRealmMap and nftUTXOsByAtomicalsID are not match")
	}
	return parentRealmName, nil
}

func (m *Memory) ParentContainerHasExist(parentContainerAtomicalsID string) (string, error) {
	if _, ok := m.nftUTXOsByAtomicalsID[parentContainerAtomicalsID]; !ok {
		return "", nil
	}
	// TODO: need to ensure if nft will be merged when they're transferred
	if len(m.nftUTXOsByAtomicalsID[parentContainerAtomicalsID]) == 0 {
		return "", nil
	}
	parentContainerName := m.nftUTXOsByAtomicalsID[parentContainerAtomicalsID][0].ContainerName
	if _, ok := m.globalNftContainerMap[parentContainerName]; !ok {
		panic("globalNftContainerMap and nftUTXOsByAtomicalsID are not match")
	}
	return parentContainerName, nil
}

func (m *Memory) NftRealmByNameHasExist(realmName string) (bool, error) {
	_, ok := m.globalNftRealmMap[realmName]
	if !ok {
		return false, nil
	}
	return true, nil
}

func (m *Memory) NftSubRealmByName(realmName, subRealm string) (bool, error) {
	_, ok := m.globalNftRealmMap[realmName]
	if !ok {
		return false, nil
	}
	_, ok = m.globalNftRealmMap[realmName][subRealm]
	if !ok {
		return false, nil
	}
	return true, nil
}

// func (m *Memory) InsertRealm(realmName string) error {
// 	m.globalNftRealmMap[realmName] = make(map[string]bool, 0)
// 	return nil
// }

// func (m *Memory) InsertSubRealm(realmName, subRealm string) error {
// 	m.globalNftRealmMap[realmName][subRealm] = true
// 	return nil
// }

func (m *Memory) NftContainerByName(containerName string) (bool, error) {
	_, ok := m.globalNftContainerMap[containerName]
	if !ok {
		return false, nil
	}
	return true, nil
}

// func (m *Memory) InsertContainer(containerName string) error {
// 	m.globalNftContainerMap[containerName] = make(map[string]bool, 0)
// 	return nil
// }

// func (m *Memory) InsertItemInContainer(containerName, itemID string) error {
// 	m.globalNftContainerMap[containerName][itemID] = true
// 	return nil
// }

func (m *Memory) FtUTXOsByLocationID(locationID string) ([]*UserFtInfo, error) {
	fts, ok := m.ftUTXOs[locationID]
	if !ok {
		return nil, nil
	}
	return fts, nil
}

func (m *Memory) InsertFtUTXO(UTXO *UserFtInfo) error {
	if _, ok := m.ftUTXOs[UTXO.LocationID]; !ok {
		m.ftUTXOs[UTXO.LocationID] = make([]*UserFtInfo, 0)
	}
	m.ftUTXOs[UTXO.LocationID] = append(m.ftUTXOs[UTXO.LocationID], UTXO)
	return nil
}

func (m *Memory) DeleteFtUTXO(locationID string) error {
	m.ftUTXOs[locationID] = nil
	return nil
}

func (m *Memory) DistributedFtByName(tickerName string) (*DistributedFtInfo, error) {
	ft, ok := m.globalDistributedFtMap[tickerName]
	if !ok {
		return nil, nil
	}
	return ft, nil
}

func (m *Memory) InsertDistributedFt(ft *DistributedFtInfo) error {
	m.globalDistributedFtMap[ft.TickerName] = ft
	return nil
}

func (m *Memory) UpdateDistributedFtAmount(tickerName string, mintTimes int64) error {
	m.globalDistributedFtMap[tickerName].MintedTimes = mintTimes
	return nil
}

func (m *Memory) DirectFtByName(tickerName string) (bool, error) {
	_, ok := m.globalDirectFtMap[tickerName]
	if !ok {
		return false, nil
	}
	return true, nil
}

func (m *Memory) InsertDirectFt(tickerName string) error {
	m.globalDirectFtMap[tickerName] = true
	return nil
}

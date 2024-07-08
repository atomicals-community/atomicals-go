package repo

import (
	"github.com/atomicals-go/repo/postsql"
)

func (m *Postgres) fetchUTXO() {
	for locationID := range m.locationIDChannel {
		ftUTXO, err := m.FtUTXOsByLocationID(locationID)
		if err != nil {
			m.ftUTXOMap.Store(locationID, err)
		} else {
			m.ftUTXOMap.Store(locationID, ftUTXO)
		}
		nftUTXO, err := m.NftUTXOsByLocationID(locationID)
		if err != nil {
			m.nftUTXOMap.Store(locationID, err)
		} else {
			m.nftUTXOMap.Store(locationID, nftUTXO)
		}
	}
}

func (m *Postgres) AddLocationIDIntoChannel(locationID string) {
	m.locationIDChannel <- locationID
}

func (m *Postgres) GetNftUTXOFromChannel(locationID string) ([]*postsql.UTXONftInfo, error) {
	for {
		res, ok := m.nftUTXOMap.LoadAndDelete(locationID)
		if ok {
			nftUTXO, ok := res.([]*postsql.UTXONftInfo)
			if !ok {
				err, ok := res.(error)
				if !ok {
					panic("invalid type")
				}
				return nil, err
			}
			return nftUTXO, nil
		}
	}
}

func (m *Postgres) GetFtUTXOFromChannel(locationID string) ([]*postsql.UTXOFtInfo, error) {
	for {
		res, ok := m.ftUTXOMap.LoadAndDelete(locationID)
		if ok {
			ftUTXO, ok := res.([]*postsql.UTXOFtInfo)
			if !ok {
				err, ok := res.(error)
				if !ok {
					panic("invalid type")
				}
				return nil, err
			}
			return ftUTXO, nil
		}
	}
}

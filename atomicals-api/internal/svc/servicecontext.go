package svc

import (
	"sync"

	"github.com/atomicals-go/atomicals-api/internal/config"
	atomicals "github.com/atomicals-go/atomicals-indexer/atomicals-core/operation"
	"github.com/atomicals-go/pkg/conf"
)

type ServiceContext struct {
	Config config.Config
	*atomicals.Atomicals
	SyncHeight               int64
	MaxBlockHeight           int64
	PendingAtomicalsAssetMap map[string]*PendingAtomicalsAsset // key: txID
}

func NewServiceContext(c config.Config, atomicalsConfigFilePath string) *ServiceContext {
	conf, err := conf.ReadJSONFromJSFile(atomicalsConfigFilePath)
	if err != nil {
		panic(err)
	}
	a := atomicals.NewAtomicalsWithSQL(conf)
	svc := &ServiceContext{
		Config:                   c,
		Atomicals:                a,
		PendingAtomicalsAssetMap: make(map[string]*PendingAtomicalsAsset, 0),
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		svc.SyncPendingAtomicalsAsset()
	}()
	wg.Wait()
	return svc
}

package svc

import (
	"sync"

	"github.com/atomicals-go/atomicals-api/internal/config"
	atomicals "github.com/atomicals-go/atomicals-core/operation"
	"github.com/atomicals-go/pkg/conf"
)

type ServiceContext struct {
	Config config.Config
	*atomicals.Atomicals
	SyncHeight            int64
	MaxBlockHeight        int64
	PendingAtomicalsAsset *PendingAtomicalsAsset
}

func NewServiceContext(c config.Config, atomicalsConfigFilePath string) *ServiceContext {
	conf, err := conf.ReadJSONFromJSFile(atomicalsConfigFilePath)
	if err != nil {
		panic(err)
	}
	a := atomicals.NewAtomicalsWithSQL(conf)
	svc := &ServiceContext{
		Config:                c,
		Atomicals:             a,
		PendingAtomicalsAsset: &PendingAtomicalsAsset{},
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		svc.SyncPendingAtomicalsAsset()
	}()
	wg.Wait()
	return svc
}

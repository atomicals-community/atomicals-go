package svc

import (
	"github.com/atomicals-go/atomicals-api/internal/config"
	atomicals "github.com/atomicals-go/atomicals-indexer/atomicals-core/operation"
	"github.com/atomicals-go/pkg/conf"
	"github.com/atomicals-go/repo"
)

type ServiceContext struct {
	Config config.Config
	*atomicals.Atomicals
	CurrentHeight            int64
	MaxBlockHeight           int64
	PendingAtomicalsAssetMap map[string]*repo.AtomicaslData // key: txID
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
		PendingAtomicalsAssetMap: make(map[string]*repo.AtomicaslData, 0),
	}
	return svc
}

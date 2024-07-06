package svc

import (
	"github.com/atomicals-go/atomicals-api/internal/config"
	"github.com/atomicals-go/repo"
)

type ServiceContext struct {
	Config config.Config
	repo.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		DB:     repo.NewSqlDB(c.SqlDNS),
	}
}

package svc

import (
	"github.com/tango/explore/internal/config"
	"github.com/tango/explore/internal/storage"
)

type ServiceContext struct {
	Config  config.Config
	Storage *storage.MemoryStorage
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		Storage: storage.NewMemoryStorage(),
	}
}

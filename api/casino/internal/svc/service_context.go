package svc

import (
	"casinoDemo/api/casino/internal/config"
	"casinoDemo/api/casino/svc/casino_svc"

	"gorm.io/gorm"
)

type ServiceContext struct {
	Config    config.Config
	CasinoSvc *casino_svc.CasinoSvc
	CasinoDb  *gorm.DB
}

func NewServiceContext(c config.Config, casinoSvc *casino_svc.CasinoSvc, casinoDb *gorm.DB) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		CasinoSvc: casinoSvc,
		CasinoDb:  casinoDb,
	}
}

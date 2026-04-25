package user

import (
	"go-api-starterkit/internal/modules/audit"
	tokenModule "go-api-starterkit/internal/modules/token"

	"gorm.io/gorm"
)

type Module struct {
	Repository Repository
	Service    *Service
	Handler    *Handler
}

func BuildModule(db *gorm.DB, auditSvc *audit.Service) *Module {
	repository := NewRepository(db)
	refreshRepo := tokenModule.NewRefreshTokenRepository(db)
	service := NewService(repository, refreshRepo, auditSvc)
	handler := NewHandler(service)

	return &Module{
		Repository: repository,
		Service:    service,
		Handler:    handler,
	}
}

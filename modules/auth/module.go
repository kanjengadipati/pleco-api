package auth

import (
	"go-auth-app/config"
	userModule "go-auth-app/modules/user"
	"go-auth-app/services"

	"gorm.io/gorm"
)

type Module struct {
	Service AuthService
	Handler *AuthHandler
}

func BuildModule(db *gorm.DB, cfg config.AppConfig, userService *userModule.Service, jwtService *services.JWTService) *Module {
	service := NewService(db, cfg, userService, jwtService)
	handler := NewHandler(service)

	return &Module{
		Service: service,
		Handler: handler,
	}
}

package routes

import (
	"go-auth-app/appsetup"
	"go-auth-app/config"

	"github.com/gin-gonic/gin"
	"go-auth-app/services"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB, cfg config.AppConfig, jwtService *services.JWTService) {
	appsetup.RegisterRoutes(router, db, cfg, jwtService)
}

package appsetup

import (
	"go-auth-app/config"
	"go-auth-app/modules/auth"
	"go-auth-app/modules/user"
	"go-auth-app/services"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB, cfg config.AppConfig, jwtService *services.JWTService) {
	api := router.Group("/")
	userModule := user.BuildModule(db)
	authModule := auth.BuildModule(db, cfg, userModule.Service, jwtService)

	auth.SetupRoutes(api, authModule.Handler, jwtService)
	user.SetupRoutes(api, userModule.Handler, jwtService)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}

func BuildRouter(db *gorm.DB, cfg config.AppConfig, jwtService *services.JWTService) *gin.Engine {
	router := gin.Default()
	RegisterRoutes(router, db, cfg, jwtService)
	return router
}

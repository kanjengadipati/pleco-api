package main

import (
	"go-auth-app/appsetup"
	"go-auth-app/config"

	"go-auth-app/services"
)

func initApp() {
	// Load env & init JWT
	config.LoadEnv()
}

func main() {
	initApp()
	appConfig := config.LoadAppConfig()
	db := config.ConnectDB(appConfig.DatabaseURL)
	appsetup.RunStartupTasks(appConfig, db)
	jwtService := services.NewJWTService(appConfig.JWTSecret)
	router := appsetup.BuildRouter(db, appConfig, jwtService)
	registerDocsRoutes(router)

	if err := router.Run(":" + appConfig.Port); err != nil {
		panic(err)
	}
}

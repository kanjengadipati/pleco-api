package main

import (
	"log"

	"go-auth-app/appsetup"
	"go-auth-app/config"
)

func main() {
	config.LoadEnv()
	appConfig := config.LoadAppConfig()
	if err := appsetup.RunMigrations(appConfig.DatabaseURL); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migration success")
}

package main

import (
	"log"

	"go-auth-app/appsetup"
	"go-auth-app/config"
)

func main() {
	config.LoadEnv()
	if err := appsetup.RunMigrations(); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migration success")
}

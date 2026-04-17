package main

import (
	"log"

	"go-auth-app/appsetup"
	"go-auth-app/config"
)

func main() {
	// Load env (WAJIB)
	config.LoadEnv()

	// Init DB (WAJIB)
	config.ConnectDB()
	log.Println("Start seeding...")
	appsetup.RunSeeds()

	log.Println("Seeding done 🚀")
}

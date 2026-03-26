package main

import (
	"go-auth-app/config"
	"go-auth-app/routes"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file loaded: %v", err)
	}

	// Connect DB
	config.ConnectDB()

	// Seed admin
	config.SeedAdmin()

	// Init router
	router := gin.Default()

	// Setup routes
	routes.SetupRoutes(router)

	// Run server
	router.Run(":8080")
}
package appsetup

import (
	"log"
	"strings"

	"go-auth-app/config"
	"go-auth-app/seeds"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunStartupTasks() {
	if shouldRunFromEnv("AUTO_RUN_MIGRATIONS") {
		if err := RunMigrations(); err != nil {
			log.Fatalf("❌ startup migrations failed: %v", err)
		}
		log.Println("✅ startup migrations completed")
	}

	if shouldRunFromEnv("AUTO_RUN_SEEDS") {
		RunSeeds()
		log.Println("✅ startup seeds completed")
	}
}

func RunMigrations() error {
	dbURL := config.DatabaseURL()
	if dbURL == "" {
		log.Fatal("❌ DATABASE_URL is not set")
	}

	m, err := migrate.New("file://migrations", dbURL)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func RunSeeds() {
	if config.DB == nil {
		log.Fatal("❌ DB is not initialized before seeding")
	}

	seeds.SeedRoles(config.DB)
	seeds.SeedPermissions(config.DB)
	seeds.SeedAdmin(config.DB)
}

func shouldRunFromEnv(key string) bool {
	value := strings.TrimSpace(strings.ToLower(config.GetEnv(key, "")))
	return value == "1" || value == "true" || value == "yes"
}

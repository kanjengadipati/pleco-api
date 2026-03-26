package config

import (
	"go-auth-app/models"
	"golang.org/x/crypto/bcrypt"
)

func SeedAdmin() {
	var user models.User

	DB.Where("email = ?", "admin@mail.com").First(&user)

	if user.ID == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), 14)

		admin := models.User{
			Name:     "Super Admin",
			Email:    "admin@mail.com",
			Password: string(hashedPassword),
			Role:     "admin",
		}

		DB.Create(&admin)
	}
}
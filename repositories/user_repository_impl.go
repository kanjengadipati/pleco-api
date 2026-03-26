package repositories

import (
	"go-auth-app/config"
	"go-auth-app/models"
)

type UserRepoDB struct{}

func (r *UserRepoDB) Create(user *models.User) error {
	return config.DB.Create(user).Error
}

func (r *UserRepoDB) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := config.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *UserRepoDB) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := config.DB.First(&user, id).Error
	return &user, err
}

func (r *UserRepoDB) Update(user *models.User) error {
	return config.DB.Save(user).Error
}
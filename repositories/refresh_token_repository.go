package repositories

import "go-auth-app/models"

// RefreshTokenRepository defines the methods required for working with refresh tokens in the data store.
type RefreshTokenRepository interface {
	Save(token *models.RefreshToken) error
	FindByUserAndDevice(userID uint, deviceID string) (*models.RefreshToken, error)
	FindByUser(userID uint) ([]models.RefreshToken, error)
	DeleteByID(id uint) error
	DeleteByUser(userID uint) error
}

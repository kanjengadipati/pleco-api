package repositories

import (
	"go-auth-app/config"
	"go-auth-app/models"
)

// RefreshTokenRepoDB implements the RefreshTokenRepository interface using a GORM DB instance.
type RefreshTokenRepoDB struct{}

// NewRefreshTokenRepo returns a pointer to a struct that satisfies the RefreshTokenRepository interface.
// We're using the global DB from config in this pattern (to match the style in user_repository_impl.go).
func NewRefreshTokenRepo() *RefreshTokenRepoDB {
	return &RefreshTokenRepoDB{}
}

func (r *RefreshTokenRepoDB) Save(token *models.RefreshToken) error {
	return config.DB.Create(token).Error
}

// FindByUserAndDevice expects userID to be uint, deviceID to be string to match the Repository interface.
func (r *RefreshTokenRepoDB) FindByUserAndDevice(userID uint, deviceID string) (*models.RefreshToken, error) {
	var token models.RefreshToken
	err := config.DB.Where("user_id = ? AND device_id = ?", userID, deviceID).First(&token).Error
	return &token, err
}

func (r *RefreshTokenRepoDB) FindByUser(userID uint) ([]models.RefreshToken, error) {
	var tokens []models.RefreshToken
	err := config.DB.Where("user_id = ?", userID).Find(&tokens).Error
	return tokens, err
}

// DeleteByID expects id to be uint to match the Repository interface.
func (r *RefreshTokenRepoDB) DeleteByID(id uint) error {
	return config.DB.Delete(&models.RefreshToken{}, id).Error
}

// DeleteByUser expects userID to be uint to match the Repository interface.
func (r *RefreshTokenRepoDB) DeleteByUser(userID uint) error {
	return config.DB.Where("user_id = ?", userID).Delete(&models.RefreshToken{}).Error
}

package repositories

import (
	"go-auth-app/models"
)

type SocialAccountRepository interface {
	Create(socialAccount *models.SocialAccount) error
	FindByProvider(provider string, providerID string) (*models.SocialAccount, error)
}

package models

type SocialAccount struct {
	ID         uint
	UserID     uint
	Provider   string // google, github
	ProviderID string
	AvatarURL  string
}

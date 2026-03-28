package models

import (
	"time"

	"gorm.io/gorm"
)

type RefreshToken struct {
	gorm.Model
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	TokenHash string
	DeviceID  string
	UserAgent string
	IPAddress string
	ExpiredAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

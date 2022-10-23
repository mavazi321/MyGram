package models

import "time"

type SocialMedia struct {
	ID             uint   `gorm:"primaryKey"`
	Name           string `gorm:"not null"`
	SocialMediaUrl string `gorm:"not null" json:"social_media_url"`
	UserID         uint
	User           User
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

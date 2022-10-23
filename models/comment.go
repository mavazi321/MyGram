package models

import "time"

type Comment struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	User      User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	PhotoID   uint   `json:"photo_id"`
	Photo     Photo  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Message   string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

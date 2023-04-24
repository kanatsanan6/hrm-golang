package model

import "time"

type User struct {
	ID                uint `gorm:"primaryKey"`
	FirstName         string
	LastName          string
	Email             string `gorm:"uniqueIndex"`
	EncryptedPassword string
	CreatedAt         time.Time `gorm:"autoCreateTime"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime"`
}

package model

import "time"

type User struct {
	ID                uint `gorm:"primaryKey"`
	FirstName         string
	LastName          string
	Email             string `gorm:"uniqueIndex"`
	EncryptedPassword string
	CompanyID         *uint
	Company           Company   `gorm:"foreignKey:CompanyID"`
	CreatedAt         time.Time `gorm:"autoCreateTime"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime"`
}

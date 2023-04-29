package model

import "time"

type Company struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	Users     []User
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

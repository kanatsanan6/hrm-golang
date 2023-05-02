package model

import "time"

type Leave struct {
	ID          uint `gorm:"primaryKey"`
	Description string
	Status      string
	StartDate   time.Time
	EndDate     time.Time
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

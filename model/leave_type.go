package model

import "time"

type LeaveType struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	Usage     int
	Max       int
	UserID    uint
	User      User `gorm:"foreignKey:UserID"`
	Leaves    []Leave
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

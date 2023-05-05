package model

import "time"

type Leave struct {
	ID          uint `gorm:"primaryKey"`
	Description string
	Status      string
	StartDate   time.Time
	EndDate     time.Time
	UserID      uint
	User        User `gorm:"foreignKey:UserID"`
	LeaveTypeID uint
	LeaveType   LeaveType `gorm:"foreignKey:LeaveTypeID"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

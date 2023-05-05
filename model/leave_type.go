package model

import "time"

var DefaultLeaveType = []map[string]interface{}{
	{"name": "vacation_leave", "max": 10},
	{"name": "extra_vacation", "max": 5},
	{"name": "sick_leave", "max": 10},
	{"name": "business_leave", "max": 3},
}

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

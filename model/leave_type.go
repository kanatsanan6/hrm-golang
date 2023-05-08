package model

import "time"

var DefaultLeaveType = []map[string]interface{}{
	{"name": "vacation_leave", "max": 10},
	{"name": "extra_vacation", "max": 5},
	{"name": "sick_leave", "max": 10},
	{"name": "business_leave", "max": 3},
}

type LeaveType struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Usage     int       `json:"usages"`
	Max       int       `json:"max"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

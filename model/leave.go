package model

import "time"

type LeaveStruct struct {
	ID          int64     `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	UserID      int64     `json:"user_id"`
	User        User      `json:"user"`
	LeaveTypeID int64     `json:"leave_type_id"`
	LeaveType   LeaveType `json:"leave_type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func LeaveResponse(leave Leave) LeaveStruct {
	return LeaveStruct{
		ID:          leave.ID,
		Description: leave.Description,
		Status:      leave.Status,
		StartDate:   leave.StartDate,
		EndDate:     leave.EndDate,
		User:        User{},
		LeaveType:   LeaveType{},
		CreatedAt:   leave.CreatedAt,
		UpdatedAt:   leave.UpdatedAt,
	}
}

type Leave struct {
	ID          int64     `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	UserID      int64     `json:"user_id"`
	LeaveTypeID int64     `json:"leave_type_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

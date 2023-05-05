package queries

import (
	"time"

	"github.com/kanatsanan6/hrm/model"
)

type LeaveStruct struct {
	ID          uint      `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	User        *UserType `json:"user_id"`
	LeaveType   LeaveType `json:"leave_type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateLeaveArgs struct {
	Description string
	Status      string
	StartDate   time.Time
	EndDate     time.Time
	UserID      uint
	LeaveTypeID uint
}

func (q *SQLQueries) CreateLeave(args CreateLeaveArgs) (model.Leave, error) {
	leave := model.Leave{
		Description: args.Description,
		Status:      args.Status,
		StartDate:   args.StartDate,
		EndDate:     args.EndDate,
		LeaveTypeID: args.LeaveTypeID,
		UserID:      args.UserID,
	}

	if err := q.DB.Create(&leave).Error; err != nil {
		return model.Leave{}, err
	}
	return leave, nil
}

func (q *SQLQueries) GetLeaves(user *model.User) ([]LeaveStruct, error) {
	var leaves []model.Leave
	err := q.DB.Where("user_id = ?", user.ID).Preload("LeaveType").Find(&leaves).Error
	if err != nil {
		return []LeaveStruct{}, err
	}

	var result []LeaveStruct
	for _, leave := range leaves {
		result = append(result, LeaveStruct{
			ID:          leave.ID,
			Description: leave.Description,
			Status:      leave.Status,
			StartDate:   leave.StartDate,
			EndDate:     leave.EndDate,
			CreatedAt:   leave.CreatedAt,
			UpdatedAt:   leave.UpdatedAt,
			LeaveType: LeaveType{
				ID:        leave.LeaveType.ID,
				Name:      leave.LeaveType.Name,
				Usage:     leave.LeaveType.Usage,
				Max:       leave.LeaveType.Max,
				CreatedAt: leave.LeaveType.CreatedAt,
				UpdatedAt: leave.LeaveType.UpdatedAt,
			},
		})
	}
	return result, nil
}

package queries

import (
	"time"

	"github.com/kanatsanan6/hrm/model"
)

type LeaveStruct struct {
	ID          uint            `json:"id"`
	Description string          `json:"description"`
	Status      string          `json:"status"`
	StartDate   time.Time       `json:"start_date"`
	EndDate     time.Time       `json:"end_date"`
	LeaveType   model.LeaveType `json:"leave_type"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
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

func (q *SQLQueries) GetLeaves(user *model.User) []LeaveStruct {
	leaves := []LeaveStruct{}
	for _, leave := range user.Leaves {
		leaves = append(leaves, LeaveStruct{
			ID:          leave.ID,
			Description: leave.Description,
			LeaveType:   leave.LeaveType,
			Status:      leave.Status,
			StartDate:   leave.StartDate,
			EndDate:     leave.EndDate,
			CreatedAt:   leave.CreatedAt,
			UpdatedAt:   leave.UpdatedAt,
		})
	}
	return leaves
}

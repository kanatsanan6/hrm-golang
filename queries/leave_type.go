package queries

import (
	"time"

	"github.com/kanatsanan6/hrm/model"
)

type LeaveType struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Usage     int       `json:"usage"`
	Max       int       `json:"max"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateLeaveTypeArgs struct {
	Name   string
	Usage  int
	Max    int
	UserID uint
}

func (q *SQLQueries) CreateLeaveType(args CreateLeaveTypeArgs) (model.LeaveType, error) {
	leaveType := model.LeaveType{
		Name:   args.Name,
		Usage:  args.Usage,
		Max:    args.Max,
		UserID: args.UserID,
	}

	if err := q.DB.Create(&leaveType).Error; err != nil {
		return model.LeaveType{}, err
	}
	return leaveType, nil
}

func (q *SQLQueries) FindUserLeaveTypeByName(user model.User, name string) (model.LeaveType, error) {
	var leaveType model.LeaveType
	if err := q.DB.Where("user_id = ? AND name = ?", user.ID, name).First(&leaveType).Error; err != nil {
		return model.LeaveType{}, err
	}
	return leaveType, nil
}

func (q *SQLQueries) GetUserLeaveTypes(user model.User) []LeaveType {
	var result []LeaveType
	for _, leaveType := range user.LeaveTypes {
		result = append(result, LeaveType{
			ID:        leaveType.ID,
			Name:      leaveType.Name,
			Usage:     leaveType.Max - leaveType.Usage,
			Max:       leaveType.Max,
			CreatedAt: leaveType.CreatedAt,
			UpdatedAt: leaveType.UpdatedAt,
		})
	}
	return result
}

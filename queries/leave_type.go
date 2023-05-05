package queries

import "github.com/kanatsanan6/hrm/model"

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

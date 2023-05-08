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
	UserID int64
}

func (q *SQLQueries) CreateLeaveType(args CreateLeaveTypeArgs) (model.LeaveType, error) {
	var leaveType model.LeaveType
	query := `
	INSERT INTO leave_types (name, usage, max, user_id)
	VALUES ($1, $2, $3, $4)
	RETURNING id, name, usage, max, user_id, created_at, updated_at
	`
	row := q.DB.QueryRow(query, args.Name, args.Usage, args.Max, args.UserID)
	err := row.Scan(
		&leaveType.ID,
		&leaveType.Name,
		&leaveType.Usage,
		&leaveType.Max,
		&leaveType.UserID,
		&leaveType.CreatedAt,
		&leaveType.UpdatedAt,
	)
	return leaveType, err
}

func (q *SQLQueries) FindUserLeaveTypeByName(user model.User, name string) (model.LeaveType, error) {
	var leaveType model.LeaveType
	row := q.DB.QueryRow(`SELECT * FROM leave_types WHERE (user_id) = ($1) AND (name) = ($2) LIMIT 1`, user.ID, name)
	err := row.Scan(
		&leaveType.ID,
		&leaveType.Name,
		&leaveType.Usage,
		&leaveType.Max,
		&leaveType.UserID,
		&leaveType.CreatedAt,
		&leaveType.UpdatedAt,
	)
	return leaveType, err
}

func (q *SQLQueries) GetUserLeaveTypes(user model.User) ([]model.LeaveType, error) {
	var leaveTypes []model.LeaveType
	row, err := q.DB.Query(`SELECT * FROM leave_types WHERE (user_id) = ($1)`, user.ID)
	if err != nil {
		return []model.LeaveType{}, err
	}

	for row.Next() {
		var leaveType model.LeaveType
		err := row.Scan(
			&leaveType.ID,
			&leaveType.Name,
			&leaveType.Usage,
			&leaveType.Max,
			&leaveType.UserID,
			&leaveType.CreatedAt,
			&leaveType.UpdatedAt,
		)
		if err != nil {
			return []model.LeaveType{}, err
		}
		leaveTypes = append(leaveTypes, leaveType)
	}

	return leaveTypes, nil
}

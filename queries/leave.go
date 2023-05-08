package queries

import (
	"time"

	"github.com/kanatsanan6/hrm/model"
)

type LeaveStruct struct {
	ID          int64     `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	User        UserType  `json:"user_id"`
	LeaveType   LeaveType `json:"leave_type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func LeaveResponse(leave model.Leave) *LeaveStruct {
	// return &LeaveStruct{
	// 	ID:          leave.ID,
	// 	Description: leave.Description,
	// 	Status:      leave.Status,
	// 	StartDate:   leave.StartDate,
	// 	EndDate:     leave.EndDate,
	// 	CreatedAt:   leave.CreatedAt,
	// 	UpdatedAt:   leave.UpdatedAt,
	// 	LeaveType: LeaveType{
	// 		// ID:        leave.LeaveType.ID,
	// 		Name:      leave.LeaveType.Name,
	// 		Usage:     leave.LeaveType.Usage,
	// 		Max:       leave.LeaveType.Max,
	// 		CreatedAt: leave.LeaveType.CreatedAt,
	// 		UpdatedAt: leave.LeaveType.UpdatedAt,
	// 	},
	// }
	return &LeaveStruct{}
}

type CreateLeaveArgs struct {
	Description string
	Status      string
	StartDate   time.Time
	EndDate     time.Time
	UserID      int64
	LeaveTypeID int64
}

func (q *SQLQueries) CreateLeave(args CreateLeaveArgs) (model.Leave, error) {
	var leave model.Leave
	query := `
	INSERT INTO leaves (
		description,
		status,
		start_date,
		end_date,
		leave_type_id,
		user_id
	) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *
	`
	row := q.DB.QueryRow(
		query,
		args.Description,
		args.Status,
		args.StartDate,
		args.EndDate,
		args.LeaveTypeID,
		args.UserID,
	)
	err := row.Scan(
		&leave.ID,
		&leave.Description,
		&leave.Status,
		&leave.StartDate,
		&leave.EndDate,
		&leave.CreatedAt,
		&leave.UpdatedAt,
		&leave.UserID,
		&leave.LeaveTypeID,
	)
	return leave, err
}

func (q *SQLQueries) GetLeaves(user *model.User) ([]model.Leave, error) {
	var leaves []model.Leave
	row, err := q.DB.Query(`SELECT * FROM leaves WHERE (user_id) = ($1)`, user.ID)
	if err != nil {
		return []model.Leave{}, err
	}
	for row.Next() {
		var leave model.Leave
		err := row.Scan(
			&leave.ID,
			&leave.Description,
			&leave.Status,
			&leave.StartDate,
			&leave.EndDate,
			&leave.CreatedAt,
			&leave.UpdatedAt,
			&leave.UserID,
			&leave.LeaveTypeID,
		)
		if err != nil {
			return []model.Leave{}, err
		}
		leaves = append(leaves, leave)
	}
	return leaves, nil
}

func (q *SQLQueries) GetLeaveByID(id int64) (model.Leave, error) {
	var leave model.Leave
	row := q.DB.QueryRow(`SELECT * FROM leaves WHERE (id) = ($1)`, id)
	err := row.Scan(
		&leave.ID,
		&leave.Description,
		&leave.Status,
		&leave.StartDate,
		&leave.EndDate,
		&leave.CreatedAt,
		&leave.UpdatedAt,
		&leave.UserID,
		&leave.LeaveTypeID,
	)
	return leave, err
}

type UpdateLeaveArgs struct {
	ID     int64
	Status string
}

func (q *SQLQueries) UpdateLeave(args UpdateLeaveArgs) (model.Leave, error) {
	var leave model.Leave
	query := `
	UPDATE leaves
	SET
	status = COALESCE($2, status)
	WHERE id = $1 RETURNING *
	`
	row := q.DB.QueryRow(query, args.ID, args.Status)
	err := row.Scan(
		&leave.ID,
		&leave.Description,
		&leave.Status,
		&leave.StartDate,
		&leave.EndDate,
		&leave.CreatedAt,
		&leave.UpdatedAt,
		&leave.UserID,
		&leave.LeaveTypeID,
	)
	return leave, err
}

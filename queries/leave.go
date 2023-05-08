package queries

import (
	"time"

	"github.com/kanatsanan6/hrm/model"
)

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

func (q *SQLQueries) GetLeaves(user *model.User) ([]model.LeaveStruct, error) {
	var leaves []model.LeaveStruct
	query := `
	SELECT l.*,
	u.* AS user,
	lt.* AS leave_type
	FROM leaves l
	JOIN leave_types lt ON l.leave_type_id = lt.ID
	JOIN users u ON l.user_id = u.ID
	WHERE l.user_id = ($1)
	`
	row, err := q.DB.Query(query, user.ID)
	if err != nil {
		return []model.LeaveStruct{}, err
	}
	for row.Next() {
		var leave model.LeaveStruct
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
			&leave.User.ID,
			&leave.User.FirstName,
			&leave.User.LastName,
			&leave.User.Email,
			&leave.User.EncryptedPassword,
			&leave.User.CreatedAt,
			&leave.User.UpdatedAt,
			&leave.User.CompanyID,
			&leave.User.ResetPasswordToken,
			&leave.User.Role,
			&leave.LeaveType.ID,
			&leave.LeaveType.Name,
			&leave.LeaveType.Usage,
			&leave.LeaveType.Max,
			&leave.LeaveType.UserID,
			&leave.LeaveType.CreatedAt,
			&leave.LeaveType.UpdatedAt,
		)
		if err != nil {
			return []model.LeaveStruct{}, err
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

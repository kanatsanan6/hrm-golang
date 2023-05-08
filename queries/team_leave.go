package queries

import (
	"github.com/kanatsanan6/hrm/model"
)

func (q *SQLQueries) GetTeamLeaves(company *model.Company) ([]LeaveStruct, error) {
	var leaves []LeaveStruct
	query := `
		SELECT
			l.*,
			u.id AS "user.id",
			u.first_name AS "user.first_name",
			u.last_name AS "user.last_name",
			u.email AS "user.email",
			u.company_id AS "user.company_id",
			u.role AS "user.role",
			u.created_at AS "user.created_at",
			u.updated_at AS "user.updated_at",
			lt.id AS "leave_type.id",
			lt.name AS "leave_type.name",
			lt.usage AS "leave_type.usage",
			lt.max AS "leave_type.max",
			lt.created_at AS "leave_type.created_at",
			lt.updated_at AS "leave_type.updated_at"
		FROM leaves l
		INNER JOIN users u ON u.id = l.user_id
		INNER JOIN leave_types lt ON lt.id = l.leave_type_id
		WHERE u.company_id = $1
	`
	row, err := q.DB.Query(query, company.ID)
	for row.Next() {
		var leave LeaveStruct
		row.Scan(
			&leave.ID,
			&leave.Description,
			&leave.Status,
			&leave.StartDate,
			&leave.EndDate,
			&leave.User,
			&leave.LeaveType,
			&leave.CreatedAt,
			&leave.UpdatedAt,
		)
		leaves = append(leaves, leave)
	}
	return leaves, err
}

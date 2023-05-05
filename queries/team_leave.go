package queries

import (
	"github.com/kanatsanan6/hrm/model"
)

func (q *SQLQueries) GetTeamLeaves(company *model.Company) ([]LeaveStruct, error) {
	var companyResult *model.Company
	err := q.DB.Where("ID = ?", company.ID).
		Preload("Users").
		Preload("Users.Leaves").
		Preload("Users.Leaves.LeaveType").
		First(&companyResult).Error
	if err != nil {
		return []LeaveStruct{}, err
	}

	var result []LeaveStruct
	for _, user := range companyResult.Users {
		for _, leave := range user.Leaves {
			result = append(result, LeaveStruct{
				ID:          leave.ID,
				Description: leave.Description,
				Status:      leave.Status,
				StartDate:   leave.StartDate,
				EndDate:     leave.EndDate,
				User: &UserType{
					ID:        user.ID,
					Email:     user.Email,
					FirstName: user.FirstName,
					LastName:  user.LastName,
					Role:      user.Role,
					CreatedAt: user.CreatedAt,
					UpdatedAt: user.UpdatedAt,
				},
				LeaveType: LeaveType{
					ID:        leave.LeaveType.ID,
					Name:      leave.LeaveType.Name,
					Usage:     leave.LeaveType.Usage,
					Max:       leave.LeaveType.Max,
					CreatedAt: leave.LeaveType.CreatedAt,
					UpdatedAt: leave.LeaveType.UpdatedAt,
				},
				CreatedAt: leave.CreatedAt,
				UpdatedAt: leave.UpdatedAt,
			})
		}
	}

	return result, nil
}

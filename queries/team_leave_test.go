package queries_test

import (
	"fmt"
	"testing"

	"github.com/kanatsanan6/hrm/queries"
	"github.com/kanatsanan6/hrm/utils"
	"github.com/stretchr/testify/assert"
)

func TestSQLQueries_GetTeamLeaves(t *testing.T) {
	company, _ := testQueries.CreateCompany(queries.CreateCompanyArgs{
		Name: utils.RandomString(10),
	})
	admin, _ := testQueries.CreateUser(queries.CreateUserArgs{
		Email:             utils.RandomEmail(),
		EncryptedPassword: utils.RandomString(16),
		FirstName:         utils.RandomString(10),
		LastName:          utils.RandomString(10),
		CompanyID:         &company.ID,
		Role:              "admin",
	})
	adminLeaveType, _ := testQueries.CreateLeaveType(queries.CreateLeaveTypeArgs{
		Name:   "vacation_leave",
		Usage:  0,
		Max:    10,
		UserID: admin.ID,
	})
	member, _ := testQueries.CreateUser(queries.CreateUserArgs{
		Email:             utils.RandomEmail(),
		EncryptedPassword: utils.RandomString(16),
		FirstName:         utils.RandomString(10),
		LastName:          utils.RandomString(10),
		CompanyID:         &company.ID,
		Role:              "member",
	})
	memberLeaveType, _ := testQueries.CreateLeaveType(queries.CreateLeaveTypeArgs{
		Name:   "vacation_leave",
		Usage:  0,
		Max:    10,
		UserID: member.ID,
	})

	startDate, _ := utils.StringToDateTime("2022-01-01")
	endDate, _ := utils.StringToDateTime("2022-01-02")
	adminLeave, _ := testQueries.CreateLeave(queries.CreateLeaveArgs{
		Description: utils.RandomString(10),
		Status:      "pending",
		StartDate:   startDate,
		EndDate:     endDate,
		LeaveTypeID: adminLeaveType.ID,
		UserID:      admin.ID,
	})
	memberLeave, _ := testQueries.CreateLeave(queries.CreateLeaveArgs{
		Description: utils.RandomString(10),
		Status:      "pending",
		StartDate:   startDate,
		EndDate:     endDate,
		LeaveTypeID: memberLeaveType.ID,
		UserID:      member.ID,
	})

	result, err := testQueries.GetTeamLeaves(&company)
	fmt.Println(result)
	assert.NoError(t, err)
	assert.Equal(t, adminLeave.ID, result[0].ID)
	assert.Equal(t, memberLeave.ID, result[1].ID)
}

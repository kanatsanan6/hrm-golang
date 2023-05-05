package queries_test

import (
	"testing"
	"time"

	"github.com/kanatsanan6/hrm/queries"
	"github.com/kanatsanan6/hrm/utils"
	"github.com/stretchr/testify/assert"
)

func TestSQLQueries_CreateLeave(t *testing.T) {
	user := GenerateUser()
	leaveType, err := testQueries.CreateLeaveType(queries.CreateLeaveTypeArgs{
		Name:   "vacation_leave",
		Usage:  0,
		Max:    10,
		UserID: user.ID,
	})
	assert.NoError(t, err)

	description := utils.RandomString(10)
	status := "pending"
	startDate := time.Now()
	endDate := time.Now().Add(24 * time.Hour)
	leave, err := testQueries.CreateLeave(queries.CreateLeaveArgs{
		Description: description,
		Status:      status,
		StartDate:   startDate,
		EndDate:     endDate,
		LeaveTypeID: leaveType.ID,
		UserID:      user.ID,
	})

	assert.NoError(t, err)
	assert.Equal(t, description, leave.Description)
	assert.Equal(t, status, leave.Status)
	assert.Equal(t, startDate, leave.StartDate)
	assert.Equal(t, endDate, leave.EndDate)
	assert.Equal(t, leaveType.ID, leave.LeaveTypeID)
	assert.Equal(t, user.ID, leave.UserID)
}

func TestSQLQueries_GetLeaves(t *testing.T) {
	user := GenerateUser()
	leaveType, err := testQueries.CreateLeaveType(queries.CreateLeaveTypeArgs{
		Name:   "vacation_leave",
		Usage:  0,
		Max:    10,
		UserID: user.ID,
	})
	assert.NoError(t, err)

	leave1, err := testQueries.CreateLeave(queries.CreateLeaveArgs{
		Description: utils.RandomString(10),
		Status:      "pending",
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(24 * time.Hour),
		LeaveTypeID: leaveType.ID,
		UserID:      user.ID,
	})
	assert.NoError(t, err)

	leave2, err := testQueries.CreateLeave(queries.CreateLeaveArgs{
		Description: utils.RandomString(10),
		Status:      "pending",
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(24 * time.Hour),
		LeaveTypeID: leaveType.ID,
		UserID:      user.ID,
	})
	assert.NoError(t, err)

	updatedUser, err := testQueries.FindUserByEmail(user.Email)
	assert.NoError(t, err)

	result := testQueries.GetLeaves(&updatedUser)
	assert.Equal(t, result[0].ID, leave1.ID)
	assert.Equal(t, result[1].ID, leave2.ID)

}

package queries_test

import (
	"testing"

	"github.com/kanatsanan6/hrm/model"
	"github.com/kanatsanan6/hrm/queries"
	"github.com/kanatsanan6/hrm/utils"
	"github.com/stretchr/testify/assert"
)

func GenerateLeaveType() model.LeaveType {
	user := GenerateUser()
	leaveType, _ := testQueries.CreateLeaveType(queries.CreateLeaveTypeArgs{
		Name:   utils.RandomString(10),
		Usage:  0,
		Max:    10,
		UserID: user.ID,
	})

	return leaveType
}

func TestSQLQueries_CreateLeaveType(t *testing.T) {
	user := GenerateUser()
	leaveType, err := testQueries.CreateLeaveType(queries.CreateLeaveTypeArgs{
		Name:   "vacation_leave",
		Usage:  0,
		Max:    10,
		UserID: user.ID,
	})
	assert.NoError(t, err)
	assert.Equal(t, "vacation_leave", leaveType.Name)
	assert.Equal(t, 0, leaveType.Usage)
	assert.Equal(t, 10, leaveType.Max)
	assert.Equal(t, user.ID, leaveType.UserID)
}

func TestQueries_FindUserLeaveTypeByName(t *testing.T) {
	user := GenerateUser()
	leaveType, _ := testQueries.CreateLeaveType(queries.CreateLeaveTypeArgs{
		Name:   utils.RandomString(10),
		Usage:  0,
		Max:    10,
		UserID: user.ID,
	})

	result, err := testQueries.FindUserLeaveTypeByName(user, leaveType.Name)
	assert.NoError(t, err)
	assert.Equal(t, result.Name, leaveType.Name)
	assert.Equal(t, result.UserID, leaveType.UserID)
}

func TestQueries_GetUserLeaveTypes(t *testing.T) {
	user := GenerateUser()
	user2 := GenerateUser()
	leaveType1, _ := testQueries.CreateLeaveType(queries.CreateLeaveTypeArgs{
		Name:   utils.RandomString(10),
		Usage:  0,
		Max:    10,
		UserID: user.ID,
	})
	leaveType2, _ := testQueries.CreateLeaveType(queries.CreateLeaveTypeArgs{
		Name:   utils.RandomString(10),
		Usage:  0,
		Max:    10,
		UserID: user.ID,
	})
	testQueries.CreateLeaveType(queries.CreateLeaveTypeArgs{
		Name:   utils.RandomString(10),
		Usage:  0,
		Max:    10,
		UserID: user2.ID,
	})

	leaveTypes, err := testQueries.GetUserLeaveTypes(user)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(leaveTypes))
	assert.Equal(t, leaveType1.ID, leaveTypes[0].ID)
	assert.Equal(t, leaveType2.ID, leaveTypes[1].ID)
}

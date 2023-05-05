package queries_test

import (
	"fmt"
	"testing"

	"github.com/kanatsanan6/hrm/model"
	"github.com/kanatsanan6/hrm/queries"
	"github.com/kanatsanan6/hrm/utils"
	"github.com/stretchr/testify/assert"
)

func GenerateLeaveType() *model.LeaveType {
	user := GenerateUser()
	leaveType, _ := testQueries.CreateLeaveType(queries.CreateLeaveTypeArgs{
		Name:   utils.RandomString(10),
		Usage:  0,
		Max:    10,
		UserID: user.ID,
	})

	return &leaveType
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
	leaveType, err := testQueries.CreateLeaveType(queries.CreateLeaveTypeArgs{
		Name:   "vacation_leave",
		Usage:  0,
		Max:    10,
		UserID: user.ID,
	})
	assert.NoError(t, err)

	result, err := testQueries.FindUserLeaveTypeByName(*user, leaveType.Name)
	fmt.Println(result)
	assert.NoError(t, err)
	assert.Equal(t, result.ID, leaveType.ID)
}

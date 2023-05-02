package queries_test

import (
	"testing"
	"time"

	"github.com/kanatsanan6/hrm/queries"
	"github.com/kanatsanan6/hrm/utils"
	"github.com/stretchr/testify/assert"
)

func TestSQLQueries_CreateLeave(t *testing.T) {
	description := utils.RandomString(10)
	status := "pending"
	startDate := time.Now()
	startPeriod := "full"
	endDate := time.Now().Add(24 * time.Hour)
	endPeriod := "full"

	args := queries.CreateLeaveArgs{
		Description: description,
		Status:      status,
		StartDate:   startDate,
		StartPeriod: startPeriod,
		EndDate:     endDate,
		EndPeriod:   endPeriod,
	}
	leave, err := testQueries.CreateLeave(args)

	assert.NoError(t, err)
	assert.Equal(t, description, leave.Description)
	assert.Equal(t, status, leave.Status)
	assert.Equal(t, startDate, leave.StartDate)
	assert.Equal(t, startPeriod, leave.StartPeriod)
	assert.Equal(t, endDate, leave.EndDate)
	assert.Equal(t, endPeriod, leave.EndPeriod)
}

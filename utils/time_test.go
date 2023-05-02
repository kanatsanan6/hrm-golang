package utils_test

import (
	"testing"
	"time"

	"github.com/kanatsanan6/hrm/utils"
	"github.com/stretchr/testify/assert"
)

func TestStringToDateTime(t *testing.T) {
	dateString := "2023-05-23"
	result, err := utils.StringToDateTime(dateString)
	expected := time.Date(2023, 5, 23, 0, 0, 0, 0, time.UTC)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)

	dateString = "23-05-2023"
	_, err = utils.StringToDateTime(dateString)
	assert.Error(t, err)
}

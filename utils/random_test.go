package utils_test

import (
	"regexp"
	"testing"

	"github.com/kanatsanan6/hrm/utils"
	"github.com/stretchr/testify/assert"
)

func TestRandomNumber(t *testing.T) {
	result := utils.RandomNumber(1, 10)

	assert.NotNil(t, result)
	assert.Greater(t, result, 1)
	assert.Less(t, result, 10)
}

func TestRandomString(t *testing.T) {
	result := utils.RandomString(10)

	assert.NotNil(t, result)
	assert.Equal(t, 10, len(result))
}

func TestRandomEmail(t *testing.T) {
	result := utils.RandomEmail()

	pattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(pattern, result)

	assert.True(t, match)
}

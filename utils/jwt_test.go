package utils_test

import (
	"testing"

	"github.com/kanatsanan6/hrm/utils"
	"github.com/stretchr/testify/assert"
)

func TestGenerateJWT(t *testing.T) {
	email := utils.RandomEmail()

	token, claims, err := utils.GenerateJWT(email)

	assert.NoError(t, err)
	assert.NotNil(t, token)
	assert.Equal(t, claims.Email, email)
}

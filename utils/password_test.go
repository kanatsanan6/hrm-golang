package utils_test

import (
	"testing"

	"github.com/kanatsanan6/hrm/utils"
	"github.com/stretchr/testify/assert"
)

func TestPassword(t *testing.T) {
	password := utils.RandomString(10)

	hash, err := utils.Encrypt(password)

	assert.NoError(t, err)
	assert.NotNil(t, hash)

	matched := utils.ComparePasswords(hash, password)

	assert.True(t, matched)

	unmatched := utils.ComparePasswords(hash, "invalid")

	assert.False(t, unmatched)
}

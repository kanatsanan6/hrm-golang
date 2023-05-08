package queries_test

import (
	"testing"

	"github.com/kanatsanan6/hrm/model"
	"github.com/kanatsanan6/hrm/queries"
	"github.com/kanatsanan6/hrm/utils"
	"github.com/stretchr/testify/assert"
)

func GenerateUser() model.User {
	company, _ := testQueries.CreateCompany(queries.CreateCompanyArgs{
		Name: utils.RandomString(10),
	})
	user, _ := testQueries.CreateUser(queries.CreateUserArgs{
		Email:             utils.RandomEmail(),
		EncryptedPassword: utils.RandomString(16),
		FirstName:         utils.RandomString(10),
		LastName:          utils.RandomString(10),
		Role:              "admin",
		CompanyID:         &company.ID,
	})
	return user
}

func TestSQLQueries_CreateUser(t *testing.T) {
	company, err := testQueries.CreateCompany(queries.CreateCompanyArgs{
		Name: utils.RandomString(10),
	})
	assert.NoError(t, err)

	email := utils.RandomString(10)
	encryptedPassword := utils.RandomString(16)
	firstName := utils.RandomString(10)
	lastName := utils.RandomString(10)
	role := "admin"

	user, err := testQueries.CreateUser(queries.CreateUserArgs{
		Email:             email,
		EncryptedPassword: encryptedPassword,
		FirstName:         firstName,
		LastName:          lastName,
		CompanyID:         &company.ID,
		Role:              role,
	})

	assert.NoError(t, err)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, firstName, user.FirstName)
	assert.Equal(t, lastName, user.LastName)
	assert.Equal(t, role, user.Role)
	assert.Equal(t, company.ID, *user.CompanyID)
}

func TestSQLQueries_DeleteUser(t *testing.T) {
	user := GenerateUser()
	err := testQueries.DeleteUser(user.ID)
	assert.Error(t, err)
}

func TestSQLQueries_FindUserByID(t *testing.T) {
	user := GenerateUser()
	result, err := testQueries.FindUserByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.Email, result.Email)
	assert.Equal(t, user.FirstName, result.FirstName)
}

func TestSQLQueries_FindUserByEmail(t *testing.T) {
	user := GenerateUser()
	result, err := testQueries.FindUserByEmail(user.Email)
	assert.NoError(t, err)
	assert.Equal(t, user.Email, result.Email)
	assert.Equal(t, user.FirstName, result.FirstName)
}

func TestSQLQueries_FindUserByResetPasswordToken(t *testing.T) {
	user := GenerateUser()
	token := utils.RandomString(10)
	_, err := testQueries.UpdateUser(queries.UpdateUserArgs{
		ID:                 user.ID,
		ResetPasswordToken: &token,
	})
	assert.NoError(t, err)

	res, err := testQueries.FindUserByResetPasswordToken(token)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, res.ID)

}

func TestSQLQueries_UpdateUser(t *testing.T) {
	user := GenerateUser()
	token := utils.RandomString(10)
	res, err := testQueries.UpdateUser(queries.UpdateUserArgs{
		ID:                 user.ID,
		ResetPasswordToken: &token,
	})

	assert.NoError(t, err)
	assert.Equal(t, token, *res.ResetPasswordToken)
	assert.Equal(t, user.EncryptedPassword, res.EncryptedPassword)
	assert.Equal(t, user.CompanyID, res.CompanyID)
}

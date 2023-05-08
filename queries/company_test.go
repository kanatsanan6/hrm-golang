package queries_test

import (
	"testing"

	"github.com/kanatsanan6/hrm/queries"
	"github.com/kanatsanan6/hrm/utils"
	"github.com/stretchr/testify/assert"
)

func TestSQLQueries_CreateCompany(t *testing.T) {
	name := utils.RandomString(10)
	company, err := testQueries.CreateCompany(queries.CreateCompanyArgs{
		Name: name,
	})

	assert.NoError(t, err)
	assert.Equal(t, name, company.Name)
}

func TestSQLQueries_FindCompanyByID(t *testing.T) {
	company, err := testQueries.CreateCompany(queries.CreateCompanyArgs{
		Name: utils.RandomString(10),
	})
	assert.NoError(t, err)

	result, err := testQueries.FindCompanyByID(company.ID)
	assert.NoError(t, err)
	assert.Equal(t, company.ID, result.ID)
}

func TestSQLQueries_GetUsers(t *testing.T) {
	company, _ := testQueries.CreateCompany(queries.CreateCompanyArgs{
		Name: utils.RandomString(10),
	})
	admin, _ := testQueries.CreateUser(queries.CreateUserArgs{
		Email:             utils.RandomEmail(),
		EncryptedPassword: utils.RandomString(16),
		FirstName:         utils.RandomString(10),
		LastName:          utils.RandomString(10),
		Role:              "admin",
		CompanyID:         &company.ID,
	})
	member, _ := testQueries.CreateUser(queries.CreateUserArgs{
		Email:             utils.RandomEmail(),
		EncryptedPassword: utils.RandomString(16),
		FirstName:         utils.RandomString(10),
		LastName:          utils.RandomString(10),
		Role:              "member",
		CompanyID:         &company.ID,
	})

	result, err := testQueries.GetUsers(company.ID)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, admin.ID, result[0].ID)
	assert.Equal(t, member.ID, result[1].ID)
}

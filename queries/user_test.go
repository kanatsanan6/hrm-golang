package queries_test

import (
	"github.com/kanatsanan6/hrm/model"
	"github.com/kanatsanan6/hrm/queries"
	"github.com/kanatsanan6/hrm/utils"
)

func GenerateUser() *model.User {
	user, _ := testQueries.CreateUser(queries.CreateUserArgs{
		Email:             utils.RandomEmail(),
		EncryptedPassword: utils.RandomString(16),
		FirstName:         utils.RandomString(10),
		LastName:          utils.RandomString(10),
		Role:              "admin",
	})

	return &user
}

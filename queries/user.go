package queries

import "github.com/kanatsanan6/hrm/model"

type CreateUserArgs struct {
	Email             string
	EncryptedPassword string
	FirstName         string
	LastName          string
}

func (q *Queries) CreateUser(args CreateUserArgs) (model.User, error) {
	user := model.User{
		Email:             args.Email,
		EncryptedPassword: args.EncryptedPassword,
		FirstName:         args.FirstName,
		LastName:          args.LastName,
	}

	if err := q.DB.Create(&user).Error; err != nil {
		return model.User{}, err
	}
	return user, nil
}

package queries

import "github.com/kanatsanan6/hrm/model"

type CreateUserArgs struct {
	Email             string
	EncryptedPassword string
	FirstName         string
	LastName          string
}

func (q *SQLQueries) CreateUser(args CreateUserArgs) (model.User, error) {
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

func (q *SQLQueries) FindUserByEmail(email string) (model.User, error) {
	var user model.User
	if err := q.DB.Where("Email = ?", email).First(&user).Error; err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (q *SQLQueries) UpdateUserCompanyID(user model.User, id uint) error {
	user.CompanyID = &id
	return q.DB.Save(&user).Error
}

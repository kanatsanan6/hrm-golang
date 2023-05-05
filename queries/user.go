package queries

import (
	"time"

	"github.com/kanatsanan6/hrm/model"
)

type UserType struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Role      string    `json:"role"`
	CompanyID *uint     `json:"company_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserArgs struct {
	Email             string
	EncryptedPassword string
	FirstName         string
	LastName          string
	CompanyID         *uint
	Role              string
}

func (q *SQLQueries) CreateUser(args CreateUserArgs) (model.User, error) {
	user := model.User{
		Email:             args.Email,
		EncryptedPassword: args.EncryptedPassword,
		FirstName:         args.FirstName,
		LastName:          args.LastName,
		CompanyID:         args.CompanyID,
		Role:              args.Role,
	}

	if err := q.DB.Create(&user).Error; err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (q *SQLQueries) DeleteUser(user model.User) error {
	if err := q.DB.Delete(&user).Error; err != nil {
		return err
	}
	return nil
}

func (q *SQLQueries) FindUserByID(id uint) (model.User, error) {
	var user model.User
	if err := q.DB.Where("ID = ?", id).First(&user).Error; err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (q *SQLQueries) FindUserByEmail(email string) (model.User, error) {
	var user model.User
	err := q.DB.Where("Email = ?", email).
		Preload("Company").
		Preload("Leaves").
		Preload("LeaveTypes").
		First(&user).
		Error

	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (q *SQLQueries) FindUserByForgetPasswordToken(token string) (model.User, error) {
	var user model.User
	if err := q.DB.Where("reset_password_token = ?", token).First(&user).Error; err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (q *SQLQueries) UpdateUserCompanyID(user model.User, id uint) error {
	user.CompanyID = &id
	return q.DB.Save(&user).Error
}

func (q *SQLQueries) UpdateUserForgetPasswordToken(user model.User, token string) error {
	user.ResetPasswordToken = &token
	return q.DB.Save(&user).Error
}

func (q *SQLQueries) UpdateUserPassword(user model.User, hash string) error {
	user.EncryptedPassword = hash
	return q.DB.Save(&user).Error
}

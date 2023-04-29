package queries

import (
	"github.com/kanatsanan6/hrm/model"
	"gorm.io/gorm"
)

type Queries interface {
	CreateUser(args CreateUserArgs) (model.User, error)
	FindUserByEmail(email string) (model.User, error)
	FindUserByForgetPasswordToken(token string) (model.User, error)
	UpdateUserCompanyID(user model.User, id uint) error
	UpdateUserForgetPasswordToken(user model.User, token string) error
	UpdateUserPassword(user model.User, hash string) error
	CreateCompany(args CreateCompanyArgs) (model.Company, error)
	FindCompanyByID(id uint) (model.Company, error)
}

type SQLQueries struct {
	DB *gorm.DB
}

func NewQueries(db *gorm.DB) Queries {
	return &SQLQueries{DB: db}
}

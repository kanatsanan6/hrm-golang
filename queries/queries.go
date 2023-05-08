package queries

import (
	"github.com/jmoiron/sqlx"
	"github.com/kanatsanan6/hrm/model"
)

type Queries interface {
	// User
	FindUserByID(id int64) (model.User, error)
	FindUserByEmail(email string) (model.User, error)
	FindUserByResetPasswordToken(token string) (model.User, error)
	CreateUser(args CreateUserArgs) (model.User, error)
	DeleteUser(id int64) error
	UpdateUser(args UpdateUserArgs) (model.User, error)

	// Company
	CreateCompany(args CreateCompanyArgs) (model.Company, error)
	FindCompanyByID(id int64) (model.Company, error)
	GetUsers(companyID int64) ([]model.User, error)

	// Leave
	CreateLeave(args CreateLeaveArgs) (model.Leave, error)
	GetLeaves(user *model.User) ([]model.LeaveStruct, error)
	GetLeaveByID(id int64) (model.Leave, error)
	UpdateLeave(args UpdateLeaveArgs) (model.Leave, error)

	// Leave Type
	CreateLeaveType(args CreateLeaveTypeArgs) (model.LeaveType, error)
	FindUserLeaveTypeByName(user model.User, name string) (model.LeaveType, error)
	GetUserLeaveTypes(user model.User) ([]model.LeaveType, error)
}

type SQLQueries struct {
	DB *sqlx.DB
}

func NewQueries(db *sqlx.DB) Queries {
	return &SQLQueries{DB: db}
}

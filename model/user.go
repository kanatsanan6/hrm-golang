package model

import (
	"time"

	"github.com/kanatsanan6/hrm/utils"
)

type User struct {
	ID                 int64     `json:"id"`
	FirstName          string    `json:"first_name"`
	LastName           string    `json:"last_name"`
	Email              string    `json:"email"`
	EncryptedPassword  string    `json:"-"`
	ResetPasswordToken *string   `json:"-"`
	CompanyID          *int64    `json:"company_id"`
	Role               string    `json:"role"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func (u *User) GenerateResetPasswordToken() string {
	if u.ResetPasswordToken != nil {
		return *u.ResetPasswordToken
	}
	return utils.RandomString(16)
}

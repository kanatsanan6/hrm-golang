package model

import (
	"time"

	"github.com/kanatsanan6/hrm/utils"
)

type User struct {
	ID                 uint `gorm:"primaryKey"`
	FirstName          string
	LastName           string
	Email              string `gorm:"uniqueIndex"`
	EncryptedPassword  string
	ResetPasswordToken *string
	CompanyID          *uint
	Company            Company   `gorm:"foreignKey:CompanyID"`
	CreatedAt          time.Time `gorm:"autoCreateTime"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime"`
}

func (u *User) GenerateResetPasswordToken() string {
	if u.ResetPasswordToken != nil {
		return *u.ResetPasswordToken
	}
	return utils.RandomString(16)
}

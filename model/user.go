package model

import (
	"fmt"
	"time"

	"github.com/kanatsanan6/hrm/utils"
	"gorm.io/gorm"
)

type User struct {
	ID                 uint `gorm:"primaryKey"`
	FirstName          string
	LastName           string
	Email              string `gorm:"uniqueIndex"`
	EncryptedPassword  string
	ResetPasswordToken *string
	CompanyID          *uint
	Role               string
	Leaves             []Leave
	LeaveTypes         []LeaveType
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

func (u *User) AfterCreate(tx *gorm.DB) error {
	fmt.Println("from callback")
	return nil
}

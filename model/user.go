package model

import (
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
	for _, lType := range DefaultLeaveType {
		leaveType := LeaveType{
			Name:   lType["name"].(string),
			Usage:  0,
			Max:    lType["max"].(int),
			UserID: u.ID,
		}
		if err := tx.Create(&leaveType).Error; err != nil {
			return err
		}
	}
	return nil
}

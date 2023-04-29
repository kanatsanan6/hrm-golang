package queries

import (
	"github.com/kanatsanan6/hrm/model"
	"gorm.io/gorm"
)

type CreateCompanyArgs struct {
	Name string
}

func (q *SQLQueries) CreateCompany(args CreateCompanyArgs) (model.Company, error) {
	company := model.Company{Name: args.Name}

	if err := q.DB.Create(&company).Error; err != nil {
		return model.Company{}, err
	}
	return company, nil
}

func (q *SQLQueries) FindCompanyByID(id uint) (model.Company, error) {
	var company model.Company
	if err := q.DB.Where("ID = ?", id).Preload("Users", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at ASC")
	}).First(&company).Error; err != nil {
		return model.Company{}, err
	}
	return company, nil
}

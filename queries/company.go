package queries

import "github.com/kanatsanan6/hrm/model"

type CreateCompanyArgs struct {
	Name string
}

func (q *Queries) CreateCompany(args CreateCompanyArgs) (model.Company, error) {
	company := model.Company{Name: args.Name}

	if err := q.DB.Create(&company).Error; err != nil {
		return model.Company{}, err
	}
	return company, nil
}

package queries

import (
	"github.com/kanatsanan6/hrm/model"
)

type CreateCompanyArgs struct {
	Name string
}

func (q *SQLQueries) CreateCompany(args CreateCompanyArgs) (model.Company, error) {
	var company model.Company
	query := `
	INSERT INTO companies (name) VALUES ($1)
	RETURNING id, name, created_at, updated_at
	`
	row := q.DB.QueryRow(query, args.Name)
	err := row.Scan(
		&company.ID,
		&company.Name,
		&company.CreatedAt,
		&company.UpdatedAt,
	)
	return company, err
}

func (q *SQLQueries) FindCompanyByID(id int64) (model.Company, error) {
	var company model.Company
	row := q.DB.QueryRow(`SELECT * FROM companies WHERE (id) = ($1) LIMIT 1`, id)
	err := row.Scan(
		&company.ID,
		&company.Name,
		&company.CreatedAt,
		&company.UpdatedAt,
	)
	return company, err
}

func (q *SQLQueries) GetUsers(companyID int64) ([]model.User, error) {
	var users []model.User
	query := `
	SELECT id, email, first_name, last_name, company_id, role, created_at, updated_at
	FROM users WHERE (company_id) = ($1)
	`
	row, err := q.DB.Query(query, companyID)
	if err != nil {
		return []model.User{}, err
	}
	for row.Next() {
		var user model.User
		err := row.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.CompanyID,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return []model.User{}, err
		}
		users = append(users, user)
	}
	return users, nil
}

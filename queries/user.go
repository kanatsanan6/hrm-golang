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
	CompanyID         *int64
	Role              string
}

func (q *SQLQueries) CreateUser(args CreateUserArgs) (model.User, error) {
	var user model.User
	query := `
	INSERT INTO users (
		email,
		encrypted_password,
		first_name,
		last_name,
		company_id,
		role
	)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, email, first_name, last_name, company_id, role, created_at, updated_at
	`
	row := q.DB.QueryRow(
		query,
		args.Email,
		args.EncryptedPassword,
		args.FirstName,
		args.LastName,
		args.CompanyID,
		args.Role,
	)
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
	return user, err
}

func (q *SQLQueries) DeleteUser(id int64) error {
	err := q.DB.QueryRow(`DELETE FROM Users WHERE (user_id) = ($1)`, id)
	return err.Err()
}

func (q *SQLQueries) FindUserByID(id int64) (model.User, error) {
	var user model.User
	query := `
	SELECT
	id,
	email,
	first_name, last_name, company_id, reset_password_token, role, created_at, updated_at
	FROM users WHERE (id) = ($1) LIMIT 1
	`
	row := q.DB.QueryRow(query, id)
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.CompanyID,
		&user.ResetPasswordToken,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	return user, err
}

func (q *SQLQueries) FindUserByEmail(email string) (model.User, error) {
	var user model.User
	query := `
	SELECT id, email, first_name, last_name, company_id, encrypted_password, reset_password_token, role, created_at, updated_at
	FROM users WHERE (email) = ($1) LIMIT 1
	`
	row := q.DB.QueryRow(query, email)
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.CompanyID,
		&user.EncryptedPassword,
		&user.ResetPasswordToken,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	return user, err
}

func (q *SQLQueries) FindUserByResetPasswordToken(token string) (model.User, error) {
	var user model.User
	query := `
	SELECT id, email, first_name, last_name, company_id, reset_password_token, role, created_at, updated_at
	FROM users WHERE (reset_password_token) = ($1) LIMIT 1
	`
	row := q.DB.QueryRow(query, token)
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.CompanyID,
		&user.ResetPasswordToken,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	return user, err
}

type UpdateUserArgs struct {
	ID                 int64
	CompanyID          *int64
	ResetPasswordToken *string
	Password           *string
}

func (q *SQLQueries) UpdateUser(args UpdateUserArgs) (model.User, error) {
	var user model.User
	query := `
	UPDATE users
	SET
	company_id = COALESCE($2, company_id),
	reset_password_token = COALESCE($3, reset_password_token),
	encrypted_password = COALESCE($4, encrypted_password)
	WHERE id = $1
	RETURNING id, email, first_name, last_name, company_id, reset_password_token, role, created_at, updated_at
	`
	row := q.DB.QueryRow(query, args.ID, args.CompanyID, args.ResetPasswordToken, args.Password)
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.CompanyID,
		&user.ResetPasswordToken,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	return user, err
}

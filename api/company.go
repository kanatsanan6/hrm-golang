package api

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kanatsanan6/hrm/model"
	"github.com/kanatsanan6/hrm/queries"
	"github.com/kanatsanan6/hrm/utils"
)

type CompanyType struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func CompanyResponse(company model.Company) CompanyType {
	return CompanyType{
		ID:        company.ID,
		Name:      company.Name,
		CreatedAt: company.CreatedAt,
		UpdatedAt: company.UpdatedAt,
	}
}

type CreateCompanyBody struct {
	Name string `json:"name" validate:"required"`
}

func (s *Server) createCompany(c *fiber.Ctx) error {
	var body CreateCompanyBody
	if err := c.BodyParser(&body); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	if err := utils.ValidateStruct(body); len(err) != 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err)
	}

	email := c.Locals("email").(string)

	user, err := s.Queries.FindUserByEmail(email)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, err.Error())
	}

	company, err := s.Queries.CreateCompany(queries.CreateCompanyArgs{Name: body.Name})
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	if err := s.Queries.UpdateUserCompanyID(user, company.ID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	return utils.JsonResponse(c, fiber.StatusCreated, CompanyResponse(company))
}

func (s *Server) getUsers(c *fiber.Ctx) error {
	email := c.Locals("email").(string)

	user, err := s.Queries.FindUserByEmail(email)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, err.Error())
	}

	company, err := s.Queries.FindCompanyByID(*user.CompanyID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, err.Error())
	}

	users := []UserType{}
	for _, user := range company.Users {
		users = append(users, userResponse(user))
	}
	return utils.JsonResponse(c, fiber.StatusOK, users)
}

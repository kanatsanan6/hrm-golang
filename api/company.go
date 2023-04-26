package api

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kanatsanan6/hrm/model"
	"github.com/kanatsanan6/hrm/queries"
	"github.com/kanatsanan6/hrm/utils"
)

type companyType struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func companyResponse(company model.Company) companyType {
	return companyType{
		ID:        company.ID,
		Name:      company.Name,
		CreatedAt: company.CreatedAt,
		UpdatedAt: company.UpdatedAt,
	}
}

type CreateCompanyBody struct {
	Name string `json:"name"`
}

func (s *Server) createCompany(c *fiber.Ctx) error {
	var body CreateCompanyBody
	if err := c.BodyParser(&body); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	email := c.Locals("email").(string)

	user, err := s.Queries.FindUserByEmail(email)
	if err != nil {
		utils.ErrorResponse(c, fiber.StatusNotFound, err.Error())
	}

	company, err := s.Queries.CreateCompany(queries.CreateCompanyArgs{Name: body.Name})
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	user.CompanyID = &company.ID
	s.Queries.DB.Save(&user)

	return utils.JsonResponse(c, fiber.StatusCreated, companyResponse(company))
}

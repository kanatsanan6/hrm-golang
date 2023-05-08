package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kanatsanan6/hrm/model"
	"github.com/kanatsanan6/hrm/utils"
)

func (s *Server) getUsers(c *fiber.Ctx) error {
	user := c.Locals("user").(model.User)

	company, err := s.Queries.FindCompanyByID(*user.CompanyID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, err.Error())
	}

	users, err := s.Queries.GetUsers(company.ID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, err.Error())
	}
	return utils.JsonResponse(c, fiber.StatusOK, users)
}

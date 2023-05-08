package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kanatsanan6/hrm/model"
	"github.com/kanatsanan6/hrm/utils"
)

func (s *Server) getLeaveTypes(c *fiber.Ctx) error {
	user := c.Locals("user").(model.User)
	leaveType, _ := s.Queries.GetUserLeaveTypes(user)
	return utils.JsonResponse(c, fiber.StatusOK, leaveType)
}

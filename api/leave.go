package api

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kanatsanan6/hrm/model"
	"github.com/kanatsanan6/hrm/queries"
	"github.com/kanatsanan6/hrm/utils"
)

type LeaveType struct {
	ID          uint      `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	LeaveType   string    `json:"leave_type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func leaveResponse(leave model.Leave) LeaveType {
	return LeaveType{
		ID:          leave.ID,
		Description: leave.Description,
		LeaveType:   leave.LeaveType,
		Status:      leave.Status,
		StartDate:   leave.StartDate,
		EndDate:     leave.EndDate,
		CreatedAt:   leave.CreatedAt,
		UpdatedAt:   leave.UpdatedAt,
	}
}

type LeaveBody struct {
	Description string `json:"description" validate:"required"`
	LeaveType   string `json:"leave_type" validate:"required"`
	StartDate   string `json:"start_date" validate:"required"`
	EndDate     string `json:"end_date" validate:"required"`
}

func (s *Server) createLeave(c *fiber.Ctx) error {
	var body LeaveBody
	if err := c.BodyParser(&body); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	if err := utils.ValidateStruct(body); len(err) != 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err)
	}

	startDate, err := utils.StringToDateTime(body.StartDate)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err)
	}

	endDate, err := utils.StringToDateTime(body.EndDate)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err)
	}

	user := c.Locals("user").(model.User)

	leave, err := s.Queries.CreateLeave(queries.CreateLeaveArgs{
		Description: body.Description,
		Status:      "pending",
		StartDate:   startDate,
		EndDate:     endDate,
		LeaveType:   body.LeaveType,
		UserID:      user.ID,
	})
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	return utils.JsonResponse(c, fiber.StatusCreated, leaveResponse(leave))
}

func (s *Server) getLeaves(c *fiber.Ctx) error {
	user := c.Locals("user").(model.User)
	leaves := s.Queries.GetLeaves(&user)
	return utils.JsonResponse(c, fiber.StatusOK, leaves)
}

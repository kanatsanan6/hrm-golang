package api

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kanatsanan6/hrm/model"
	"github.com/kanatsanan6/hrm/queries"
	"github.com/kanatsanan6/hrm/utils"
)

type LeaveStruct struct {
	ID          int64           `json:"id"`
	Description string          `json:"description"`
	Status      string          `json:"status"`
	StartDate   time.Time       `json:"start_date"`
	EndDate     time.Time       `json:"end_date"`
	LeaveType   model.LeaveType `json:"leave_type"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

func LeaveResponse(leave model.Leave, leaveType model.LeaveType) *LeaveStruct {
	return &LeaveStruct{
		ID:          leave.ID,
		Description: leave.Description,
		Status:      leave.Status,
		StartDate:   leave.StartDate,
		EndDate:     leave.EndDate,
		LeaveType:   leaveType,
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
	leaveType, err := s.Queries.FindUserLeaveTypeByName(user, body.LeaveType)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "leave_type not found")
	}

	var status string
	if user.Role == "admin" {
		status = "approved"
	} else {
		status = "pending"
	}

	leave, err := s.Queries.CreateLeave(queries.CreateLeaveArgs{
		Description: body.Description,
		Status:      status,
		StartDate:   startDate,
		EndDate:     endDate,
		LeaveTypeID: leaveType.ID,
		UserID:      user.ID,
	})
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	return utils.JsonResponse(c, fiber.StatusCreated, LeaveResponse(leave, leaveType))
}

func (s *Server) getLeaves(c *fiber.Ctx) error {
	user := c.Locals("user").(model.User)
	leaves, err := s.Queries.GetLeaves(&user)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	return utils.JsonResponse(c, fiber.StatusOK, leaves)
}

type UpdateLeaveBody struct {
	Status string `json:"status" validate:"required"`
}

type UpdateLeaveParam struct {
	ID uint `json:"id" validate:"required"`
}

func (s *Server) updateLeaveStatus(c *fiber.Ctx) error {
	var body UpdateLeaveBody
	if err := c.BodyParser(&body); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	if err := utils.ValidateStruct(body); len(err) != 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err)
	}

	var params UpdateLeaveParam
	if err := c.ParamsParser(&params); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	if err := utils.ValidateStruct(params); len(err) != 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err)
	}

	// leave, err := s.Queries.GetLeaveByID(params.ID)
	// if err != nil {
	// 	return utils.ErrorResponse(c, fiber.StatusNotFound, err.Error())
	// }

	// if err := s.Queries.UpdateLeaveStatus(&leave, body.Status); err != nil {
	// 	return utils.ErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
	// }
	return utils.JsonResponse(c, fiber.StatusCreated, "test")
}

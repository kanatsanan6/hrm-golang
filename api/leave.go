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
	StartPeriod string    `json:"start_period"`
	EndDate     time.Time `json:"end_date"`
	EndPeriod   string    `json:"end_period"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

func leaveResponse(leave model.Leave) LeaveType {
	return LeaveType{
		ID:          leave.ID,
		Description: leave.Description,
		Status:      leave.Status,
		StartDate:   leave.StartDate,
		StartPeriod: leave.StartPeriod,
		EndDate:     leave.EndDate,
		EndPeriod:   leave.EndPeriod,
		CreatedAt:   leave.CreatedAt,
		UpdatedAt:   leave.UpdatedAt,
	}
}

type LeaveBody struct {
	Description string    `json:"description" validate:"required"`
	StartDate   time.Time `json:"start_date" validate:"required"`
	StartPeriod string    `json:"start_period" validate:"required"`
	EndDate     time.Time `json:"end_date" validate:"required"`
	EndPeriod   string    `json:"end_period" validate:"required"`
}

func (s *Server) createLeave(c *fiber.Ctx) error {
	var body LeaveBody
	if err := c.BodyParser(&body); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	if err := utils.ValidateStruct(body); len(err) != 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err)
	}

	leave, err := s.Queries.CreateLeave(queries.CreateLeaveArgs{
		Description: body.Description,
		Status:      "pending",
		StartDate:   body.StartDate,
		StartPeriod: body.StartPeriod,
		EndDate:     body.EndDate,
		EndPeriod:   body.EndPeriod,
	})
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	return utils.JsonResponse(c, fiber.StatusCreated, leaveResponse(leave))
}

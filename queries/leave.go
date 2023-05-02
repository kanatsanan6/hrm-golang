package queries

import (
	"time"

	"github.com/kanatsanan6/hrm/model"
)

type CreateLeaveArgs struct {
	Description string
	Status      string
	StartDate   time.Time
	StartPeriod string
	EndDate     time.Time
	EndPeriod   string
}

func (q *SQLQueries) CreateLeave(args CreateLeaveArgs) (model.Leave, error) {
	leave := model.Leave{
		Description: args.Description,
		Status:      args.Status,
		StartDate:   args.StartDate,
		StartPeriod: args.StartPeriod,
		EndDate:     args.EndDate,
		EndPeriod:   args.EndPeriod,
	}

	if err := q.DB.Create(&leave).Error; err != nil {
		return model.Leave{}, err
	}
	return leave, nil
}

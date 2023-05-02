package queries

import (
	"time"

	"github.com/kanatsanan6/hrm/model"
)

type CreateLeaveArgs struct {
	Description string
	Status      string
	StartDate   time.Time
	EndDate     time.Time
}

func (q *SQLQueries) CreateLeave(args CreateLeaveArgs) (model.Leave, error) {
	leave := model.Leave{
		Description: args.Description,
		Status:      args.Status,
		StartDate:   args.StartDate,
		EndDate:     args.EndDate,
	}

	if err := q.DB.Create(&leave).Error; err != nil {
		return model.Leave{}, err
	}
	return leave, nil
}

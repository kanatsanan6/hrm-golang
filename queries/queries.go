package queries

import (
	"gorm.io/gorm"
)

type Queries struct {
	DB *gorm.DB
}

func NewQueries(db *gorm.DB) *Queries {
	return &Queries{DB: db}
}

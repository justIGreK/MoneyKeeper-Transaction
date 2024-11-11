package models

import "time"

type CreateTransaction struct {
	Category string
	UserID   string  `validate:"required"`
	Name     string  `validate:"required"`
	Cost     float64 `validate:"required"`
}

type Transaction struct {
	ID       string    `bson:"_id,omitempty"`
	UserID   string    `bson:"user_id"`
	Category string    `bson:"category"`
	Name     string    `bson:"name"`
	Cost     float64   `bson:"cost"`
	Date     time.Time `bson:"date"`
}

type TimeFrame struct {
	StartDate time.Time
	EndDate   time.Time
}

type CreateTimeFrame struct {
	StartDate string
	EndDate   string
}

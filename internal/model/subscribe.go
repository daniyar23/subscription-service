package model

import (
	"github.com/google/uuid"
)

// Subscription represents a subscription for a user to a service.
type Subscription struct {
	ID          uuid.UUID  `json:"id"`
	ServiceName string     `json:"service_name"`
	Price       int        `json:"price"`
	UserID      uuid.UUID  `json:"user_id"`
	StartDate   MonthYear  `json:"start_date"`
	EndDate     *MonthYear `json:"end_date"`
}

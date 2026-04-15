package subscriptions

import (
	"github.com/google/uuid"
)

type SubscriptionsCreateRequest struct {
	ServiceName string    `json:"service_name" validate:"required"`
	Price       int       `json:"price" validate:"required"`
	UserID      uuid.UUID `json:"user_id" validate:"required"`
	StartDate   string    `json:"start_date" validate:"required"`
	EndDate     *string   `json:"end_date"`
}

type SubscriptionsUpdateRequest struct {
	ServiceName string  `json:"service_name" validate:"required"`
	Price       int     `json:"price" validate:"required"`
	StartDate   string  `json:"start_date" validate:"required"`
	EndDate     *string `json:"end_date"`
}

type SubscriptionTotalRequest struct {
	ServiceName string    `json:"service_name"`
	UserID      uuid.UUID `json:"user_id"`
	From        string    `json:"from"`
	To          string    `json:"to"`
}

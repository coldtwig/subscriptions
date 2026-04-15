package subscriptions

import (
	"time"

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

type SubscriptionResponse struct {
	ID          uint       `json:"id"`
	ServiceName string     `json:"service_name"`
	Price       int        `json:"price"`
	UserID      uuid.UUID  `json:"user_id"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

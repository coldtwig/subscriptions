package subscriptions

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ServiceName string    `json:"service_name" validate:"required"`
	Price       int       `json:"price" validate:"required"`
	UserId      uuid.UUID `json:"user_id" validate:"required"`
	StartDate   time.Time `json:"start_date" validate:"required"`
	EndDate     time.Time `json:"end_date"`
}

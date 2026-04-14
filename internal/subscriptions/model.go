package subscriptions

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Subscription struct {
	gorm.Model
	ServiceName string     `json:"service_name" validate:"required"`
	Price       int        `json:"price" validate:"required"`
	UserID      uuid.UUID  `json:"user_id" validate:"required"`
	StartDate   time.Time  `json:"start_date" validate:"required"`
	EndDate     *time.Time `json:"end_date"`
}

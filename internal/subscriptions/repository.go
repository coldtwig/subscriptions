package subscriptions

import (
	"subscriptions/internal/db"

	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

type SubscriptionsRepository struct {
	Database *db.DB
}

func NewSubscriptionsRepository(database *db.DB) *SubscriptionsRepository {
	return &SubscriptionsRepository{
		Database: database,
	}
}

func (repo *SubscriptionsRepository) Create(subscription *Subscription) (*Subscription, error) {
	result := repo.Database.DB.Create(subscription)
	if result.Error != nil {
		return nil, result.Error
	}

	return subscription, nil
}

func (repo *SubscriptionsRepository) GetAll() ([]Subscription, error) {
	var subscriptions []Subscription

	result := repo.Database.Table("subscriptions").
		Where("deleted_at is null").
		Scan(&subscriptions)

	if result.Error != nil {
		return nil, result.Error
	}

	return subscriptions, nil
}

func (repo *SubscriptionsRepository) Update(subscription *Subscription) (*Subscription, error) {
	result := repo.Database.DB.Clauses(clause.Returning{}).Updates(subscription)
	if result.Error != nil {
		return nil, result.Error
	}

	return subscription, nil
}

func (repo *SubscriptionsRepository) Delete(id uint) error {
	result := repo.Database.DB.Delete(&Subscription{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *SubscriptionsRepository) FindForTotal(subTotal *SubscriptionTotalFilter) ([]Subscription, error) {
	var subscriptions []Subscription

	query := repo.Database.DB.Model(&Subscription{}).
		Where("deleted_at is null").
		Where("start_date <= ?", subTotal.To).
		Where("(end_date is null OR end_date >= ?)", subTotal.From)

	if subTotal.ServiceName != "" {
		query = query.Where("service_name = ?", subTotal.ServiceName)
	}

	if subTotal.UserID != uuid.Nil {
		query = query.Where("user_id = ?", subTotal.UserID)
	}

	result := query.Find(&subscriptions)
	if result.Error != nil {
		return nil, result.Error
	}

	return subscriptions, nil
}

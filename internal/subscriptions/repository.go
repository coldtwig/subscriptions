package subscriptions

import (
	"subscriptions/internal/db"

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

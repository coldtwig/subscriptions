package subscriptions

import "time"
import "log"

type SubscriptionsService struct {
	SubscriptionsRepository *SubscriptionsRepository
}

func NewSubscriptionsService(subscriptionsRepository *SubscriptionsRepository) *SubscriptionsService {
	return &SubscriptionsService{
		SubscriptionsRepository: subscriptionsRepository,
	}
}

func (service *SubscriptionsService) SumAll(subTotal *SubscriptionTotalFilter) (int, error) {
	subscriptions, err := service.SubscriptionsRepository.FindForTotal(subTotal)
	if err != nil {
		log.Printf(
			"ERROR service sum total failed user_id=%s service_name=%s from=%s to=%s err=%v",
			subTotal.UserID,
			subTotal.ServiceName,
			subTotal.From.Format("01-2006"),
			subTotal.To.Format("01-2006"),
			err,
		)
		return 0, err
	}

	sum := calculateTotalForPeriod(subscriptions, subTotal.From, subTotal.To)
	log.Printf(
		"INFO total calculated user_id=%s service_name=%s from=%s to=%s subscriptions=%d total=%d",
		subTotal.UserID,
		subTotal.ServiceName,
		subTotal.From.Format("01-2006"),
		subTotal.To.Format("01-2006"),
		len(subscriptions),
		sum,
	)
	return sum, nil
}

func calculateTotalForPeriod(subscriptions []Subscription, from time.Time, to time.Time) int {
	sum := 0

	for _, s := range subscriptions {
		months := overlapMonths(s, from, to)
		sum += s.Price * months
	}

	return sum
}

func overlapMonths(subscription Subscription, from time.Time, to time.Time) int {
	overlapStart := time.Date(from.Year(), from.Month(), 1, 0, 0, 0, 0, time.UTC)
	subscriptionStart := time.Date(subscription.StartDate.Year(), subscription.StartDate.Month(), 1, 0, 0, 0, 0, time.UTC)
	if subscriptionStart.After(overlapStart) {
		overlapStart = subscriptionStart
	}

	subscriptionEnd := to
	if subscription.EndDate != nil && subscription.EndDate.Before(to) {
		subscriptionEnd = *subscription.EndDate
	}

	overlapEnd := time.Date(to.Year(), to.Month(), 1, 0, 0, 0, 0, time.UTC)
	subscriptionEndMonth := time.Date(subscriptionEnd.Year(), subscriptionEnd.Month(), 1, 0, 0, 0, 0, time.UTC)
	if subscriptionEndMonth.Before(overlapEnd) {
		overlapEnd = subscriptionEndMonth
	}

	if overlapStart.After(overlapEnd) {
		return 0
	}

	yearDiff := overlapEnd.Year() - overlapStart.Year()
	monthDiff := int(overlapEnd.Month()) - int(overlapStart.Month())
	return yearDiff*12 + monthDiff + 1
}

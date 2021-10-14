package service

import (
	"go-subscriptions-workflow/types"
	"time"
)

const (
	TaskQueueName            = "SubscriptionsTaskQueue"
	QuerySubscriptionState   = "QuerySubscriptionState"
	SignalCancelSubscription = "SignalCancelSubscription"
)

type SubscriptionState struct {
	ID          string
	UserID      string
	Price       float64
	Features    []*Feature
	Activations int
	ActivatedAt time.Time
	ExpiresAt   time.Time
	Canceled    bool
	CanceledAt  *time.Time
	Disabled    bool
	DisabledAt  *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Feature struct {
	Name string
}

func (s *SubscriptionState) HasExpired(t time.Time) bool {
	return s.ExpiresAt.Before(t)
}

func (s *SubscriptionState) Expiration() time.Duration {
	return s.ExpiresAt.Sub(s.ActivatedAt)
}

func NewState(subscription *types.SubscriptionOutput) SubscriptionState {
	state := SubscriptionState{
		ID:          subscription.ID,
		UserID:      subscription.UserID,
		Price:       subscription.Price,
		Activations: subscription.Activations,
		ActivatedAt: subscription.ActivatedAt,
		ExpiresAt:   subscription.ExpiresAt,
		Canceled:    subscription.Canceled,
		CanceledAt:  subscription.CanceledAt,
		Disabled:    subscription.Disabled,
		DisabledAt:  subscription.DisabledAt,
		CreatedAt:   subscription.CreatedAt,
		UpdatedAt:   subscription.UpdatedAt,
	}
	state.Features = make([]*Feature, 0, len(subscription.Features))
	for index := range subscription.Features {
		feature := &Feature{Name: subscription.Features[index].Name}
		state.Features = append(state.Features, feature)
	}
	return state
}

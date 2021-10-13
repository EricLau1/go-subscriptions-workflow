package service

import (
	"context"
	"fmt"
	"go-subscriptions-workflow/db"
	"go-subscriptions-workflow/services/subscriptions/models"
	"go-subscriptions-workflow/services/subscriptions/rules"
	"go-subscriptions-workflow/services/subscriptions/store"
	userssvc "go-subscriptions-workflow/services/users/service"
	"go-subscriptions-workflow/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type SubscriptionsClient interface {
	GetSubscriptions(ctx context.Context) ([]*types.SubscriptionOutput, error)
}

type SubscriptionsServiceServer interface {
	Start(ctx context.Context, req *types.StartSubscriptionRequest) error
	Charge(ctx context.Context, req *types.ChargeSubscriptionRequest) error
	Cancel(ctx context.Context, req *types.CancelSubscriptionRequest) error
	Disable(ctx context.Context, req *types.DisableSubscriptionRequest) error
	GetSubscriptions(ctx context.Context) ([]*types.SubscriptionOutput, error)
}

type subscriptionsService struct {
	usersService       userssvc.UsersService
	subscriptionsStore store.SubscriptionsStore
}

func NewSubscriptionsClient(dbConn db.Connection, usersService userssvc.UsersService) SubscriptionsClient {
	return &subscriptionsService{
		usersService:       usersService,
		subscriptionsStore: store.NewSubscriptionsStore(dbConn.DB()),
	}
}

func NewSubscriptionsServiceServer(dbConn db.Connection, usersService userssvc.UsersService) SubscriptionsServiceServer {
	return &subscriptionsService{
		usersService:       usersService,
		subscriptionsStore: store.NewSubscriptionsStore(dbConn.DB()),
	}
}

func (s *subscriptionsService) Start(ctx context.Context, req *types.StartSubscriptionRequest) error {
	user, err := s.usersService.GetUser(ctx, req.UserID)
	if err != nil {
		return err
	}
	if user.Balance < rules.DefaultPrice {
		return fmt.Errorf("insufficient funds to subscribe: user_id=%v, balance=%f, price=%f", user.ID, user.Balance, rules.DefaultPrice)
	}
	userID, _ := primitive.ObjectIDFromHex(user.ID)
	subscriptions, err := s.subscriptionsStore.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}
	for index := range subscriptions {
		if !subscriptions[index].Canceled && !subscriptions[index].Disabled {
			return fmt.Errorf("subscription already started: subscription_id=%v, user_id=%v",
				subscriptions[index].ID, subscriptions[index].UserID)
		}
	}

	id := primitive.NewObjectID()

	subscription := &models.Subscription{
		ID:     id,
		UserID: userID,
		Price:  rules.DefaultPrice,
		Features: []*models.Feature{
			{
				Name: "downloads",
			},
			{
				Name: "uploads",
			},
		},
		Activations: 1,
		ActivatedAt: time.Now(),
		ExpiresAt:   time.Now().Add(rules.DefaultExpiration),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = s.subscriptionsStore.Create(ctx, subscription)
	if err != nil {
		return err
	}

	return nil
}

func (s *subscriptionsService) Charge(ctx context.Context, req *types.ChargeSubscriptionRequest) error {
	id, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		return err
	}
	subscription, err := s.subscriptionsStore.Get(ctx, id)
	if err != nil {
		return err
	}
	user, err := s.usersService.GetUser(ctx, subscription.UserID.Hex())
	if err != nil {
		return err
	}
	if user.Balance < subscription.Price {
		return rules.ErrInsufficientFunds
	}
	debit := new(types.DebitInput)
	debit.Amount = subscription.Price
	debit.UserID = subscription.UserID.Hex()
	_, err = s.usersService.Debit(ctx, debit)
	if err != nil {
		return err
	}
	subscription.Activations++
	subscription.ActivatedAt = time.Now()
	subscription.ExpiresAt = time.Now().Add(rules.DefaultExpiration)
	subscription.UpdatedAt = time.Now()
	return s.subscriptionsStore.Update(ctx, subscription)
}

func (s *subscriptionsService) Cancel(ctx context.Context, req *types.CancelSubscriptionRequest) error {
	id, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		return err
	}
	subscription, err := s.subscriptionsStore.Get(ctx, id)
	if err != nil {
		return err
	}
	if subscription.UserID.Hex() != req.UserID {
		return fmt.Errorf("invalid user to cancel subscription: user_id=%v, subscription_id=%v",
			req.UserID, req.ID)
	}
	if subscription.Canceled {
		return nil
	}
	subscription.Canceled = true
	canceledAt := time.Now()
	subscription.CanceledAt = &canceledAt
	subscription.UpdatedAt = time.Now()
	return s.subscriptionsStore.Update(ctx, subscription)
}

func (s *subscriptionsService) Disable(ctx context.Context, req *types.DisableSubscriptionRequest) error {
	id, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		return err
	}
	subscription, err := s.subscriptionsStore.Get(ctx, id)
	if err != nil {
		return err
	}
	subscription.Disabled = true
	disabledAt := time.Now()
	subscription.DisabledAt = &disabledAt
	subscription.UpdatedAt = time.Now()
	return s.subscriptionsStore.Update(ctx, subscription)
}

func (s *subscriptionsService) GetSubscriptions(ctx context.Context) ([]*types.SubscriptionOutput, error) {
	subscriptions, err := s.subscriptionsStore.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*types.SubscriptionOutput, 0, len(subscriptions))
	for index := range subscriptions {
		out = append(out, subscriptions[index].Out())
	}
	return out, nil
}

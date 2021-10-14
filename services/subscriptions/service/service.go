package service

import (
	"context"
	"fmt"
	"go-subscriptions-workflow/db"
	"go-subscriptions-workflow/services/subscriptions/models"
	"go-subscriptions-workflow/services/subscriptions/shared"
	"go-subscriptions-workflow/services/subscriptions/store"
	userssvc "go-subscriptions-workflow/services/users/service"
	"go-subscriptions-workflow/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.temporal.io/sdk/client"
	"log"
	"time"
)

type SubscriptionsClient interface {
	GetSubscriptions(ctx context.Context) ([]*types.SubscriptionOutput, error)
}

type SubscriptionsServiceServer interface {
	Start(ctx context.Context, req *types.StartSubscriptionRequest) (*types.SubscriptionOutput, error)
	Charge(ctx context.Context, req *types.ChargeSubscriptionRequest) (*types.SubscriptionOutput, error)
	Cancel(ctx context.Context, req *types.CancelSubscriptionRequest) (*types.SubscriptionOutput, error)
	Disable(ctx context.Context, req *types.DisableSubscriptionRequest) (*types.SubscriptionOutput, error)
	GetSubscriptions(ctx context.Context) ([]*types.SubscriptionOutput, error)
}

type subscriptionsService struct {
	usersService       userssvc.UsersService
	subscriptionsStore store.SubscriptionsStore
	temporalClient     client.Client
}

func NewSubscriptionsClient(dbConn db.Connection, usersService userssvc.UsersService) SubscriptionsClient {
	return &subscriptionsService{
		usersService:       usersService,
		subscriptionsStore: store.NewSubscriptionsStore(dbConn.DB()),
	}
}

func NewSubscriptionsServiceServer(dbConn db.Connection, usersService userssvc.UsersService, temporalClient client.Client) SubscriptionsServiceServer {
	return &subscriptionsService{
		usersService:       usersService,
		subscriptionsStore: store.NewSubscriptionsStore(dbConn.DB()),
		temporalClient:     temporalClient,
	}
}

func (s *subscriptionsService) Start(ctx context.Context, req *types.StartSubscriptionRequest) (*types.SubscriptionOutput, error) {
	user, err := s.usersService.GetUser(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	if user.Balance < shared.DefaultPrice {
		return nil, fmt.Errorf("insufficient funds to subscribe: user_id=%v, balance=%f, price=%f", user.ID, user.Balance, shared.DefaultPrice)
	}
	userID, _ := primitive.ObjectIDFromHex(user.ID)
	subscriptions, err := s.subscriptionsStore.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	for index := range subscriptions {
		if !subscriptions[index].Canceled && !subscriptions[index].Disabled {
			return nil, fmt.Errorf("subscription already started: subscription_id=%v, user_id=%v",
				subscriptions[index].ID, subscriptions[index].UserID)
		}
	}

	id := primitive.NewObjectID()

	subscription := &models.Subscription{
		ID:     id,
		UserID: userID,
		Price:  shared.DefaultPrice,
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
		ExpiresAt:   time.Now().Add(shared.DefaultExpiration),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = s.subscriptionsStore.Create(ctx, subscription)
	if err != nil {
		return nil, err
	}

	out := subscription.Out()

	state := NewState(out)

	options := client.StartWorkflowOptions{
		ID:                 state.ID,
		TaskQueue:          TaskQueueName,
		WorkflowRunTimeout: time.Hour * 24 * 30 * 6,
	}

	we, err := s.temporalClient.ExecuteWorkflow(ctx, options, SubscriptionsWorkflow, state, &Activities{s})
	if err != nil {
		return nil, err
	}

	log.Printf("execute workflow: ID=%v, RunID=%v\n", we.GetID(), we.GetRunID())

	return out, nil
}

func (s *subscriptionsService) Charge(ctx context.Context, req *types.ChargeSubscriptionRequest) (*types.SubscriptionOutput, error) {
	id, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		return nil, err
	}
	subscription, err := s.subscriptionsStore.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	user, err := s.usersService.GetUser(ctx, subscription.UserID.Hex())
	if err != nil {
		return nil, err
	}
	if user.Balance < subscription.Price {
		return nil, shared.ErrInsufficientFunds
	}

	debit := new(types.DebitInput)
	debit.Amount = subscription.Price
	debit.UserID = subscription.UserID.Hex()

	_, err = s.usersService.Debit(ctx, debit)
	if err != nil {
		return nil, err
	}

	subscription.Activations++
	subscription.ActivatedAt = time.Now()
	subscription.ExpiresAt = time.Now().Add(shared.DefaultExpiration)
	subscription.UpdatedAt = time.Now()

	err = s.subscriptionsStore.Update(ctx, subscription)
	if err != nil {
		return nil, err
	}

	log.Println("subscription charged: ", subscription.ID.Hex())

	return subscription.Out(), nil
}

func (s *subscriptionsService) Cancel(ctx context.Context, req *types.CancelSubscriptionRequest) (*types.SubscriptionOutput, error) {
	id, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		return nil, err
	}
	subscription, err := s.subscriptionsStore.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if subscription.UserID.Hex() != req.UserID {
		return nil, fmt.Errorf("invalid user to cancel subscription: user_id=%v, subscription_id=%v",
			req.UserID, req.ID)
	}
	if subscription.Canceled {
		return nil, nil
	}

	subscription.Canceled = true
	canceledAt := time.Now()
	subscription.CanceledAt = &canceledAt
	subscription.UpdatedAt = time.Now()

	err = s.subscriptionsStore.Update(ctx, subscription)
	if err != nil {
		return nil, err
	}

	log.Println("subscription canceled: ", subscription.ID.Hex())

	err = s.temporalClient.SignalWorkflow(ctx, req.ID, "", SignalCancelSubscription, subscription.Canceled)
	if err != nil {
		return nil, err
	}

	return subscription.Out(), err
}

func (s *subscriptionsService) Disable(ctx context.Context, req *types.DisableSubscriptionRequest) (*types.SubscriptionOutput, error) {
	id, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		return nil, err
	}
	subscription, err := s.subscriptionsStore.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	subscription.Disabled = true
	disabledAt := time.Now()
	subscription.DisabledAt = &disabledAt
	subscription.UpdatedAt = time.Now()

	err = s.subscriptionsStore.Update(ctx, subscription)
	if err != nil {
		return nil, err
	}

	log.Println("subscription disabled: ", subscription.ID.Hex())

	return subscription.Out(), err
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

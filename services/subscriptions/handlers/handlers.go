package handlers

import (
	"context"
	"encoding/json"
	"go-subscriptions-workflow/rmq"
	"go-subscriptions-workflow/services/subscriptions/service"
	"go-subscriptions-workflow/types"
)

type subscriptionsHandlers struct {
	svc service.SubscriptionsServiceServer
}

func Register(svc service.SubscriptionsServiceServer, consumer rmq.Consumer) {
	h := &subscriptionsHandlers{svc: svc}
	consumer.HandleFunc(rmq.NewHandleMessageType(&types.StartSubscriptionRequest{}), h.HandleStartSubscription)
	consumer.HandleFunc(rmq.NewHandleMessageType(&types.CancelSubscriptionRequest{}), h.HandleCancelSubscription)
}

func (h *subscriptionsHandlers) HandleStartSubscription(ctx context.Context, data []byte) error {
	var req types.StartSubscriptionRequest
	err := json.Unmarshal(data, &req)
	if err != nil {
		return err
	}
	return h.svc.Start(ctx, &req)
}

func (h *subscriptionsHandlers) HandleCancelSubscription(ctx context.Context, data []byte) error {
	var req types.CancelSubscriptionRequest
	err := json.Unmarshal(data, &req)
	if err != nil {
		return err
	}
	return h.svc.Cancel(ctx, &req)
}
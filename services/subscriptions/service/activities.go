package service

import (
	"context"
	"go-subscriptions-workflow/types"
)

type Activities struct {
	svc SubscriptionsServiceServer
}

func (a *Activities) Charge(ctx context.Context, state SubscriptionState) (SubscriptionState, error) {
	out, err := a.svc.Charge(ctx, &types.ChargeSubscriptionRequest{ID: state.ID})
	if err != nil {
		return state, HandleError(err)
	}
	return NewState(out), nil
}

func (a *Activities) Disable(ctx context.Context, state SubscriptionState) (SubscriptionState, error) {
	out, err := a.svc.Disable(ctx, &types.DisableSubscriptionRequest{ID: state.ID})
	if err != nil {
		return state, HandleError(err)
	}
	return NewState(out), nil
}


package service

import (
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"strings"
	"time"
)

func SubscriptionsWorkflow(ctx workflow.Context, state SubscriptionState, activities *Activities) (SubscriptionState, error) {

	logger := workflow.GetLogger(ctx)

	logger.Debug("subscription workflow started.", "id", state.ID)

	err := workflow.SetQueryHandler(ctx, QuerySubscriptionState, func() (SubscriptionState, error) {
		return state, nil
	})
	if err != nil {
		return state, err
	}

	cancelSelector := workflow.NewSelector(ctx)
	cancelChannel := workflow.GetSignalChannel(ctx, SignalCancelSubscription)
	cancelSelector.AddReceive(cancelChannel, func(ch workflow.ReceiveChannel, _ bool) {
		var cancelSignal bool
		ch.Receive(ctx, &cancelSignal)
		state.Canceled = true
		state.CreatedAt = time.Now()
	})

	ao := workflow.ActivityOptions{
		StartToCloseTimeout:    time.Second * 10,
		RetryPolicy:            &temporal.RetryPolicy{
			MaximumAttempts:        3,
			NonRetryableErrorTypes: []string{ErrInsufficientFunds.Error()},
		},
	}

	ctx = workflow.WithActivityOptions(ctx, ao)

	for {
		_, err = workflow.AwaitWithTimeout(ctx, state.Expiration(), func() bool {
			return state.Canceled || state.Disabled
		})

		logger.Debug("subscription expired", "id", state.ID, "expires_at", state.ExpiresAt.String())


		err = workflow.ExecuteActivity(ctx, activities.Charge, state).Get(ctx, &state)

		if err != nil {
			if strings.Contains(err.Error(), ErrInsufficientFunds.Error()) {
				break
			}
			return state, err
		}

		for cancelSelector.HasPending() {
			cancelSelector.Select(ctx)
		}

		if state.Canceled {
			break
		}
	}

	if !state.Canceled {
		err = workflow.ExecuteActivity(ctx, activities.Disable, state).Get(ctx, &state)
		if err != nil {
			return state, err
		}
	}

	logger.Debug("subscription workflow finished.", "id", state.ID)

	return state, nil
}
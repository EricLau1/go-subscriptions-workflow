package service

import (
	"go-subscriptions-workflow/services/subscriptions/shared"
	"go.temporal.io/sdk/temporal"
)

var ErrInsufficientFunds = temporal.NewNonRetryableApplicationError(shared.ErrInsufficientFunds.Error(), "user_poor", shared.ErrInsufficientFunds, nil)

func HandleError(err error) error {
	switch err.Error() {
	case shared.ErrInsufficientFunds.Error():
		return ErrInsufficientFunds
	default:
		return err
	}
}

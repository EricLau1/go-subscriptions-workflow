package types

import (
	"time"
)

type CreateUserInput struct {
	Email    string `json:"email" validate:"email,required"`
	Password string `json:"password" validate:"required"`
}

type UserOutput struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LoginInput struct {
	CreateUserInput
}

type LoginOutput struct {
	*UserOutput
	Token string `json:"token"`
}

type CreditInput struct {
	UserID string  `json:"user_id"`
	Amount float64 `json:"amount" validate:"gte=0"`
}

type DebitInput struct {
	CreditInput
}

type StartSubscriptionRequest struct {
	UserID string `json:"id"`
}

type SubscriptionOutput struct {
	ID          string           `json:"id"`
	UserID      string           `json:"user_id"`
	Price       float64          `json:"price"`
	Features    []*FeatureOutput `json:"features"`
	Activations int              `json:"activations"`
	ActivatedAt time.Time        `json:"activated_at"`
	ExpiresAt   time.Time        `json:"expires_at"`
	Canceled    bool             `json:"canceled"`
	CanceledAt  *time.Time       `json:"canceled_at"`
	Disabled    bool             `json:"disabled"`
	DisabledAt  *time.Time       `json:"disabled_at"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

type FeatureOutput struct {
	Name string `json:"name"`
}

type ChargeSubscriptionRequest struct {
	ID string `json:"id"`
}

type CancelSubscriptionRequest struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
}

type DisableSubscriptionRequest struct {
	ID string `json:"id"`
}

package models

import (
	"go-subscriptions-workflow/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Feature struct {
	Name string
}

func (f *Feature) Out() *types.FeatureOutput {
	return &types.FeatureOutput{Name: f.Name}
}

type Subscription struct {
	ID          primitive.ObjectID `bson:"_id"`
	UserID      primitive.ObjectID `bson:"user_id"`
	Price       float64            `bson:"price"`
	Features    []*Feature         `bson:"features"`
	Activations int                `bson:"activations"`
	ActivatedAt time.Time          `bson:"activated_at"`
	ExpiresAt   time.Time          `bson:"expires_at"`
	Canceled    bool               `bson:"canceled"`
	CanceledAt  *time.Time         `bson:"canceled_at"`
	Disabled    bool               `bson:"disabled"`
	DisabledAt  *time.Time         `bson:"disabled_at"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}

func (s *Subscription) Out() *types.SubscriptionOutput {
	out := &types.SubscriptionOutput{
		ID:          s.ID.Hex(),
		UserID:      s.UserID.Hex(),
		Price:       s.Price,
		Activations: s.Activations,
		ActivatedAt: s.ActivatedAt,
		ExpiresAt:   s.ExpiresAt,
		Canceled:    s.Canceled,
		CanceledAt:  s.CanceledAt,
		Disabled:    s.Disabled,
		DisabledAt:  s.DisabledAt,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
	out.Features = make([]*types.FeatureOutput, 0, len(s.Features))
	for index := range s.Features {
		out.Features = append(out.Features, s.Features[index].Out())
	}
	return out
}

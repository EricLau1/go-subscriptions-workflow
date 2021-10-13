package models

import (
	"go-subscriptions-workflow/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id"`
	Email     string             `bson:"email"`
	Password  string             `bson:"password"`
	Balance   float64            `bson:"balance"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

func (u *User) Out() *types.UserOutput {
	return &types.UserOutput{
		ID:        u.ID.Hex(),
		Email:     u.Email,
		Password:  u.Password,
		Balance:   u.Balance,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

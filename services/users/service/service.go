package service

import (
	"context"
	"fmt"
	"go-subscriptions-workflow/db"
	"go-subscriptions-workflow/security/passwords"
	"go-subscriptions-workflow/security/tokens"
	"go-subscriptions-workflow/services/users/models"
	"go-subscriptions-workflow/services/users/store"
	"go-subscriptions-workflow/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type UsersService interface {
	CreateUser(ctx context.Context, in *types.CreateUserInput) (*types.UserOutput, error)
	Login(ctx context.Context, in *types.LoginInput) (*types.LoginOutput, error)
	GetUser(ctx context.Context, id string) (*types.UserOutput, error)
	Credit(ctx context.Context, in *types.CreditInput) (*types.UserOutput, error)
	Debit(ctx context.Context, in *types.DebitInput) (*types.UserOutput, error)
}

type usersService struct {
	usersStore store.UsersStore
}

func NewUsersService(dbConn db.Connection) UsersService {
	return &usersService{usersStore: store.NewUsersStore(dbConn.DB())}
}

func (s *usersService) CreateUser(ctx context.Context, in *types.CreateUserInput) (*types.UserOutput, error) {
	_, err := s.usersStore.GetByEmail(ctx, in.Email)
	if err == mongo.ErrNoDocuments {
		password, err := passwords.New(in.Password)
		if err != nil {
			return nil, err
		}
		user := &models.User{
			ID:        primitive.NewObjectID(),
			Email:     in.Email,
			Password:  password,
			Balance:   0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err = s.usersStore.Create(ctx, user)
		if err != nil {
			return nil, err
		}
		return user.Out(), nil
	}
	if err == nil {
		return nil, fmt.Errorf("email already registered: %s", in.Email)
	}
	return nil, err
}

func (s *usersService) Login(ctx context.Context, in *types.LoginInput) (*types.LoginOutput, error) {
	user, err := s.usersStore.GetByEmail(ctx, in.Email)
	if err != nil {
		return nil, err
	}
	err = passwords.OK(user.Password, in.Password)
	if err != nil {
		return nil, err
	}
	token, err := tokens.New(user.ID.Hex())
	if err != nil {
		return nil, err
	}
	out := new(types.LoginOutput)
	out.UserOutput = user.Out()
	out.Token = token
	return out, nil
}

func (s *usersService) GetUser(ctx context.Context, id string) (*types.UserOutput, error) {
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	user, err := s.usersStore.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user.Out(), nil
}

func (s *usersService) Credit(ctx context.Context, in *types.CreditInput) (*types.UserOutput, error) {
	userID, err := primitive.ObjectIDFromHex(in.UserID)
	if err != nil {
		return nil, err
	}
	if in.Amount <= 0 {
		return nil, fmt.Errorf("invalid amount to credit: %f", in.Amount)
	}
	user, err := s.usersStore.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	user.Balance += in.Amount
	user.UpdatedAt = time.Now()
	err = s.usersStore.Update(ctx, user)
	if err != nil {
		return nil, err
	}
	return user.Out(), nil
}

func (s *usersService) Debit(ctx context.Context, in *types.DebitInput) (*types.UserOutput, error) {
	userID, err := primitive.ObjectIDFromHex(in.UserID)
	if err != nil {
		return nil, err
	}
	if in.Amount <= 0 {
		return nil, fmt.Errorf("invalid amount to debit: %f", in.Amount)
	}
	user, err := s.usersStore.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	user.Balance -= in.Amount
	user.UpdatedAt = time.Now()
	err = s.usersStore.Update(ctx, user)
	if err != nil {
		return nil, err
	}
	return user.Out(), nil
}

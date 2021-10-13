package store

import (
	"context"
	"go-subscriptions-workflow/services/subscriptions/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type SubscriptionsStore interface {
	Create(ctx context.Context, subscription *models.Subscription) error
	Update(ctx context.Context, subscription *models.Subscription) error
	Get(ctx context.Context, id primitive.ObjectID) (*models.Subscription, error)
	GetAll(ctx context.Context) ([]*models.Subscription, error)
	GetByUserID(ctx context.Context, userID primitive.ObjectID) ([]*models.Subscription, error)
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type subscriptionsStore struct {
	coll *mongo.Collection
}

func NewSubscriptionsStore(dbConn *mongo.Database) SubscriptionsStore {
	return &subscriptionsStore{coll: dbConn.Collection("subscriptions")}
}

func (s *subscriptionsStore) Create(ctx context.Context, subscription *models.Subscription) error {
	result, err := s.coll.InsertOne(ctx, subscription)
	if err != nil {
		return err
	}
	log.Printf("subscription created: %+v\n", result)
	return nil
}

func (s *subscriptionsStore) Update(ctx context.Context, subscription *models.Subscription) error {

	update := bson.M{
		"$set": bson.M{
			"price":        subscription.Price,
			"features":     subscription.Features,
			"activations":  subscription.Activations,
			"activated_at": subscription.ActivatedAt,
			"expires_at":   subscription.ExpiresAt,
			"canceled":     subscription.Canceled,
			"canceled_at":  subscription.CanceledAt,
			"disabled":     subscription.Disabled,
			"disabled_at":  subscription.DisabledAt,
			"updated_at":   subscription.UpdatedAt,
		},
	}

	result, err := s.coll.UpdateByID(ctx, subscription.ID, update)
	if err != nil {
		return err
	}
	log.Printf("subscription updated: %+v\n", result)
	return nil
}

func (s *subscriptionsStore) Get(ctx context.Context, id primitive.ObjectID) (*models.Subscription, error) {
	var subscription models.Subscription
	err := s.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&subscription)
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

func (s *subscriptionsStore) GetAll(ctx context.Context) ([]*models.Subscription, error) {
	cursor, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var subscriptions []*models.Subscription
	err = cursor.All(ctx, &subscriptions)
	if err != nil {
		return nil, err
	}
	return subscriptions, nil
}

func (s *subscriptionsStore) GetByUserID(ctx context.Context, userID primitive.ObjectID) ([]*models.Subscription, error) {
	cursor, err := s.coll.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var subscriptions []*models.Subscription
	err = cursor.All(ctx, &subscriptions)
	if err != nil {
		return nil, err
	}
	return subscriptions, nil
}

func (s *subscriptionsStore) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := s.coll.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	log.Printf("subscription deleted: %+v\n", result)
	return nil
}

package store

import (
	"context"
	"go-subscriptions-workflow/services/users/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type UsersStore interface {
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
	Get(ctx context.Context, id primitive.ObjectID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetAll(ctx context.Context) ([]*models.User, error)
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type usersStore struct {
	coll *mongo.Collection
}

func NewUsersStore(dbConn *mongo.Database) UsersStore {
	return &usersStore{coll: dbConn.Collection("users")}
}

func (s *usersStore) Create(ctx context.Context, user *models.User) error {
	result, err := s.coll.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	log.Printf("user created: %+v\n", result)
	return nil
}

func (s *usersStore) Update(ctx context.Context, user *models.User) error {

	update := bson.M{
		"$set": bson.M{
			"email":      user.Email,
			"password":   user.Password,
			"balance":    user.Balance,
			"updated_at": user.UpdatedAt,
		},
	}

	result, err := s.coll.UpdateByID(ctx, user.ID, update)
	if err != nil {
		return err
	}
	log.Printf("user updated: %+v\n", result)
	return nil
}

func (s *usersStore) Get(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	var user models.User
	err := s.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *usersStore) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := s.coll.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *usersStore) GetAll(ctx context.Context) ([]*models.User, error) {
	cursor, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var users []*models.User
	err = cursor.All(ctx, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *usersStore) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := s.coll.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return err
	}
	log.Printf("user deleted: %+v\n", result)
	return nil
}
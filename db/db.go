package db

import (
	"context"
	"go-subscriptions-workflow/util"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
)

type Connection interface {
	Ping(ctx context.Context) error
	Close(ctx context.Context)
	DB() *mongo.Database
}

type connection struct {
	mongoClient *mongo.Client
}

func New(ctx context.Context, cfg Config) Connection {
	var conn connection
	var err error
	opts := options.Client().ApplyURI(cfg.URI())
	conn.mongoClient, err = mongo.Connect(ctx, opts)
	util.PanicOnError(err)
	return &conn
}

func (c *connection) Ping(ctx context.Context) error {
	return c.mongoClient.Ping(ctx, readpref.Primary())
}

func (c *connection) Close(ctx context.Context) {
	err := c.mongoClient.Disconnect(ctx)
	if err != nil {
		log.Println("error on disconnect database:", err)
	}
}

func (c *connection) DB() *mongo.Database {
	return c.mongoClient.Database(database)
}

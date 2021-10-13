package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"go-subscriptions-workflow/db"
	"go-subscriptions-workflow/rmq"
	"go-subscriptions-workflow/services/subscriptions/handlers"
	"go-subscriptions-workflow/services/subscriptions/service"
	"go-subscriptions-workflow/services/subscriptions/shared"
	userssvc "go-subscriptions-workflow/services/users/service"
	"go-subscriptions-workflow/util"
	"log"
	"time"
)

func init() {
	err := godotenv.Load(util.GetEnvFilePath())
	util.PanicOnError(err)
	db.LoadConfigFromEnv()
	rmq.LoadConfigFromEnv()
}

func main() {
	dbCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	dbConn := db.New(dbCtx, db.NewConfig())
	defer dbConn.Close(dbCtx)

	err := dbConn.Ping(dbCtx)
	util.PanicOnError(err)
	log.Println("mongodb connected!")

	rmqConn := rmq.New(rmq.NewConfig())
	defer rmqConn.Close()
	log.Println("rabbitmq connected!")

	err = rmqConn.ExchangeDeclare(&rmq.ExchangeOptions{
		Name:    shared.ExchangeName,
		Kind:    amqp.ExchangeDirect,
		Durable: true,
	})
	util.PanicOnError(err)
	log.Println("workflows exchange declared!")

	err = rmqConn.QueueDeclare(&rmq.QueueOptions{
		Name:    shared.QueueName,
		Durable: true,
		BindOptions: &rmq.QueueBindOptions{
			ExchangeName: shared.ExchangeName,
		},
	})
	util.PanicOnError(err)
	log.Println("subscriptions queue declared!")

	consumer := rmqConn.NewConsumer()

	usersService := userssvc.NewUsersService(dbConn)
	subscriptionsService := service.NewSubscriptionsServiceServer(dbConn, usersService)
	handlers.Register(subscriptionsService, consumer)

	log.Println("subscriptions service is running...")

	err = consumer.Listen(context.Background(), &rmq.ConsumerOptions{
		QueueName: shared.QueueName,
	})
	util.PanicOnError(err)
}

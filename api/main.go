package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"go-subscriptions-workflow/api/handlers"
	"go-subscriptions-workflow/db"
	"go-subscriptions-workflow/rmq"
	subssvc "go-subscriptions-workflow/services/subscriptions/service"
	userssvc "go-subscriptions-workflow/services/users/service"
	"go-subscriptions-workflow/util"
	"go.temporal.io/sdk/client"
	"log"
	"time"
)

var port int

func init() {
	err := godotenv.Load(util.GetEnvFilePath())
	util.PanicOnError(err)
	db.LoadConfigFromEnv()
	rmq.LoadConfigFromEnv()
	flag.IntVar(&port, "port", 8080, "api port")
	flag.Parse()
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

	temporalClient, err := client.NewClient(client.Options{})
	util.PanicOnError(err)
	defer temporalClient.Close()
	log.Println("temporal client connected!")

	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())

	usersService := userssvc.NewUsersService(dbConn)
	handlers.RegisterUsersHandlers(usersService, app)
	subsClient := subssvc.NewSubscriptionsClient(dbConn, usersService, temporalClient)
	handlers.RegisterSubscriptionsHandlers(subsClient, rmqConn.NewProducer(), app)

	err = app.Listen(fmt.Sprintf(":%d", port))
	util.PanicOnError(err)
}

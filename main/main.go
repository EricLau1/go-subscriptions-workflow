package main

import (
	"context"
	"github.com/joho/godotenv"
	"go-subscriptions-workflow/db"
	"go-subscriptions-workflow/util"
	"log"
	"time"
)

func init() {
	err := godotenv.Load(util.GetEnvFilePath())
	util.PanicOnError(err)
	db.LoadConfigFromEnv()
}

func main() {
	dbCtx, cancel := context.WithTimeout(context.Background(), time.Second * 10)
	defer cancel()

	dbConn := db.New(dbCtx, db.NewConfig())
	defer dbConn.Close(dbCtx)

	err := dbConn.Ping(dbCtx)
	util.PanicOnError(err)
	log.Println("mongodb connected!")
}

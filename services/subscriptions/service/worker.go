package service

import (
	"go-subscriptions-workflow/util"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"log"
)

func Worker(temporalClient client.Client, svc SubscriptionsServiceServer) {
	log.Println("subscriptions worker starting")
	w := worker.New(temporalClient, TaskQueueName, worker.Options{})
	w.RegisterWorkflow(SubscriptionsWorkflow)
	w.RegisterActivity(&Activities{svc: svc})
	util.PanicOnError(w.Run(worker.InterruptCh()))
}
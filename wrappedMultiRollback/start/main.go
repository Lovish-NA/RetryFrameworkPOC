package main

import (
	"context"
	"log"
	app2 "money-transfer-project-template-go/app/wrappedMultiRollback"

	"go.temporal.io/sdk/client"
)

// @@@SNIPSTART money-transfer-project-template-go-start-workflow
func main() {
	// Create the client object just once per process
	c, err := client.Dial(client.Options{})

	if err != nil {
		log.Fatalln("Unable to create Temporal client:", err)
	}

	defer c.Close()

	options := client.StartWorkflowOptions{
		ID:        "pay-invoice-701",
		TaskQueue: app2.OrderTaskQueueName,
	}

	we, err := c.ExecuteWorkflow(context.Background(), options, app2.OrderWorkflow)

	if err != nil {
		// err
		log.Fatalln("Unable to start the Workflow:", err)
	}

	log.Printf("WorkflowID: %s RunID: %s\n", we.GetID(), we.GetRunID())

	var result string

	err = we.Get(context.Background(), &result)

	if err != nil {
		log.Fatalln("Unable to get Workflow result:", err)
	}

	log.Println(result)
}

// @@@SNIPEND

package main

import (
	"context"
	"log"
	"money-transfer-project-template-go/app/Main"
	app2 "money-transfer-project-template-go/app/sagaMain/workflows"

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
		TaskQueue: "money-transfer-task-queue",
	}

	input := Main.PaymentDetails{
		SourceAccount: "85-150",
		TargetAccount: "43-812",
		Amount:        250,
		ReferenceID:   "12345",
	}

	we, err := c.ExecuteWorkflow(context.Background(), options, app2.MoneyTransfer, input)

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

package main

import (
	"context"
	"log"
	"money-transfer-project-template-go/app/rollbackWithSaga"
	"money-transfer-project-template-go/app/rollbackWithSaga/activities"
	app2 "money-transfer-project-template-go/app/rollbackWithSaga/workflows"

	"go.temporal.io/sdk/client"
)

func main() {
	// Create the Temporal client
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create Temporal client:", err)
	}
	defer c.Close()

	// Define workflow input
	input := activities.OrderDetails{
		OrderID: "701",
		Amount:  199.99,
		Address: "221B Baker Street",
	}

	// Define workflow start options
	options := client.StartWorkflowOptions{
		ID:        "order-workflow-701",
		TaskQueue: rollbackWithSaga.OrderTaskQueueName, // Define this constant in your saga package
	}

	// Start the workflow
	we, err := c.ExecuteWorkflow(context.Background(), options, app2.OrderWorkflow, input)
	if err != nil {
		log.Fatalln("Unable to start the Workflow:", err)
	}

	log.Printf("Started Workflow. WorkflowID: %s RunID: %s\n", we.GetID(), we.GetRunID())

	var result string
	err = we.Get(context.Background(), &result)
	if err != nil {
		log.Fatalln("Unable to get Workflow result:", err)
	}

	log.Println("Workflow completed with result:", result)
}

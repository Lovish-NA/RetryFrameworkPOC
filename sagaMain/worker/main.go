package main

import (
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"money-transfer-project-template-go/app/sagaMain/activities"
	"money-transfer-project-template-go/app/sagaMain/workflows"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		panic(err)
	}
	defer c.Close()

	w := worker.New(c, "money-transfer-task-queue", worker.Options{})

	w.RegisterWorkflow(workflows.MoneyTransfer)
	w.RegisterActivity(activities.Withdraw)
	w.RegisterActivity(activities.Deposit)
	w.RegisterActivity(activities.Refund)

	if err := w.Run(worker.InterruptCh()); err != nil {
		panic(err)
	}
}

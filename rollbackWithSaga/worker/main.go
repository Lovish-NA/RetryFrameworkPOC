package main

import (
	"log"
	"money-transfer-project-template-go/app/rollbackWithSaga"
	activities "money-transfer-project-template-go/app/rollbackWithSaga/activities"
	app "money-transfer-project-template-go/app/rollbackWithSaga/workflows"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

// @@@SNIPSTART money-transfer-project-template-go-worker
func main() {

	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, rollbackWithSaga.OrderTaskQueueName, worker.Options{})
	w.RegisterWorkflow(app.OrderWorkflow)
	w.RegisterActivity(activities.ProcessOrder)
	w.RegisterActivity(activities.CancelOrder)
	w.RegisterActivity(activities.ChargePayment)
	w.RegisterActivity(activities.RefundPayment)
	w.RegisterActivity(activities.ShipOrder)
	w.RegisterActivity(activities.CancelShipment)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}

// @@@SNIPEND

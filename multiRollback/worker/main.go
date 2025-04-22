package main

import (
	"log"
	app "money-transfer-project-template-go/app/multiRollback"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	_ "money-transfer-project-template-go/app/multiRollback"
)

// @@@SNIPSTART money-transfer-project-template-go-worker
func main() {

	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "ORDER_TASK_QUEUE", worker.Options{})
	w.RegisterWorkflow(app.OrderWorkflow)
	w.RegisterActivity(app.ProcessOrder)
	w.RegisterActivity(app.CancelOrder)
	w.RegisterActivity(app.ChargePayment)
	w.RegisterActivity(app.RefundPayment)
	w.RegisterActivity(app.ShipOrder)
	w.RegisterActivity(app.CancelShipment)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}

// @@@SNIPEND

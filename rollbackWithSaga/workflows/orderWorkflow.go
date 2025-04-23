package workflows

import (
	"go.temporal.io/sdk/workflow"

	"money-transfer-project-template-go/app/rollbackWithSaga/activities"
	"money-transfer-project-template-go/app/rollbackWithSaga/saga"
)

func OrderWorkflow(ctx workflow.Context, input activities.OrderDetails) (string, error) {
	steps := []saga.WorkflowStep{
		{
			Name:             "ProcessOrder",
			Activity:         activities.ProcessOrder,
			Args:             []interface{}{input},
			FallbackActivity: activities.CancelOrder,
			RetryThreshold:   3,
		},
		{
			Name:             "ChargePayment",
			Activity:         activities.ChargePayment,
			Args:             []interface{}{input},
			FallbackActivity: activities.RefundPayment,
			RetryThreshold:   3,
		},
		{
			Name:             "ShipOrder",
			Activity:         activities.ShipOrder,
			Args:             []interface{}{input},
			FallbackActivity: activities.CancelShipment,
			RetryThreshold:   3,
		},
	}

	err := saga.ExecuteSagaWithRollback(ctx, steps)
	if err != nil {
		return "", err
	}

	return "Order completed", nil
}

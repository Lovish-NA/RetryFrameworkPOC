package workflows

import (
	"go.temporal.io/sdk/workflow"

	"money-transfer-project-template-go/app/sagaMain/activities"
	"money-transfer-project-template-go/app/sagaMain/saga"
)

func MoneyTransfer(ctx workflow.Context, input activities.PaymentDetails) (string, error) {
	steps := []saga.WorkflowStep{
		{
			Name:             "Withdraw",
			Activity:         activities.Withdraw,
			Args:             []interface{}{input},
			FallbackActivity: activities.Refund,
			RetryThreshold:   3,
		},
		{
			Name:             "Deposit",
			Activity:         activities.Deposit,
			Args:             []interface{}{input},
			FallbackActivity: activities.Refund,
			RetryThreshold:   3,
		},
	}

	err := saga.ExecuteSagaWithRollback(ctx, steps)
	if err != nil {
		return "", err
	}

	return "Transfer complete", nil
}

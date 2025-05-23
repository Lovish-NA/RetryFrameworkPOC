package app

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// @@@SNIPSTART OrderWorkflow-project-template-go-workflow
func OrderWorkflow(ctx workflow.Context) (string, error) {

	// RetryPolicy specifies how to automatically handle retries if an Activity fails.
	retrypolicy := &temporal.RetryPolicy{
		InitialInterval:        time.Second,
		BackoffCoefficient:     1.0,
		MaximumInterval:        100 * time.Second,
		MaximumAttempts:        5, // 0 is unlimited retries
		NonRetryableErrorTypes: []string{"InvalidAccountError", "InsufficientFundsError"},
	}

	options := workflow.ActivityOptions{
		// Timeout options specify when to automatically timeout Activity functions.
		StartToCloseTimeout: time.Minute,
		// Optionally provide a customized RetryPolicy.
		// Temporal retries failed Activities by default.
		RetryPolicy: retrypolicy,
	}

	// Apply the options.
	ctx = workflow.WithActivityOptions(ctx, options)

	err := workflow.ExecuteActivity(ctx, ProcessOrder).Get(ctx, nil)
	if err != nil {
		// Invoke rollback workflow
		workflow.ExecuteChildWorkflow(ctx, RollbackWorkflow, "ProcessOrder").Get(ctx, nil)
		return "", err
	}

	err = workflow.ExecuteActivity(ctx, ChargePayment).Get(ctx, nil)
	if err != nil {
		// Invoke rollback workflow
		workflow.ExecuteChildWorkflow(ctx, RollbackWorkflow, "ChargePayment").Get(ctx, nil)
		return "", err
	}

	err = workflow.ExecuteActivity(ctx, ShipOrder).Get(ctx, nil)
	if err != nil {
		// Invoke rollback workflow
		workflow.ExecuteChildWorkflow(ctx, RollbackWorkflow, "ShipOrder").Get(ctx, nil)
		return "", err
	}

	return "", nil
}

// @@@SNIPEND

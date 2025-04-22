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

	rm := NewRollbackManager(ctx)
	ae := NewActivityExecutor(ctx, rm)

	// Step 1: Process Order
	err := ae.ExecuteWithRollback(ProcessOrder, CancelOrder)
	if err != nil {
		rm.ExecuteRollback()
		return "", err
	}

	// Step 2: Charge Payment
	err = ae.ExecuteWithRollback(ChargePayment, RefundPayment)
	if err != nil {
		rm.ExecuteRollback()
		return "", err
	}

	// Step 3: Ship Order
	err = ae.ExecuteWithRollback(ShipOrder, CancelShipment)
	if err != nil {
		rm.ExecuteRollback()
		return "", err
	}

	return "", nil
}

// @@@SNIPEND

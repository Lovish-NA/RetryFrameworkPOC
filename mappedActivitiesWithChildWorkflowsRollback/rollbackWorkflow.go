package app

import (
	"fmt"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"time"
)

// RollbackWorkflow is the workflow that performs rollbacks based on the mappings
func RollbackWorkflow(ctx workflow.Context, failedActivity string) error {

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

	mappings, err := LoadMappings("mappedActivitiesWithChildWorkflowsRollback/activities.json")
	// Print the loaded mappings for debugging
	fmt.Println("Loaded mappings:", mappings)
	if err != nil {
		return err
	}

	rollbackStack := []string{}
	for activity, rollbackActivity := range mappings.Mappings {
		rollbackStack = append(rollbackStack, rollbackActivity)
		fmt.Println("Added to rollback stack", "activity", activity, "rollbackActivity", rollbackActivity)
		if activity == failedActivity {
			break
		}
	}

	for i := len(rollbackStack) - 1; i >= 0; i-- {
		rollbackActivity := rollbackStack[i]
		err := workflow.ExecuteActivity(ctx, rollbackActivity).Get(ctx, nil)
		if err != nil {
			workflow.GetLogger(ctx).Error("Rollback failed", "activity", rollbackActivity, "error", err)
		}
	}

	return nil
}

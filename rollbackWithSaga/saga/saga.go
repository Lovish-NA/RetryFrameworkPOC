package saga

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type WorkflowStep struct {
	Name             string
	Activity         interface{}
	Args             []interface{}
	FallbackActivity interface{}
	RetryThreshold   int
	NonRetriableErrs []string
}

func ExecuteSagaWithRollback(ctx workflow.Context, steps []WorkflowStep) error {
	logger := workflow.GetLogger(ctx)
	var executedSteps []WorkflowStep

	retryPolicy := &temporal.RetryPolicy{
		InitialInterval:        time.Second,
		BackoffCoefficient:     1.0,
		MaximumInterval:        10 * time.Second,
		MaximumAttempts:        3,
		NonRetryableErrorTypes: []string{"invalidOrder"},
	}

	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
		RetryPolicy:         retryPolicy,
	}

	stepCtx := workflow.WithActivityOptions(ctx, activityOptions)
	for _, step := range steps {
		logger.Info("Executing step", "step", step.Name)
		err := workflow.ExecuteActivity(stepCtx, step.Activity, step.Args...).Get(stepCtx, nil)
		if err != nil {
			logger.Error("Step failed", "step", step.Name, "error", err)

			// Attempt fallback for current step
			if step.FallbackActivity != nil {
				logger.Info("Running fallback", "step", step.Name)
				fallbackErr := workflow.ExecuteActivity(stepCtx, step.FallbackActivity, step.Args...).Get(stepCtx, nil)
				if fallbackErr != nil {
					logger.Error("Fallback activity failed", "step", step.Name, "error", fallbackErr)
				}
			}

			// Rollback previously executed steps in reverse order
			for i := len(executedSteps) - 1; i >= 0; i-- {
				prevStep := executedSteps[i]
				if prevStep.FallbackActivity != nil {
					logger.Info("Rolling back previous step", "step", prevStep.Name)
					rollbackErr := workflow.ExecuteActivity(stepCtx, prevStep.FallbackActivity, prevStep.Args...).Get(stepCtx, nil)
					if rollbackErr != nil {
						logger.Error("Rollback activity failed", "step", prevStep.Name, "error", rollbackErr)
					}
				}
			}

			return fmt.Errorf("step %s failed: %w", step.Name, err)
		}

		executedSteps = append(executedSteps, step)
	}
	return nil
}

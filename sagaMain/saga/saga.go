package saga

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func ExecuteSagaWithRollback(ctx workflow.Context, steps []WorkflowStep) error {
	var executed []WorkflowStep

	for _, step := range steps {
		policy := &temporal.RetryPolicy{
			InitialInterval:        time.Second,
			BackoffCoefficient:     2.0,
			MaximumInterval:        10 * time.Second,
			MaximumAttempts:        int32(step.RetryThreshold),
			NonRetryableErrorTypes: step.NonRetriableErrs,
		}

		opts := workflow.ActivityOptions{
			StartToCloseTimeout: time.Minute,
			RetryPolicy:         policy,
		}

		stepCtx := workflow.WithActivityOptions(ctx, opts)

		err := workflow.ExecuteActivity(stepCtx, step.Activity, step.Args...).Get(stepCtx, nil)
		if err != nil {
			workflow.GetLogger(ctx).Error("Step failed", "step", step.Name, "error", err)

			if step.FallbackActivity != nil {
				_ = runFallback(ctx, step)
			}

			for i := len(executed) - 1; i >= 0; i-- {
				if executed[i].FallbackActivity != nil {
					_ = runFallback(ctx, executed[i])
				}
			}
			return fmt.Errorf("step %s failed: %w", step.Name, err)
		}

		executed = append(executed, step)
	}
	return nil
}

func runFallback(ctx workflow.Context, step WorkflowStep) error {
	return workflow.ExecuteActivity(ctx, step.FallbackActivity, step.Args...).Get(ctx, nil)
}

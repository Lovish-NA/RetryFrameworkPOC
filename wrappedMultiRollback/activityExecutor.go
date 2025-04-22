package app

import "go.temporal.io/sdk/workflow"

// ActivityExecutor wraps the logic for executing activities and registering rollbacks
type ActivityExecutor struct {
	ctx workflow.Context
	rm  *RollbackManager
}

func NewActivityExecutor(ctx workflow.Context, rm *RollbackManager) *ActivityExecutor {
	return &ActivityExecutor{
		ctx: ctx,
		rm:  rm,
	}
}

func (ae *ActivityExecutor) ExecuteWithRollback(activity, rollbackActivity interface{}, args ...interface{}) error {
	// Register the rollback activity first
	ae.rm.Add(rollbackActivity, args...)
	// Execute the main activity
	return workflow.ExecuteActivity(ae.ctx, activity, args...).Get(ae.ctx, nil)
}

package app

import (
	"go.temporal.io/sdk/workflow"
)

type RollbackManager struct {
	ctx       workflow.Context
	rollbacks []func() error
}

func NewRollbackManager(ctx workflow.Context) *RollbackManager {
	return &RollbackManager{
		ctx:       ctx,
		rollbacks: []func() error{},
	}
}

func (rm *RollbackManager) Add(rollbackActivity interface{}, args ...interface{}) {
	rm.rollbacks = append(rm.rollbacks, func() error {
		return workflow.ExecuteActivity(rm.ctx, rollbackActivity, args...).Get(rm.ctx, nil)
	})
}

func (rm *RollbackManager) ExecuteRollback(inParallel bool) {
	if !inParallel {
		// Compensate in Last-In-First-Out order, to undo in the reverse order that activities were applied.
		for i := len(rm.rollbacks) - 1; i >= 0; i-- {
			errCompensation := rm.rollbacks[i]()
			if errCompensation != nil {
				workflow.GetLogger(rm.ctx).Error("Executing compensation failed", "Error", errCompensation)
			}
		}
	} else {
		selector := workflow.NewSelector(rm.ctx)
		for i := 0; i < len(rm.rollbacks); i++ {
			execution := rm.rollbacks[i]
			selector.AddFuture(workflow.ExecuteActivity(rm.ctx, execution), func(f workflow.Future) {
				if errCompensation := f.Get(rm.ctx, nil); errCompensation != nil {
					workflow.GetLogger(rm.ctx).Error("Executing compensation failed", "Error", errCompensation)
				}
			})
		}
		for range rm.rollbacks {
			selector.Select(rm.ctx)
		}
	}
}

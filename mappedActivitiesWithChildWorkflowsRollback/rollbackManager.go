package app

import (
	"go.temporal.io/sdk/workflow"
)

type RollbackManager struct {
	ctx       workflow.Context
	rollbacks []func()
}

func NewRollbackManager(ctx workflow.Context) *RollbackManager {
	return &RollbackManager{
		ctx:       ctx,
		rollbacks: []func(){},
	}
}

func (rm *RollbackManager) Add(rollbackActivity interface{}, args ...interface{}) {
	rm.rollbacks = append(rm.rollbacks, func() {
		_ = workflow.ExecuteActivity(rm.ctx, rollbackActivity, args...).Get(rm.ctx, nil)
	})
}

func (rm *RollbackManager) ExecuteRollback() {
	for i := len(rm.rollbacks) - 1; i >= 0; i-- {
		rm.rollbacks[i]()
	}
}

package executor

import (
	"fmt"
)

var (
	executorInstance Executor
	EManager         executor_manager
)

func init() {
	executorInstance = NewExecutor()
}

type executor_manager struct {
}

func (s *executor_manager) Execute() error {
	fmt.Println("this is Execute call from Manager ")
	return nil
}

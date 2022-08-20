package judge

import (
	"fmt"

	"github.com/cranemont/judge-manager/judge/config"
)

type Runner interface {
	Run(task *Task, out chan<- string)
}

type runner struct {
	sandbox Sandbox
	option  *config.RunOption
}

func NewRunner(sandbox Sandbox, option *config.RunOption) *runner {
	return &runner{sandbox, option}
}

func (r *runner) Run(task *Task, out chan<- string) {
	fmt.Println("RUN! from runner")
	// r.sandbox.Execute()
	dir := task.GetDir()
	out <- "task " + dir + " running done"
	// 채널로 결과반환
}

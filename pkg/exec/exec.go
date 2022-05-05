package exec

import (
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
)

// Run 执行指令
func Run(item command.OperationItem) command.RunErr {
	executor, err := runner.ParseExecutor(item.B, item.Logger)
	if err != nil {
		return command.RunErr{Err: err}
	}
	return runner.RemoteRun(runner.RemoteRunItem{
		ManifestContent:     item.B,
		Logger:              item.Logger,
		Cmd:                 executor.Script,
		RecordErrServerList: true,
	})
}

func SURun(item command.OperationItem) command.RunErr {
	executor, err := runner.ParseExecutor(item.B, item.Logger)
	if err != nil {
		return command.RunErr{Err: err}
	}
	return runner.RemoteRun(runner.RemoteRunItem{
		ManifestContent:     item.B,
		Logger:              item.Logger,
		Cmd:                 executor.Script,
		RunShellFunc:        runner.RunOnNodeWithChangeToRoot,
		RecordErrServerList: true,
	})
}

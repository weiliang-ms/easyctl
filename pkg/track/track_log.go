package track

import (
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"gopkg.in/yaml.v2"
)

// TailLogExecutor 追踪日志执行器
type TailLogExecutor struct {
	TailLog struct {
		LogPath string `yaml:"log-path"`
		Whence  int    `yaml:"whence"`
		Offset  int64  `yaml:"offset"`
	} `yaml:"tail-log"`
}

// TaiLog 多级追踪日志
func TaiLog(item command.OperationItem) error {

	servers, err := runner.ParseServerList(item.B, item.Logger)
	stopCh := make(chan struct{})
	if err != nil {
		return err
	}

	executor, err := parseTailLogExecutor(item.B)
	if err != nil {
		return err
	}

	executor.Tail(servers, stopCh)

	return nil
}

// Tail tail日志
func (tail TailLogExecutor) Tail(servers []runner.ServerInternal, stopCh <-chan struct{}) {
	for _, v := range servers {
		go v.TailFile(tail.TailLog.LogPath, tail.TailLog.Offset, tail.TailLog.Whence, stopCh)
	}

	<-stopCh
}

func parseTailLogExecutor(b []byte) (TailLogExecutor, error) {
	executor := TailLogExecutor{}
	err := yaml.Unmarshal(b, &executor)
	if err != nil {
		return TailLogExecutor{}, err
	}

	return executor, nil
}

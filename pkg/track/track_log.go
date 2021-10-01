package track

import (
	"github.com/sirupsen/logrus"
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"gopkg.in/yaml.v2"
)

type TailLogExecutor struct {
	TailLog struct {
		LogPath string `yaml:"log-path"`
		Whence  int    `yaml:"whence"`
		Offset  int64  `yaml:"offset"`
	} `yaml:"tail-log"`
}

func TaiLog(config []byte, logger *logrus.Logger) error {

	servers, err := runner.ParseServerList(config)
	stopCh := make(chan struct{})
	if err != nil {
		return err
	}

	executor, err := parseTailLogExecutor(config)
	if err != nil {
		return err
	}

	executor.Tail(servers, stopCh)

	return nil
}

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

package install

import (
	"fmt"
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"k8s.io/klog"
)

type Task struct {
	Servers []runner.ServerInternal
	Cmd     string
	Debug   bool
	ErrChan chan error
}

type Executor struct {
	Servers []runner.ServerInternal
	Meta    interface{}
}

type Interface interface {
	Combine(servers []runner.ServerInternal) Executor
	Detect(executor Executor, debug bool) error
	HandPackage(executor Executor, debug bool) error
	Compile(executor Executor, debug bool) error
}

type TaskFnc func(Task) error

func Install(i Interface, b []byte, debug bool) error {

	server, err := ParseServerList(b, debug)
	if err != nil {
		return err
	}
	executor := i.Combine(server)
	if err != nil {
		return err
	}

	if err := i.Detect(executor, debug); err != nil {
		return err
	}
	if err := i.HandPackage(executor, debug); err != nil {
		return err
	}
	if err := i.Compile(executor, debug); err != nil {
		return err
	}

	return nil
}

func ParseServerList(b []byte, debug bool) ([]runner.ServerInternal, error) {

	klog.Infoln("解析主机列表...")
	servers, err := runner.ParseServerList(b)
	if err != nil {
		return []runner.ServerInternal{}, err
	}

	if debug {
		fmt.Printf("%+v", servers)
	}

	return servers, nil
}

package scan

import (
	"github.com/sirupsen/logrus"
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"strings"
)

//go:generate mockery --name=HandleOSInterface
type HandleOSInterface interface {
	DoRequest(r runner.RunItem) runner.ShellResult
	GetHostName(s runner.ServerInternal, l *logrus.Logger) (result string, err error)
	GetKernelVersion(s runner.ServerInternal, l *logrus.Logger) (result string, err error)
	GetSystemVersion(s runner.ServerInternal, l *logrus.Logger) (result string, err error)
	GetCPUInfo(s runner.ServerInternal, l *logrus.Logger) (result string, err error)
	GetCPULoadAverage(s runner.ServerInternal, l *logrus.Logger) (result string, err error)
	GetMemoryInfo(s runner.ServerInternal, l *logrus.Logger) (result string, err error)
	GetMountPointInfo(s runner.ServerInternal, l *logrus.Logger) (result string, err error)
}

type OsExecutor struct{}

func (osExecutor OsExecutor) DoRequest(s runner.ServerInternal, r runner.RunItem) (result string, err error) {
	re := s.ReturnRunResult(r)
	if re.Err != nil {
		return "", re.Err
	}
	return strings.TrimSuffix(re.StdOut, "\n"), nil
}

func (osExecutor OsExecutor) GetHostName(s runner.ServerInternal, l *logrus.Logger) (result string, err error) {
	defer l.Debugf("[%s] hostname -> %s", s.Host, result)
	return osExecutor.DoRequest(s, runner.RunItem{Logger: l, Cmd: PrintHostnameShell})
}

func (osExecutor OsExecutor) GetKernelVersion(s runner.ServerInternal, l *logrus.Logger) (result string, err error) {
	defer l.Debugf("[%s] kernel version -> %s", s.Host, result)
	return osExecutor.DoRequest(s, runner.RunItem{Logger: l, Cmd: PrintKernelVersionShell})
}

func (osExecutor OsExecutor) GetSystemVersion(s runner.ServerInternal, l *logrus.Logger) (result string, err error) {
	defer l.Debugf("[%s] system version -> %s", s.Host, result)
	return osExecutor.DoRequest(s, runner.RunItem{Logger: l, Cmd: PrintOSVersionShell})
}

func (osExecutor OsExecutor) GetCPUInfo(s runner.ServerInternal, l *logrus.Logger) (result string, err error) {
	defer l.Debugf("[%s] cpu info -> %s", s.Host, result)
	return osExecutor.DoRequest(s, runner.RunItem{Logger: l, Cmd: PrintCPUInfoShell})
}

func (osExecutor OsExecutor) GetCPULoadAverage(s runner.ServerInternal, l *logrus.Logger) (result string, err error) {
	defer l.Debugf("[%s] cpu loadaverage info -> %s", s.Host, result)
	return osExecutor.DoRequest(s, runner.RunItem{Logger: l, Cmd: PrintCPULoadavgShell})
}

func (osExecutor OsExecutor) GetMemoryInfo(s runner.ServerInternal, l *logrus.Logger) (result string, err error) {
	defer l.Debugf("[%s] memory info -> %s", s.Host, result)
	return osExecutor.DoRequest(s, runner.RunItem{Logger: l, Cmd: PrintMemInfoShell})
}

func (osExecutor OsExecutor) GetMountPointInfo(s runner.ServerInternal, l *logrus.Logger) (result string, err error) {
	defer l.Debugf("[%s] mount point info -> %s", s.Host, result)
	return osExecutor.DoRequest(s, runner.RunItem{Logger: l, Cmd: PrintMountInfoShell})
}

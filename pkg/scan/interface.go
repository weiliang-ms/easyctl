package scan

import (
	"github.com/sirupsen/logrus"
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"strings"
)

//go:generate mockery --name=HandleOSInterface
type HandleOSInterface interface {
	DoRequest(doRequestItem runner.DoRequestItem) (string, error)
	GetHostName(s runner.ServerInternal, l *logrus.Logger) (result string, err error)
	GetKernelVersion(s runner.ServerInternal, l *logrus.Logger) (result string, err error)
	GetSystemVersion(s runner.ServerInternal, l *logrus.Logger) (result string, err error)
	GetCPUInfo(s runner.ServerInternal, l *logrus.Logger) (result string, err error)
	GetCPULoadAverage(s runner.ServerInternal, l *logrus.Logger) (result string, err error)
	GetMemoryInfo(s runner.ServerInternal, l *logrus.Logger) (result string, err error)
	GetMountPointInfo(s runner.ServerInternal, l *logrus.Logger) (result string, err error)
}

type Handler struct{}

func (h Handler) DoRequest(doRequestItem runner.DoRequestItem) (string, error) {
	re := doRequestItem.S.ReturnRunResult(doRequestItem.R)
	if re.Err != nil && !doRequestItem.Mock {
		return "", re.Err
	}
	return strings.TrimSuffix(re.StdOut, "\n"), nil
}

func (h Handler) GetHostName(s runner.ServerInternal, l *logrus.Logger) (result string, err error) {
	defer l.Debugf("[%s] hostname -> %s", s.Host, result)
	return h.DoRequest(runner.DoRequestItem{S: s, R: runner.RunItem{Logger: l, Cmd: PrintHostnameShell}})
}

func (h Handler) GetKernelVersion(s runner.ServerInternal, l *logrus.Logger) (result string, err error) {
	defer l.Debugf("[%s] kernel version -> %s", s.Host, result)
	return h.DoRequest(runner.DoRequestItem{S: s, R: runner.RunItem{Logger: l, Cmd: PrintKernelVersionShell}})
}

func (h Handler) GetSystemVersion(s runner.ServerInternal, l *logrus.Logger) (result string, err error) {
	defer l.Debugf("[%s] system version -> %s", s.Host, result)
	return h.DoRequest(runner.DoRequestItem{S: s, R: runner.RunItem{Logger: l, Cmd: PrintOSVersionShell}})
}

func (h Handler) GetCPUInfo(s runner.ServerInternal, l *logrus.Logger) (result string, err error) {
	defer l.Debugf("[%s] cpu info -> %s", s.Host, result)
	return h.DoRequest(runner.DoRequestItem{S: s, R: runner.RunItem{Logger: l, Cmd: PrintCPUInfoShell}})
}

func (h Handler) GetCPULoadAverage(s runner.ServerInternal, l *logrus.Logger) (result string, err error) {
	defer l.Debugf("[%s] cpu loadaverage info -> %s", s.Host, result)
	return h.DoRequest(runner.DoRequestItem{S: s, R: runner.RunItem{Logger: l, Cmd: PrintCPULoadavgShell}})
}

func (h Handler) GetMemoryInfo(s runner.ServerInternal, l *logrus.Logger) (result string, err error) {
	defer l.Debugf("[%s] memory info -> %s", s.Host, result)
	return h.DoRequest(runner.DoRequestItem{S: s, R: runner.RunItem{Logger: l, Cmd: PrintMemInfoShell}})
}

func (h Handler) GetMountPointInfo(s runner.ServerInternal, l *logrus.Logger) (result string, err error) {
	defer l.Debugf("[%s] mount point info -> %s", s.Host, result)
	return h.DoRequest(runner.DoRequestItem{S: s, R: runner.RunItem{Logger: l, Cmd: PrintMountInfoShell}})
}

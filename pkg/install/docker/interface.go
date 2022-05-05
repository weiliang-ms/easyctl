package docker

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/weiliang-ms/easyctl/pkg/install"
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"github.com/weiliang-ms/easyctl/pkg/runner/scp"
	strings2 "github.com/weiliang-ms/easyctl/pkg/util/strings"
	"time"
)

//go:generate mockery --name=HandlerInterface
type HandlerInterface interface {
	Boot(server runner.ServerInternal, local bool, logger *logrus.Logger, timeout time.Duration) install.BootErr
	Prune(server runner.ServerInternal, local bool, logger *logrus.Logger, timeout time.Duration) install.PruneErr
	Exec(cmd string, server runner.ServerInternal, local bool, logger *logrus.Logger, timeout time.Duration) error
	SetUpRuntime(cmd string, server runner.ServerInternal, local bool, logger *logrus.Logger, timeout time.Duration) install.SetUpRuntimeErr
	SetSystemd(cmd string, server runner.ServerInternal, local bool, logger *logrus.Logger, timeout time.Duration) install.SetSystemdErr
	SetConfig(cmd string, server runner.ServerInternal, local bool, logger *logrus.Logger, timeout time.Duration) install.SetConfigErr
	Install(cmd string, server runner.ServerInternal, local bool, logger *logrus.Logger, timeout time.Duration) install.InstallErr
	Detect(cmd string, server runner.ServerInternal, local bool, logger *logrus.Logger, timeout time.Duration) error
	HandPackage(server runner.ServerInternal, filePath string, local bool, logger *logrus.Logger, timeout time.Duration) error
}

type Handler struct{}

func (h Handler) Exec(cmd string, server runner.ServerInternal, local bool, logger *logrus.Logger, timeout time.Duration) error {
	if local {
		logger.Info("[docker] 本地执行...")
		return runner.LocalRun(cmd, logger).Err
	}

	return runner.RunOnNode(cmd, server, timeout, logger).Err
}

func (h Handler) Prune(server runner.ServerInternal, local bool, logger *logrus.Logger, timeout time.Duration) install.PruneErr {
	host := h.host(server, local)
	return install.PruneErr{
		Host: host,
		Err:  h.Exec(PruneDockerShell, server, local, logger, timeout),
	}
}

func (h Handler) SetUpRuntime(cmd string, server runner.ServerInternal, local bool, logger *logrus.Logger, timeout time.Duration) install.SetUpRuntimeErr {
	host := h.host(server, local)
	return install.SetUpRuntimeErr{
		Host: host,
		Err:  h.Exec(cmd, server, local, logger, timeout),
	}
}

func (h Handler) Boot(server runner.ServerInternal, local bool, logger *logrus.Logger, timeout time.Duration) install.BootErr {
	host := h.host(server, local)
	return install.BootErr{
		Host: host,
		Err:  h.Exec(bootDockerShell, server, local, logger, timeout),
	}
}

func (h Handler) SetSystemd(cmd string, server runner.ServerInternal, local bool, logger *logrus.Logger, timeout time.Duration) install.SetSystemdErr {
	host := h.host(server, local)
	return install.SetSystemdErr{
		Host: host,
		Err:  h.Exec(cmd, server, local, logger, timeout),
	}
}

func (h Handler) SetConfig(cmd string, server runner.ServerInternal, local bool, logger *logrus.Logger, timeout time.Duration) install.SetConfigErr {
	host := h.host(server, local)
	return install.SetConfigErr{
		Host: host,
		Err:  h.Exec(cmd, server, local, logger, timeout),
	}
}

func (h Handler) Install(cmd string, server runner.ServerInternal, local bool, logger *logrus.Logger, timeout time.Duration) install.InstallErr {
	host := h.host(server, local)
	return install.InstallErr{
		Host: host,
		Err:  h.Exec(cmd, server, local, logger, timeout),
	}
}

func (h Handler) Detect(cmd string, server runner.ServerInternal, local bool, logger *logrus.Logger, timeout time.Duration) error {
	return nil
}

func (h Handler) HandPackage(server runner.ServerInternal, filePath string, local bool, logger *logrus.Logger, timeout time.Duration) error {
	if local {
		shell := fmt.Sprintf("cp %s /tmp/%s", filePath, filePath)
		logger.Infof("分发包至本地: %s", shell)
		return runner.LocalRun(shell, logger).Err
	}

	return scp.Scp(&scp.ScpItem{
		Logger: logger,
		SftpExecutor: scp.SftpExecutor{
			SrcPath: filePath,
			DstPath: fmt.Sprintf("/tmp/%s", strings2.SubFileName(filePath)),
			Mode:    0755,
			Server:  server,
		},
		SftpConnectTimeout: 0})
}

func (h Handler) host(server runner.ServerInternal, local bool) string {
	if local {
		return "localhost"
	}
	return server.Host
}

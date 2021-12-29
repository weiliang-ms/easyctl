package deny

import (
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
)

// Selinux 关闭selinux
func Selinux(item command.OperationItem) command.RunErr {
	return runner.RemoteRun(runner.RemoteRunItem{
		B:                   item.B,
		Logger:              item.Logger,
		Cmd:                 closeSELinuxShell,
		RecordErrServerList: false,
	})
}

// todo confirm
//func confirm(reader *bufio.Reader) (string, error) {
//	for {
//		input, err := reader.ReadString('\n')
//		if err != nil {
//			return "", err
//		}
//		input = strings.TrimSpace(input)
//
//		if input != "" && (input == "yes" || input == "no") {
//			return input, nil
//		}
//	}
//}

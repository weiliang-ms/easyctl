package deny

import (
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
)

// Selinux 关闭selinux
func Selinux(item command.OperationItem) error {
	return runner.RemoteRun(item.B, item.Logger, closeSELinuxShell)
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

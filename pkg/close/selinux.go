package close

import (
	"bufio"
	"fmt"
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"github.com/weiliang-ms/easyctl/pkg/util"
	"os"
	"strings"
	"sync"
)

const closeSELinuxShell = `
if [ "$(getenforce)" == "Disabled" ];then
	echo "已关闭，无需重复关闭" > 2
	exit 3
fi
setenforce 0
sed -i 's/SELINUX=enforcing/SELINUX=disabled/' /etc/selinux/config
`

func (ac *Actuator) SeLinux() {
	ac.parseServerList().seLinuxCmd().execute("关闭selinux", util.AlreadyCloseCode)
	// todo: 多机按退出状态提示重启
	ac.interact("selinux已关闭，请问是否重启上述主机永久生效?")
}

func (ac *Actuator) seLinuxCmd() *Actuator {
	ac.Cmd = closeSELinuxShell
	return ac
}

func (ac *Actuator) interact(msg string) {

	fmt.Printf(fmt.Sprintf("%s [yes/no]: ", msg))

	reader := bufio.NewReader(os.Stdin)
	input, err := confirm(reader)
	if err != nil {
		panic(err)
	}

	if input == "no" {
		os.Exit(0)
	}

	if len(ac.ServerList.Server) <= 0 {
		runner.Shell("reboot")
	} else {
		// 并发关闭
		var wg sync.WaitGroup
		for _, v := range ac.ServerList.Server {
			wg.Add(1)
			// todo: 本机最后执行，防止进程退出
			go func(server runner.Server) {
				server.RemoteShell("reboot")
				defer wg.Done()
			}(v)
		}

		wg.Wait()
	}
}

func confirm(reader *bufio.Reader) (string, error) {
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		input = strings.TrimSpace(input)

		if input != "" && (input == "yes" || input == "no") {
			return input, nil
		}
	}
}

// 排除本机
func excludeLocalHost(s []runner.Server) (servers []runner.Server) {
	// 获取本机IP列表
	re := runner.Shell("hostname -I")
	if re.ExitCode == 0 {
		localAddresses := strings.Split(strings.Trim(re.StdOut, "\n"), " ")
		for _, v := range s {
			if !util.SliceContain(localAddresses, v.Host) {
				servers = append(servers, v)
			}
		}
	}
	return servers
}

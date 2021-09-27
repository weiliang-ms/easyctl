package deny

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"os"
	"sort"
)

const (
	disableFirewallShell = "systemctl disable firewalld --now"
	denyPingShell        = `
sed -i "/net.ipv4.icmp_echo_ignore_all/d" /etc/sysctl.conf
echo "net.ipv4.icmp_echo_ignore_all=1"  >> /etc/sysctl.conf
sysctl -p
`
	closeSELinuxShell = `
if [ "$(getenforce)" == "Disabled" ];then
	echo "已关闭，无需重复关闭"
	exit 0
fi
setenforce 0
sed -i 's/SELINUX=enforcing/SELINUX=disabled/' /etc/selinux/config
`
)

func Item(b []byte, logger *logrus.Logger, cmd string) error {

	results, err := GetResult(b, logger, cmd)
	if err != nil {
		return err
	}
	var data [][]string

	for _, v := range results {
		data = append(data, []string{v.Host, v.Cmd, fmt.Sprintf("%d", v.Code), v.Status, v.StdOut, v.StdErrMsg})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"IP ADDRESS", "cmd", "exit code", "result", "output", "exception"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	//table.SetRowLine(true)
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.AppendBulk(data) // Add Bulk Data
	table.Render()

	return nil
}

func GetResult(b []byte, logger *logrus.Logger, cmd string) ([]runner.ShellResult, error) {

	servers, err := runner.ParseServerList(b)
	if err != nil {
		return []runner.ShellResult{}, err
	}

	executor := runner.ExecutorInternal{Servers: servers, Script: cmd}

	ch := executor.ParallelRun(logger)

	results := []runner.ShellResult{}

	for re := range ch {
		var result runner.ShellResult
		_ = mapstructure.Decode(re, &result)
		results = append(results, result)
	}

	// todo: ip地址排序
	sort.Sort(runner.ShellResultSlice(results))

	return results, nil
}

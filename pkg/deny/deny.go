package deny

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"os"
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

// Item 执行入口
func Item(b []byte, logger *logrus.Logger, cmd string) error {

	results, err := runner.GetResult(b, logger, cmd)
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

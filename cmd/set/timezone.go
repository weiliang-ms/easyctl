package set

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/runner"
)

func init() {
	setTimeZoneCmd.Flags().StringVarP(&value, "value", "v", "Asia/Shanghai", "时区")
	setTimeZoneCmd.Flags().BoolVarP(&multiNode, "multi-node", "", false, "是否配置多节点")
	setTimeZoneCmd.Flags().StringVarP(&serverListFile, "server-list", "", "server.yaml", "服务器列表")
}

// 配置时区子命令
var setTimeZoneCmd = &cobra.Command{
	Use:     "timezone",
	Short:   "easyctl set tz/timezone [value]",
	Example: "\neasyctl set tz/timezone",
	Aliases: []string{"tz"},
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		setTimeZone()
	},
}

// 配置时区
func setTimeZone() {
	if !multiNode {
		setLocalTZ()
	} else {
		setMultiNodeTZ()
	}
}

func setLocalTZ() {
	local(fmt.Sprintf("配置本机时区,时区为%s...\n", value), timezoneScript())
}

func setMultiNodeTZ() {
	list := runner.ParseServerList(serverListFile, runner.CommonServerList{})
	multiShell(list, timezoneScript())
}

func timezoneScript() string {
	return fmt.Sprintf("\\cp /usr/share/zoneinfo/%s /etc/localtime -R", value)
}

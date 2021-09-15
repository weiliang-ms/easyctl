package set

import (
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/set"
)

func init() {
	timeZoneCmd.Flags().StringVarP(&value, "value", "v", "Asia/Shanghai", "时区")
}

// 配置时区子命令
var timeZoneCmd = &cobra.Command{
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
	ac := &set.Actuator{
		ServerListFile: serverListFile,
		Value:          value,
	}

	ac.TimeZone()
}

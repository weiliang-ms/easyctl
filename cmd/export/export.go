package export

import (
	"github.com/spf13/cobra"
)

var (
	configFile string
)

// RootCmd export命令
var RootCmd = &cobra.Command{
	Use:   "export [OPTIONS] [flags]",
	Short: "导出指令集",
	Args:  cobra.ExactValidArgs(1),
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "配置文件")
	RootCmd.AddCommand(imageCmd)
	RootCmd.AddCommand(localImagesCmd)
	RootCmd.AddCommand(chartCmd)
}

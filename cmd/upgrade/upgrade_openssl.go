package upgrade

import (
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/upgrade"
)

func init() {
	upgradeOpensslCmd.Flags().StringVarP(&filePath, "file-path", "f", "", "文件路径或url")
	_ = upgradeOpensslCmd.MarkFlagRequired("file-path")
}

// install kernel
var upgradeOpensslCmd = &cobra.Command{
	Use:   "openssl [flags]",
	Short: "upgrade openssl through easyctl...",
	Run: func(cmd *cobra.Command, args []string) {
		upgradeOpenssl()
	},
	Args: cobra.NoArgs,
}

func upgradeOpenssl() {
	upgrade := upgrade.Actuator{
		ServerListFile: serverListFile,
		FilePath:       filePath,
	}
	upgrade.Openssl()
}

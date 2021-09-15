package upgrade

import (
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/upgrade"
)

var opensslDir string

func init() {
	upgradeOpenSSHCmd.Flags().StringVarP(&filePath, "file-path", "f", "", "文件路径或url")
	upgradeOpenSSHCmd.Flags().StringVarP(&opensslDir, "openssl-path", "", "", "openssl路径")
	_ = upgradeOpenSSHCmd.MarkFlagRequired("file-path")
}

// install kernel
var upgradeOpenSSHCmd = &cobra.Command{
	Use:   "openssh [flags]",
	Short: "更新openssh",
	Run: func(cmd *cobra.Command, args []string) {
		upgradeOpenSSH()
	},
	Args: cobra.NoArgs,
}

func upgradeOpenSSH() {
	upgrade := upgrade.Actuator{
		ServerListFile: serverListFile,
		FilePath:       filePath,
		OpensslDir:     opensslDir,
	}
	upgrade.OpenSSH()
}

package cmd

import (
	"easyctl/asset"
	"easyctl/util"
	"fmt"
	"github.com/spf13/cobra"
)

type harbor struct {
	Address  string
	Username string
	Password string
}

const HarborVersion = "v2.1.4"

var (
	domain string // harbor域名
	ssl    bool
)

func init() {
	installHarborCmd.Flags().BoolVarP(&offline, "offline", "", false, "是否离线安装")
	installHarborCmd.Flags().BoolVarP(&ssl, "ssl", "", false, "是否开启ssl")
	installHarborCmd.Flags().StringVarP(&domain, "domain", "", "harbor.wl.com", "域名")
	installHarborCmd.MarkFlagRequired("domain")
}

// install docker-compose
var installHarborCmd = &cobra.Command{
	Use:   "harbor [flags]",
	Short: "install harbor through easyctl...",
	Run: func(cmd *cobra.Command, args []string) {
		if offline {
			localHarborOffline()
		} else {

		}
	},
	Args: cobra.NoArgs,
}

// 单机本地离线
func localHarborOffline() {

	script, _ := asset.Asset("static/script/install_harbor.sh")
	fmt.Println("安装harbor...")
	util.Run(fmt.Sprintf("ssl=%t version=%s domain=%s %s", ssl, HarborVersion, domain, string(script)))
}

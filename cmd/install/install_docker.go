package install

import (
	"easyctl/asset"
	"easyctl/util"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	installDockerCmd.Flags().BoolVarP(&offline, "offline", "", false, "是否离线安装")
	//installDockerCmd.MarkFlagRequired("offline")
}

// install docker
var installDockerCmd = &cobra.Command{
	Use:   "docker [flags]",
	Short: "install docker through easyctl...",
	Run: func(cmd *cobra.Command, args []string) {
		if offline {
			installDockerOffline()
		} else {

		}
	},
	Args: cobra.NoArgs,
}

// 单机本地离线
func installDockerOffline() {
	// 判断是否存在离线包
	filePath := fmt.Sprintf("%s/resources/repo/%s/repo.tar.gz",
		util.CurrentPath(), docker)
	_, err := os.Stat(filePath)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	script, _ := asset.Asset("static/script/install_offline.sh")
	fmt.Printf("开始安装%s...\n", docker)
	util.Run(fmt.Sprintf("package_name=%s %s", docker, string(script)))

	// 启动
	fmt.Printf("开始启动%s...\n", docker)
	util.Run("systemctl enable docker --now")
}

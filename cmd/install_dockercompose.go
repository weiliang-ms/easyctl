package cmd

import (
	"easyctl/util"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
)

func init() {
	installDockerComposeCmd.Flags().BoolVarP(&offline, "offline", "", false, "是否离线安装")
	installCmd.AddCommand(installDockerComposeCmd)
}

// install docker-compose
var installDockerComposeCmd = &cobra.Command{
	Use:   "docker-compose [flags]",
	Short: "install docker-compose through easyctl...",
	Run: func(cmd *cobra.Command, args []string) {
		if offline {
			installDockerComposeOffline()
		} else {

		}
	},
	Args: cobra.NoArgs,
}

// 单机本地离线
func installDockerComposeOffline() {

	b, err := ioutil.ReadFile(fmt.Sprintf("%s/resources/repo/%s/%s",
		util.CurrentPath(), dockerCompose, dockerCompose))

	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	ioutil.WriteFile("/usr/bin/docker-compose", b, 0755)
	util.Run("docker-compose -v")
}

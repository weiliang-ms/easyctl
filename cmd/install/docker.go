package install

import (
	// embed
	_ "embed"
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/install"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"log"
)

//go:embed asset/docker.yaml
var dockerConfig []byte

//go:embed asset/docker-local.yaml
var dockerLocalConfig []byte

// install docker
var dockerCmd = &cobra.Command{
	Use:   "docker-ce [flags]",
	Short: "安装docker-ce",
	Run: func(cmd *cobra.Command, args []string) {
		if runErr := command.SetExecutorDefault(
			command.Item{
				Cmd:            cmd,
				Fnc:            install.Docker,
				DefaultConfig:  dockerConfigContent(),
				ConfigFilePath: configFile,
				Local:          local,
			}); runErr.Err != nil {
			log.Println(runErr.Msg)
			panic(runErr.Err)
		}
	},
}

func dockerConfigContent() []byte {
	if local {
		return dockerLocalConfig
	}
	return dockerConfig
}

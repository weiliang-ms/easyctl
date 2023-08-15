package export

import (
	// embed
	_ "embed"
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/export/harbor"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"log"
)

//go:embed asset/harbor.yaml
var harborConfig []byte

var imageCmd = &cobra.Command{
	Use:   "harbor-image-list [flags]",
	Short: "导出harbor项目内的镜像列表",
	Run: func(cmd *cobra.Command, args []string) {
		if runErr := command.SetExecutorDefault(
			command.Item{Cmd: cmd, Fnc: harbor.ImageList, ConfigFilePath: configFile, DefaultConfig: harborConfig}); runErr.Err != nil {
			log.Println(runErr.Msg)
			panic(runErr.Err)
		}
	},
}

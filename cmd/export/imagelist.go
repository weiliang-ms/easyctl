package export

import (
	// embed
	_ "embed"
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/export"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
)

//go:embed asset/harbor.yaml
var harborConfig []byte

var imageCmd = &cobra.Command{
	Use:   "harbor-image-list [flags]",
	Short: "导出harbor项目内的镜像列表",
	Run: func(cmd *cobra.Command, args []string) {
		if runErr := command.SetExecutorDefault(command.Item{Cmd: cmd, Fnc: export.HarborImageList, DefaultConfig: harborConfig}); runErr != nil {
			panic(runErr)
		}
	},
}

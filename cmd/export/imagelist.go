package export

import (
	"github.com/spf13/cobra"
)

var imageCmd = &cobra.Command{
	Use:   "harbor-image-list [flags]",
	Short: "导出harbor项目内的镜像列表",
	Run: func(cmd *cobra.Command, args []string) {

	},
	Args: cobra.NoArgs,
}

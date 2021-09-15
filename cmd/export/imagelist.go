package export

import (
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/api/export"
)

var imageCmd = &cobra.Command{
	Use:   "image-list [flags]",
	Short: "export harbor images list through easyctl...",
	Run: func(cmd *cobra.Command, args []string) {
		exportImage()
	},
	Args: cobra.NoArgs,
}

func exportImage() {
	export.ImageList(serverListFile)
}

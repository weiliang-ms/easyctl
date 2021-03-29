package cmd

import (
	"easyctl/asset"
	"easyctl/util"
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	exportRepoCmd.Flags().StringVarP(&packageName, "package-name", "", "", "rpm package name")
	exportRepoCmd.MarkFlagRequired("package-name")
	exportCmd.AddCommand(exportRepoCmd)
}

// export rpm package and its dependencies then generate local repo
var exportRepoCmd = &cobra.Command{
	Use:   "yum-repo [flags]",
	Short: "export package-name through easyctl...",
	// Example: "\neasyctl download harbor --url=https://github.com/goharbor/harbor/releases/download/v2.1.4/harbor-offline-installer-v2.1.4.tgz",
	Run: func(cmd *cobra.Command, args []string) {
		exportYumPackage()
	},
	Args: cobra.NoArgs,
}

func exportYumPackage() {
	data, _ := asset.Asset("script/export_repo.sh")
	util.Run(fmt.Sprintf("%s %s", "package_name=gcc", string(data)))
}

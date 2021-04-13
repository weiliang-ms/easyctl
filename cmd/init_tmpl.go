package cmd

import (
	"easyctl/asset"
	"easyctl/util"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
)

func init() {
	initTmplCmd.AddCommand(initServerTmplCmd)
	initTmplCmd.AddCommand(initKeepaliveTmplCmd)
}

// 初始化模板命令
var initTmplCmd = &cobra.Command{
	Use:     "init-tmpl [OPTIONS] [flags]",
	Short:   "init-tmpl xxx through easyctl",
	Example: "\neasyctl init-tmpl server\n",
	Args:    cobra.MinimumNArgs(1),
}

// init servers template
var initServerTmplCmd = &cobra.Command{
	Use: "server [flags]",
	Run: func(cmd *cobra.Command, args []string) {
		tmpl, _ := asset.Asset("static/tmpl/server.yaml")
		ioutil.WriteFile(fmt.Sprintf("%s/server.yaml", util.CurrentPath()), tmpl, 0644)
	},
	Args: cobra.NoArgs,
}

// init keepalive servers template
var initKeepaliveTmplCmd = &cobra.Command{
	Use: "keepalived [flags]",
	Run: func(cmd *cobra.Command, args []string) {
		tmpl, _ := asset.Asset("static/tmpl/keepalived.yaml")
		ioutil.WriteFile(fmt.Sprintf("%s/keepalived.yaml", util.CurrentPath()), tmpl, 0644)
	},
	Args: cobra.NoArgs,
}
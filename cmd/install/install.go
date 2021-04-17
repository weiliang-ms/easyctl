package install

import (
	"github.com/spf13/cobra"
)

var (
	domain          string // harbor域名
	ssl             bool
	offline         bool
	offlineFilePath string
	serverListFile  string
)

func init() {
	Cmd.AddCommand(keepaliveCmd)
	Cmd.AddCommand(haproxyCmd)
	Cmd.AddCommand(dockerCmd)
	Cmd.AddCommand(dockerComposeCmd)
	Cmd.AddCommand(harborCmd)
}

//
var Cmd = &cobra.Command{
	Use:   "install [OPTIONS] [flags]",
	Short: "install soft through easyctl",
	Args:  cobra.MinimumNArgs(1),
}

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

const (
	docker        = "docker-ce"
	dockerCompose = "docker-compose"
)

func init() {
	Cmd.AddCommand(installHarborCmd)
	Cmd.AddCommand(keepaliveCmd)
}

//
var Cmd = &cobra.Command{
	Use:     "install [OPTIONS] [flags]",
	Short:   "install soft through easyctl",
	Example: "\neasyctl install docker\n",
	Args:    cobra.MinimumNArgs(1),
}

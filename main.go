package main

import (
	_ "embed"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/cmd/add"
	"github.com/weiliang-ms/easyctl/cmd/clean"
	close2 "github.com/weiliang-ms/easyctl/cmd/close"
	"github.com/weiliang-ms/easyctl/cmd/export"
	"github.com/weiliang-ms/easyctl/cmd/install"
	"github.com/weiliang-ms/easyctl/cmd/secure"
	"github.com/weiliang-ms/easyctl/cmd/set"
	"github.com/weiliang-ms/easyctl/cmd/upgrade"
	"github.com/weiliang-ms/easyctl/pkg/logs"
	"math/rand"
	"os"
	"time"
)

var (
	GitTag    = "2000.01.01.release"
	BuildTime = "2000-01-01T00:00:00+0800"
	Debug     bool
)

var RootCmd = &cobra.Command{
	Use:   "easyctl",
	Short: "Easycf is a tool manage linux settings",
	Long: `A Fast and Flexible Static Site Generator built with
                love by spf13 and friends in Go.
                Complete documentation is available at http://hugo.spf13.com`,
	Args: cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func init() {

	RootCmd.PersistentFlags().BoolVarP(&Debug, "debug", "d", false, "开启debug模式")
	subCmds := []*cobra.Command{
		add.RootCmd,
		clean.RootCmd,
		close2.RootCmd,
		set.RootCmd,
		install.RootCmd,
		upgrade.Cmd,
		//cmd.InitTmplCmd,
		export.RootCmd,
		secure.RootCmd,
		versionCmd,
		completionCmd,
	}

	for _, cmd := range subCmds {
		RootCmd.AddCommand(cmd)
	}

}

func main() {

	rand.Seed(time.Now().UnixNano())

	logs.InitLogs()

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// 输出easyctl版本
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of easyctl",
	Long:  `All software has versions. This is easyctl's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("easyctl v%s -- %s\n", GitTag, BuildTime)
		//fmt.Println("easyctl  v0.6.0 -- alpha -- 2021-06-22-09:58:00")
	},
}

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Long: `To load completions:

Bash:

$ source <(yourprogram completion bash)

# To load completions for each session, execute once:
Linux:
  $ yourprogram completion bash > /etc/bash_completion.d/yourprogram
MacOS:
  $ yourprogram completion bash > /usr/local/etc/bash_completion.d/yourprogram

Zsh:

# If shell completion is not already enabled in your environment you will need
# to enable it.  You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# To load completions for each session, execute once:
$ yourprogram completion zsh > "${fpath[1]}/_yourprogram"

# You will need to start a new shell for this setup to take effect.

Fish:

$ yourprogram completion fish | source

# To load completions for each session, execute once:
$ yourprogram completion fish > ~/.config/fish/completions/yourprogram.fish
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactValidArgs(1),
	Hidden:                true,
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			cmd.Root().GenPowerShellCompletion(os.Stdout)
		}
	},
}

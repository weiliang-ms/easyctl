package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.SetHelpCommand(helpCmd)
}

// 帮助信息
var helpCmd = &cobra.Command{
	Use:   "help",
	Short: "Print the version number of easycfg",
	Long:  `All software has versions. This is easycfg's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("" +
			"Cobra is a CLI library for Go that empowers applications.\n" +
			"This application is a tool to generate the needed files\n" +
			"to quickly create a Cobra application.\n" +
			"\n" +
			"Usage:" + "\n" +
			"  easycfg [command]\n" +
			"\n" +
			"Available Commands:" + "\n" +
			"  set         Set DNS address, ulimit and so on." + "\n" +
			"  help        Help about any command" + "\n" +
			"  init        Initialize a Cobra Application\n" +
			"\n" +
			"Flags:\n" +
			"  -a, --author string    author name for copyright attribution (default \"YOUR NAME\")\n" +
			"      --config string    config file (default is $HOME/.cobra.yaml)\n" +

			"\n" +
			"Use \"easycfg [command] --help\" for more information about a command.")
	},
}

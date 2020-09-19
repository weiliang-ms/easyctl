package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
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

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func hasFlags(cmd *cobra.Command, args []string) bool {
	if len(args) > 0 {
		return true
	}
	fmt.Printf(
		"See '%s --help'.\n\nUsage:  %s\n\nExamples: \n%s\n\n",
		cmd.Short,
		cmd.Short,
		cmd.Example,
	)
	return false
}

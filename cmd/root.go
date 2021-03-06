package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	close2 "github.com/weiliang-ms/easyctl/cmd/close"
	"github.com/weiliang-ms/easyctl/cmd/install"
	"github.com/weiliang-ms/easyctl/util"
	"os"
	"strings"
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
	//RootCmd.AddCommand(install.RootCmd)
	RootCmd.AddCommand(close2.RootCmd)
	//RootCmd.AddCommand(set.RootCmd)
	//RootCmd.AddCommand(stat.RootCmd)
	RootCmd.AddCommand(install.Cmd)
	//RootCmd.AddCommand(download.Cmd)
	//RootCmd.AddCommand(upgrade.Cmd)
	//RootCmd.AddCommand(initTmplCmd)
	//RootCmd.AddCommand(scan.RootCmd)
	//RootCmd.AddCommand(export.RootCmd)
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
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

func findSuggestions(c *cobra.Command, arg string) string {
	if c.DisableSuggestions {
		return ""
	}
	if c.SuggestionsMinimumDistance <= 0 {
		c.SuggestionsMinimumDistance = 2
	}
	suggestionsString := ""
	if suggestions := c.SuggestionsFor(arg); len(suggestions) > 0 {
		suggestionsString += "\n\nDid you mean this?\n"
		for _, s := range suggestions {
			suggestionsString += fmt.Sprintf("\t%v\n", s)
		}
	}
	return suggestionsString
}

func ParseCommand(cmd *cobra.Command, args []string, validArgs []string) error {

	for _, v := range validArgs {
		validArgs = append(validArgs, strings.Split(v, "\t")[0])
	}

	for _, v := range args {
		if !util.StringInSlice(v, validArgs) {
			return fmt.Errorf("invalid argument %q for %q%s", v, cmd.CommandPath(), findSuggestions(cmd, args[0]))
		}
	}
	return nil
}

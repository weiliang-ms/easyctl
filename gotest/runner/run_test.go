package runner

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"os"
	"testing"
)

func TestRun(t *testing.T) {
	b, readErr := os.ReadFile("../../asset/config.yaml")
	if readErr != nil {
		panic(readErr)
	}

	executor, err := runner.ParseExecutor(b)
	executor.Script = "date"
	if err != nil {
		panic(err)
	}

	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	results := executor.ParallelRun(logger)

	var data [][]string

	for v := range results {
		data = append(data, []string{v.Host, v.Cmd, fmt.Sprintf("%d", v.Code), v.Status, v.StdOut, v.StdErrMsg})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"IP ADDRESS", "cmd", "exit code", "result", "output", "exception"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	//table.SetRowLine(true)
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.AppendBulk(data) // Add Bulk Data
	table.Render()
}

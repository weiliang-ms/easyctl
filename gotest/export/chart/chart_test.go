package chart

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/weiliang-ms/easyctl/pkg/export"
	"os"
	"testing"
)

func TestParseHelmRepoConfig(t *testing.T) {
	b, readErr := os.ReadFile("../../../asset/config.yaml")
	if readErr != nil {
		panic(readErr)
	}

	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	c, err := export.ParseHelmRepoConfig(b, logger)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", c)
}

func TestChartList(t *testing.T) {
	b, readErr := os.ReadFile("../../../asset/config.yaml")
	if readErr != nil {
		panic(readErr)
	}

	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	c, err := export.ParseHelmRepoConfig(b, logger)

	executor := export.ChartExecutor{
		Config: c,
		Logger: logger,
	}

	if err != nil {
		panic(err)
	}
	list, err := executor.ChartList()
	if err != nil {
		panic(err)
	}
	logger.Debugf("%v", list)
}

func TestSaveCharts(t *testing.T) {
	b, readErr := os.ReadFile("../../../asset/config.yaml")
	if readErr != nil {
		panic(readErr)
	}

	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	c, err := export.ParseHelmRepoConfig(b, logger)

	executor := export.ChartExecutor{
		Config: c,
		Logger: logger,
	}

	if err != nil {
		panic(err)
	}
	list, err := executor.ChartList()
	if err != nil {
		panic(err)
	}
	logger.Debugf("%v", list)

	executor.Save(list)
}

func TestExportChart(t *testing.T) {
	b, readErr := os.ReadFile("../../../asset/config.yaml")
	if readErr != nil {
		panic(readErr)
	}

	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	err := export.Chart(b, logger)
	if err != nil {
		panic(err)
	}
}

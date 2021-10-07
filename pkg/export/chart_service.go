package export

import "github.com/sirupsen/logrus"

type ChartService interface {
	List() ([]ChartItem, error)
}

// ChartExecutor 与helm仓库交互的执行器
type ChartExecutor struct {
	Config         ChartRepoConfig
	Logger         *logrus.Logger
	ChartListFunc  getChartListFunc
	ChartsByteFunc getChartsByteFunc
}

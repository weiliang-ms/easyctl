package export

import (
	_ "embed"
	"github.com/stretchr/testify/assert"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"testing"
)

func mockGetChartListErr(executor *ChartExecutor) ([]byte, error) {
	b := []byte("ddd")
	return b, nil
}

//go:embed mocks/chart-list.json
var chartListMockByte []byte

func mockGetChartList(executor *ChartExecutor) ([]byte, error) {

	return chartListMockByte, nil
}

func TestChart(t *testing.T) {
	// a.test error parse
	d := `
helm-repo:
  - endpoint: 10.79.160.181
  domain: harbor.chs.neusoft.com
  username: admin
  password: Harbor-12345
  preserveDir: D:\github\easyctl\asset\charts
  package: true
  repo-name: charts
`

	item := command.OperationItem{
		B:      []byte(d),
		Logger: nil,
	}
	assert.EqualError(t, Chart(item), "yaml: line 2: did not find expected '-' indicator")

	// b.测试传入错误的参数
	d = `
helm-repo:
  endpoint: 10.79.160.181
  domain: harbor.chs.neusoft.com
  username: admin
  password: Harbor-12345
  preserveDir: D:\github\easyctl\asset\charts
  package: true
  repo-name: charts
`

	options := make(map[string]interface{})
	options[GetChartListFunc] = Chart
	item.B = []byte(d)

	assert.EqualError(t, Chart(item), "getChartListFunc 入参非法")

	// c.test mock get list function err
	options[GetChartListFunc] = mockGetChartListErr
	item.OptionFunc = options
	assert.EqualError(t, Chart(item), "getChartsByteFunc 入参非法")

	// d.test mock get list function success
	options[GetChartsByteFunc] = GetChartsByte
	item.OptionFunc = options
	assert.EqualError(t, Chart(item), "invalid character 'd' looking for beginning of value")

	options[GetChartListFunc] = mockGetChartList
	item.OptionFunc = options
	assert.EqualError(t, Chart(item), "mkdir : The system cannot find the path specified.")

}

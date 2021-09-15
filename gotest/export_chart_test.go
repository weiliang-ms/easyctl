package gotest

import (
	"github.com/weiliang-ms/easyctl/pkg/export"
	"log"
	"os"
	"testing"
)

func TestExportChart(t *testing.T) {
	const chartConfig = "api/export/helm.yaml"
	if err := export.Chart(chartConfig); err != nil {
		t.Error(err.Error())
	} else {
		log.Printf("chart导出用例测试通过，清理生成文件...")
		os.RemoveAll("charts")
	}
}

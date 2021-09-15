package gotest

import (
	"github.com/weiliang-ms/easyctl/pkg/api/export"
	"os"
	"testing"
)

func TestExportImageList(t *testing.T) {
	const chartConfig = "api/export/harbor.yaml"
	export.ImageList(chartConfig)
	os.RemoveAll("images2")
}

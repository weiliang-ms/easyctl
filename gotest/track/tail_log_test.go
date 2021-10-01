package track

import (
	"github.com/sirupsen/logrus"
	"github.com/weiliang-ms/easyctl/pkg/track"
	"os"
	"testing"
)

func TestTailLog(t *testing.T) {
	b, readErr := os.ReadFile("../../asset/config.yaml")
	if readErr != nil {
		panic(readErr)
	}
	err := track.TaiLog(b, logrus.New())
	if err != nil {
		panic(err)
	}
}

package exec

import (
	"github.com/sirupsen/logrus"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"testing"
)

func TestPing(t *testing.T) {
	content := `
ping:
  - address: 10.58.3
    start: 1
    end: 255
    #port: 22
`
	err := Ping(command.OperationItem{B: []byte(content), Logger: logrus.New()})
	if err.Err != nil {
		panic(err.Err)
	}
}

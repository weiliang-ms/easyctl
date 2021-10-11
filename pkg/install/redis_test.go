package install

import (
	"github.com/sirupsen/logrus"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"github.com/weiliang-ms/easyctl/pkg/util/constant"
	"log"
	"os"
	"testing"
)

func TestRedis(t *testing.T) {
	b := `
server:
  - host: 10.10.10.1
    username: root
    password: 123456
    port: 22
excludes:
  - 192.168.235.132
redis:
  password: "ddd"
  port: 33333
`
	os.Setenv(constant.SshNoTimeout, "true")
	var item command.OperationItem
	item.Logger = logrus.New()
	item.UnitTest = true
	item.Logger.SetLevel(logrus.DebugLevel)
	item.B = []byte(b)
	runErr := Redis(item)

	if runErr.Err != nil {
		log.Println(runErr.Msg)
		panic(runErr.Err)
	}
}

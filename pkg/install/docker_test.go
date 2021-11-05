package install

import (
	"github.com/sirupsen/logrus"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"github.com/weiliang-ms/easyctl/pkg/util/constant"
	"log"
	"os"
	"testing"
)

func TestDocker(t *testing.T) {
	b := `
docker:
  package: docker-19.03.15.tgz   # 二进制安装包目录
  preserveDir: /data/lib/docker  # docker数据持久化目录
  insecureRegistries: # 非https仓库列表
    - gcr.azk8s.cn
    - quay.azk8s.cn
  registryMirrors:               # 镜像源
`
	os.Setenv(constant.SshNoTimeout, "true")
	var item command.OperationItem
	item.Logger = logrus.New()
	item.UnitTest = true
	item.Logger.SetLevel(logrus.DebugLevel)
	item.B = []byte(b)
	runErr := Docker(item)

	if runErr.Err != nil {
		log.Println(runErr.Msg)
		panic(runErr.Err)
	}
}

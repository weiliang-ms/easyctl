package install

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"github.com/weiliang-ms/easyctl/pkg/util/constant"
	"os"
	"testing"
	"time"
)

// todo 合并测试用例

func TestRedisCluster(t *testing.T) {
	os.Setenv(constant.SshNoTimeout, "true")

	content := `
server:
  - host: 10.10.10.[1:3]
    username: root
    password: 123456
    port: 22
excludes:
  - 192.168.235.132
redis-cluster:
  password: "ddd"
  cluster-type: 0 # [0] 本地伪集群 ; [1] 三节点3分片2副本 ; [2] 6节点3分片2副本
  package: /root/redis-5.0.13.tar.gz
  listenPorts:
    - 12341
    - 12342
    - 12343
    - 12344
    - 12345
    - 12346
`
	assert.Equal(t, command.RunErr{},
		RedisCluster(command.OperationItem{B: []byte(content), Logger: logrus.New(), UnitTest: true, SSHTimeout: time.Millisecond}))
}

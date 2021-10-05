package install

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"github.com/weiliang-ms/easyctl/pkg/util/errors"
	"os"
	"testing"
)

func TestParse(t *testing.T) {
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
`
	var servers []runner.ServerInternal
	for i := 1; i < 4; i++ {
		servers = append(servers, runner.ServerInternal{
			Host:     fmt.Sprintf("10.10.10.%d", i),
			Port:     "22",
			Username: "root",
			Password: "123456",
		})
	}
	config := redisClusterConfig{
		Logger:        logrus.New(),
		ConfigContent: []byte(content),
	}

	err := config.Parse()
	assert.Nil(t, err)

	assert.Equal(t, "ddd", config.Password)
	assert.Equal(t, "/root/redis-5.0.13.tar.gz", config.Package)
	assert.Equal(t, servers, config.Servers)
}

func TestDetect(t *testing.T) {
	var config redisClusterConfig
	var err error
	err = config.Detect()
	assert.Equal(t, errors.FileNotFoundErr(config.Package), err)

	f, _ := os.Create("1.tar.gz")
	config.Package = "1.tar.gz"
	config.Logger = logrus.New()

	config.CluterType = threeNodesThreeShards
	err = config.Detect()
	assert.Equal(t, errors.NumNotEqualErr("节点", 3, 0), err)

	config.CluterType = sixNodesThreeShards
	err = config.Detect()
	assert.Equal(t, errors.NumNotEqualErr("节点", 6, 0), err)

	f.Close()
	fmt.Println(os.Remove("1.tar.gz"))
}

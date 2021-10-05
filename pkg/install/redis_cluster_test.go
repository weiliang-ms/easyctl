package install

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"github.com/weiliang-ms/easyctl/pkg/util/errors"
	"os"
	"os/exec"
	"sort"
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

	actualServers := runner.InternelServersSlice{}
	err := config.Parse()
	actualServers = config.Servers
	sort.Sort(actualServers)

	assert.Nil(t, err)
	assert.Equal(t, "ddd", config.Password)
	assert.Equal(t, "/root/redis-5.0.13.tar.gz", config.Package)
	assert.Equal(t, servers, config.Servers)

	// test yaml.Unmarshal RedisClusterConfig err
	ddd := `
server:
  - host: 10.10.10.[1:3]
    username: root
    password: 123456
    port: 22
excludes:
  - 192.168.235.132
redis-cluster:
  password: "ddd"
  cluster-type: "0" # [0] 本地伪集群 ; [1] 三节点3分片2副本 ; [2] 6节点3分片2副本
  package: /root/redis-5.0.13.tar.gz
`
	ccc := redisClusterConfig{}
	ccc.ConfigContent = []byte(ddd)
	err = ccc.Parse()
	assert.Errorf(t, err, "Expected nil, but got: &yaml.TypeError{Errors:[]string{\"line 11: cannot unmarshal !!str `0` into install.RedisClusterType")

	// test yaml.Unmarshal RedisClusterConfig err
	aaa := `
server:
   host: 10.10.10.[1:3]
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
	bbb := redisClusterConfig{}
	bbb.ConfigContent = []byte(aaa)
	err = bbb.Parse()
	assert.Errorf(t, err, "[redis-cluster] 反序列化主机列表失败 -> yaml: unmarshal errors:\n  line 3: cannot unmarshal !!map into []runner.ServerExternal")
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

	// test local gcc exist
	config.CluterType = local
	err = config.Detect()
	assert.Equal(t, runner.ExecutorInternal{Script: "gcc -v", Logger: config.Logger}, config.Executor)

	// install & detect -> return nil
	exec.Command("apt install -y gcc").Run()
	err = config.Detect()
	assert.Nil(t, err)

	f.Close()
	os.Remove("1.tar.gz")
}

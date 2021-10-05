package install

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/weiliang-ms/easyctl/pkg/util/errors"
	"os"
	"testing"
)

//func RedisCluster(b []byte, logger *logrus.Logger) error {
//	var config RedisClusterConfig
//
//	logger.Info("解析redis cluster安装配置")
//	if err := yaml.Unmarshal(b, &config); err != nil {
//		return err
//	}
//
//	// 深拷贝属性
//	redisCluster := config.deepCopy()
//	redisCluster.Logger = logger
//	servers, err := runner.ParseServerList(b, logger)
//
//	if err != nil {
//		return fmt.Errorf("[redis-cluster] 反序列化主机列表失败 -> %v", err)
//	}
//	redisCluster.Servers = servers
//
//	return install(redisCluster)
//}

func TestRedisCluster(t *testing.T) {
	logger := logrus.New()
	content := `
redis-cluster:
  paasword: ""
  cluster-type: 0 # [0] 本地伪集群 ; [1] 三节点3分片2副本 ; [2] 6节点3分片2副本
  package: D:\github\easyctl\asset\install\redis-5.0.13.tar.gz
`
	err := RedisCluster([]byte(content), logger)
	assert.Nil(t, err)
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

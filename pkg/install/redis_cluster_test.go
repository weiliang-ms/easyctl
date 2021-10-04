package install

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/weiliang-ms/easyctl/pkg/util/errors"
	"github.com/weiliang-ms/easyctl/pkg/util/log"
	"os"
	"testing"
)

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

func TestInstallRedis(t *testing.T) {
	b, readErr := os.ReadFile("../../asset/install/redis.yaml")
	if readErr != nil {
		panic(readErr)
	}

	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&log.CustomFormatter{})

	if err := RedisCluster(b, logger); err != nil {
		logger.Println(err)
	}
}

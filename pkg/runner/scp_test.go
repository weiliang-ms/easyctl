package runner

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestParallelScp(t *testing.T) {
	b, err := ioutil.ReadFile("../../asset/install/redis.yaml")
	assert.Nil(t, err)

	servers, err := ParseServerList(b, logrus.New())
	assert.Nil(t, err)

	item := ScpItem{
		Servers:        servers,
		SrcPath:        "../../asset/install/redis-5.0.13.tar.gz",
		DstPath:        "/tmp/redis-5.0.13.tar.gz",
		Mode:           0644,
		Logger:         logrus.New(),
		ShowProcessBar: true,
	}
	ParallelScp(item)
}

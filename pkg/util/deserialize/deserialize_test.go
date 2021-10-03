package deserialize

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type RedisClusterConfig struct {
	RedisCluster struct {
		Paasword    string `yaml:"paasword"`
		ClusterType int    `yaml:"cluster-type"`
		Package     string `yaml:"package"`
	} `yaml:"redis-cluster"`
}

const config = `
redis-cluster:
  paasword: ""
  cluster-type: 0
  package: "redis-5.0.13.tar.gz"
`

func TestParseYamlConfig(t *testing.T) {
	object := RedisClusterConfig{}
	result, err := ParseYamlConfig([]byte(config), object)
	assert.Nil(t, err)
	object, ok := result.(RedisClusterConfig)
	fmt.Println(ok)
}

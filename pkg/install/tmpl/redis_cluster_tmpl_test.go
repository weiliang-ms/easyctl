package tmpl

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/weiliang-ms/easyctl/pkg/util/tmplutil"
	"testing"
)

func TestRedisCompileTmpl(t *testing.T) {
	content, err := tmplutil.Render(RedisCompileTmpl, tmplutil.TmplRenderData{
		"PackageName": "redis-5.0.12.tar.gz",
	})
	assert.Nil(t, err)
	const expect = `
#!/bin/bash
set -e
cd /tmp
if [ ! -f redis-5.0.12.tar.gz ];then
  echo /tmp/redis-5.0.12.tar.gz Not Found.
  exit 1
fi
tar zxvf redis-5.0.12.tar.gz
packageName=$(echo redis-5.0.12.tar.gz|sed 's#.tar.gz##g')
echo $packageName
cd $packageName
make -j $(nproc)
make install
`
	assert.Equal(t, expect, content)
}

func TestRedisClusterConfigTmpl(t *testing.T) {
	content, err := tmplutil.Render(RedisClusterConfigTmpl, tmplutil.TmplRenderData{
		"Ports":    []int{26379, 26380, 26381},
		"Password": "redis",
	})
	assert.Nil(t, err)
	assert.NotNil(t, content)
	fmt.Println(content)
}

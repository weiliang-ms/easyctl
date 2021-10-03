package tmpl

import (
	"github.com/stretchr/testify/assert"
	"github.com/weiliang-ms/easyctl/pkg/util"
	"testing"
)

func TestRedisCompileTmpl(t *testing.T) {
	content, err := util.Render(RedisCompileTmpl, util.TmplRenderData{
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
sed -i "s#\$(PREFIX)/bin#%s#g" src/Makefile
make -j $(nproc)
make install
`
	assert.Equal(t, expect, content)
}

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

func TestOpenFirewallPortTmpl(t *testing.T) {
	content, err := tmplutil.Render(OpenFirewallPortTmpl, tmplutil.TmplRenderData{
		"Ports": []int{26379, 26380},
	})
	const expect = `
#!/bin/sh
firewall-cmd --zone=public --add-port=26379/tcp --permanent || true
firewall-cmd --zone=public --add-port=26380/tcp --permanent || true
firewall-cmd --reload || true
`
	assert.Equal(t, nil, err)
	assert.Equal(t, expect, content)
	assert.NotNil(t, content)
}

func TestInitClusterTmpl(t *testing.T) {
	content, err := tmplutil.Render(InitClusterTmpl, tmplutil.TmplRenderData{
		"EndpointList": []string{
			"10.10.10.1:26379",
			"10.10.10.1:26380",
			"10.10.10.1:26381",
			"10.10.10.1:26382",
			"10.10.10.1:26383",
			"10.10.10.1:26384",
		},
	})
	expect := `
#!/bin/sh
echo "yes" | /usr/local/bin/redis-cli --cluster create \
10.10.10.1:26379 \
10.10.10.1:26380 \
10.10.10.1:26381 \
10.10.10.1:26382 \
10.10.10.1:26383 \
10.10.10.1:26384 \
--cluster-replicas 1
`
	assert.Equal(t, nil, err)
	assert.Equal(t, expect, content)
	assert.NotNil(t, content)

	// with password
	content, err = tmplutil.Render(InitClusterTmpl, tmplutil.TmplRenderData{
		"EndpointList": []string{
			"10.10.10.1:26379",
			"10.10.10.1:26380",
			"10.10.10.1:26381",
			"10.10.10.1:26382",
			"10.10.10.1:26383",
			"10.10.10.1:26384",
		},
		"Password": "1111",
	})
	expect = `
#!/bin/sh
echo "yes" | /usr/local/bin/redis-cli --cluster create \
10.10.10.1:26379 \
10.10.10.1:26380 \
10.10.10.1:26381 \
10.10.10.1:26382 \
10.10.10.1:26383 \
10.10.10.1:26384 \
--cluster-replicas 1 -a 1111
`
	assert.Equal(t, nil, err)
	assert.Equal(t, expect, content)
	assert.NotNil(t, content)
}

func TestSetRedisServiceTmpl(t *testing.T) {
	content, err := tmplutil.Render(SetRedisServiceTmpl, tmplutil.TmplRenderData{
		"Ports": []int{26379, 26380},
	})

	assert.Nil(t, err)

	fmt.Println(content)
}

/*
	MIT License

Copyright (c) 2020 xzx.weiliang

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

*/
package runner

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/weiliang-ms/easyctl/pkg/util/constant"
	"os"
	"sort"
	"testing"
)

// ipaddress
func TestParseIPAddress(t *testing.T) {
	os.Setenv(constant.SshNoTimeout, "true")
	var slice InternelServersSlice
	var expect []ServerInternal
	var err error
	const d = `
server:
 - host: 10.10.10.1
   username: root
   password: 123456
   port: 22
excludes:
 - 10.10.10.9`
	servers, err := ParseServerList([]byte(d), logrus.New())
	assert.Nil(t, err)

	slice = servers
	sort.Sort(slice)

	expect = append(expect, ServerInternal{
		Host:     "10.10.10.1",
		Port:     "22",
		Username: "root",
		Password: "123456",
	})

	assert.Equal(t, expect, servers)
}

// 测试非法的地址区间
func TestParseInvalidIPRange1(t *testing.T) {

	const d = `
server:
 - host: 10.10.[1:2].1
   username: root
   password: 123456
   port: 22
excludes:
 - 10.10.10.9`
	_, err := ParseServerList([]byte(d), logrus.New())
	assert.Equal(t, fmt.Errorf("10.10.[1:2].1 地址区间非法"), err)
}

// 测试非法的地址区间
func TestParseInvalidIPRange2(t *testing.T) {

	const d = `
server:
 - host: 10.10.10.333:222
   username: root
   password: 123456
   port: 22
excludes:
 - 10.10.10.9`
	os.Setenv(constant.SshNoTimeout, "true")
	servers, err := ParseServerList([]byte(d), logrus.New())
	assert.Nil(t, err)
	assert.Nil(t, servers)
}

// 测试非法的地址区间
func TestParseInvalidIPRange3(t *testing.T) {

	const d = `
server:
 - host: 10.10.10.1:333
   username: root
   password: 123456
   port: 22
excludes:
 - 10.10.10.9`
	servers, err := ParseServerList([]byte(d), logrus.New())
	assert.Nil(t, err)
	assert.Nil(t, servers)
}

// 测试非法的地址区间
func TestParseInvalidIPRange4(t *testing.T) {

	const d = `
server:
 - host: 10.10.10.10+333
   username: root
   password: 123456
   port: 22
excludes:
 - 10.10.10.9`
	servers, _ := ParseServerList([]byte(d), logrus.New())
	assert.Nil(t, servers)
}

// x:y
func TestParseIPRange0(t *testing.T) {
	var slice InternelServersSlice
	var expect []ServerInternal
	var err error
	const d = `
server:
 - host: 10.10.10.[1:3]
   username: root
   password: 123456
   port: 22
excludes:
 - 10.10.10.3`
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	servers, err := ParseServerList([]byte(d), logger)
	assert.Nil(t, err)

	slice = servers
	sort.Sort(slice)

	for i := 1; i < 3; i++ {
		expect = append(expect, ServerInternal{
			Host:     fmt.Sprintf("10.10.10.%d", i),
			Port:     "22",
			Username: "root",
			Password: "123456",
		})
	}

	assert.Equal(t, expect, servers)
}

// [x:y]
func TestParseIPRange1(t *testing.T) {
	var slice InternelServersSlice
	var expect []ServerInternal
	var err error
	const d = `
server:
 - host: 10.10.10.[1:3]
   username: root
   password: 123456
   port: 22
excludes:
 - 10.10.10.3`
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	servers, err := ParseServerList([]byte(d), logger)
	assert.Nil(t, err)

	slice = servers
	sort.Sort(slice)

	for i := 1; i < 3; i++ {
		expect = append(expect, ServerInternal{
			Host:     fmt.Sprintf("10.10.10.%d", i),
			Port:     "22",
			Username: "root",
			Password: "123456",
		})
	}

	assert.Equal(t, expect, servers)
}

// x-y
func TestParseIPRange2(t *testing.T) {
	var slice InternelServersSlice
	var expect []ServerInternal
	var err error
	const d = `
server:
 - host: 10.10.10.1-3
   username: root
   password: 123456
   port: 22
excludes:
 - 10.10.10.3`
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	servers, err := ParseServerList([]byte(d), logger)
	assert.Nil(t, err)

	slice = servers
	sort.Sort(slice)

	for i := 1; i < 3; i++ {
		expect = append(expect, ServerInternal{
			Host:     fmt.Sprintf("10.10.10.%d", i),
			Port:     "22",
			Username: "root",
			Password: "123456",
		})
	}

	assert.Equal(t, expect, servers)
}

// [x-y]
func TestParseIPRange3(t *testing.T) {
	var slice InternelServersSlice
	var expect []ServerInternal
	var err error
	const d = `
server:
 - host: 10.10.10.[1-3]
   username: root
   password: 123456
   port: 22
excludes:
 - 10.10.10.3`
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	servers, err := ParseServerList([]byte(d), logger)
	assert.Nil(t, err)

	slice = servers
	sort.Sort(slice)

	for i := 1; i < 3; i++ {
		expect = append(expect, ServerInternal{
			Host:     fmt.Sprintf("10.10.10.%d", i),
			Port:     "22",
			Username: "root",
			Password: "123456",
		})
	}

	assert.Equal(t, expect, servers)
}

// x..y
func TestParseIPRange4(t *testing.T) {
	var slice InternelServersSlice
	var expect []ServerInternal
	var err error
	const d = `
server:
 - host: 10.10.10.1..3
   username: root
   password: 123456
   port: 22
excludes:
 - 10.10.10.3`
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	servers, err := ParseServerList([]byte(d), logger)
	assert.Nil(t, err)

	slice = servers
	sort.Sort(slice)

	for i := 1; i < 3; i++ {
		expect = append(expect, ServerInternal{
			Host:     fmt.Sprintf("10.10.10.%d", i),
			Port:     "22",
			Username: "root",
			Password: "123456",
		})
	}

	assert.Equal(t, expect, servers)
}

// [x..y]
func TestParseIPRange5(t *testing.T) {
	var slice InternelServersSlice
	var expect []ServerInternal
	var err error
	const d = `
server:
 - host: 10.10.10.[1..3]
   username: root
   password: 123456
   port: 22
excludes:
 - 10.10.10.3`
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	servers, err := ParseServerList([]byte(d), logger)
	assert.Nil(t, err)

	slice = servers
	sort.Sort(slice)

	for i := 1; i < 3; i++ {
		expect = append(expect, ServerInternal{
			Host:     fmt.Sprintf("10.10.10.%d", i),
			Port:     "22",
			Username: "root",
			Password: "123456",
		})
	}

	assert.Equal(t, expect, servers)
}

// 测试excludes[]
func TestParseInvalidExcludes(t *testing.T) {

	const d = `
server:
 - host: 10.10.10.1:3
   username: root
   password: 123456
   port: 22
excludes:
 - 10.10.10.x`
	servers, _ := ParseServerList([]byte(d), logrus.New())
	var slice InternelServersSlice
	var expect []ServerInternal
	slice = servers
	sort.Sort(slice)
	for i := 1; i < 4; i++ {
		expect = append(expect, ServerInternal{
			Host:     fmt.Sprintf("10.10.10.%d", i),
			Port:     "22",
			Username: "root",
			Password: "123456",
		})
	}
	assert.Equal(t, expect, servers)
}

// 测试非数组类型主机列表反序列化
func TestParseServerListErrHosts(t *testing.T) {
	aaa := `
server:
  host: 10.10.10.1
  username: root
  password: 123456
  port: 22
excludes:
- 192.168.235.132
`
	servers, err := ParseServerList([]byte(aaa), logrus.New())
	assert.Equal(t, servers, []ServerInternal{})
	assert.NotNil(t, err)
}

// 解析执行器测试用例
func TestParseExecutor(t *testing.T) {
	// a.测试异常元数据类型
	const b = `
server:
  host: 10.10.10.1
  username: root
  password: 123456
  port: 22
excludes:
- 192.168.235.132
script: 1.sh
`
	_, err := ParseExecutor([]byte(b), nil)
	assert.EqualError(t, err, "yaml: unmarshal errors:\n  line 3: cannot unmarshal !!map into []runner.ServerExternal")

	// b.测试地址段
	const d = `
server:
 - host: 10.10.10.1-3
   username: root
   password: 123456
   port: 22
excludes:
- 192.168.235.132
script: 1.sh
`
	executor, err := ParseExecutor([]byte(d), nil)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(executor.Servers))
}

// 测试解析带有多个分隔符的主机列表配置
func TestParseMultiSplitCharServers(t *testing.T) {
	// a.测试异常元数据类型
	const b = `
server:
  - host: 10.10.10.,-1-3
    username: root
    password: 123456
    port: 22
excludes:
- 192.168.235.132
`
	servers, err := ParseServerList([]byte(b), nil)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(servers))
}

// 测试解析ParseExecutor地址区间异常情况
func TestParseExecutorWithErrIPRange(t *testing.T) {
	// a.测试异常元数据类型
	const b = `
server:
  - host: xxx.xxx.xxx.1-3
    username: root
    password: 123456
    port: 22
excludes:
- 192.168.235.132
`
	executor, err := ParseExecutor([]byte(b), nil)
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(executor.Servers))
}

// host为数组类型
func TestParseHostSliceType(t *testing.T) {
	const b = `
server:
   - host:
       - 192.168.1.1-3
       - 192.168.1.4
     username: root
     password: 123456
     port: 22
excludes:
 - 192.168.235.132
`
	// 反序列化
	s, err := ParseExecutor([]byte(b), logrus.New())
	assert.Equal(t, nil, err)
	assert.Equal(t, 3, len(s.Servers))
}

func TestParseHostsArray(t *testing.T) {
	const b = `
server:
   - host:
       - 192.168.69.175
       - 192.168.71.[159-162]
       - 10.10.10.1-3
     username: root
     password: 123456
     port: 22
   - host: 192.168.0.1
     username: root
     password: 123456
     port: 22
   - host:
       - 192.168.1.1-3
       - 192.168.1.4
     username: root
     password: 123456
     port: 22
`
	executor, err := ParseExecutor([]byte(b), nil)
	assert.Nil(t, err)
	assert.Equal(t, 13, len(executor.Servers))
}

func TestParseHostsArrayWithErr(t *testing.T) {
	const b = `
server:
   - host:
       - xxx.xxx.xxx.1-3
     username: root
     password: 123456
     port: 22
`
	executor, err := ParseExecutor([]byte(b), nil)
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(executor.Servers))
}

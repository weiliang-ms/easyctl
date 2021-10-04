package runner

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

// ipaddress
func TestParseIPAddress(t *testing.T) {
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

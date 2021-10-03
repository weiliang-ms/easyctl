package runner

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func TestParseServerList(t *testing.T) {
	const d = `
server:
  - host: 10.10.10.[1:3]
    username: root
    password: 123456
    port: 22
  - host: 10.10.10.[5:5]
    username: root
    password: 123456
    port: 22
  - host: 10.10.10.2
    username: root
    password: 123456
    port: 22
  - host: 10.10.10.2
    username: root
    password: 123456
    port: 22
  - host: 10.10.10.[14:12]
    username: root
    password: 123456
    port: 22
  - host: 10.10.10.[8:10]
    username: root
    password: 123456
    port: 22
excludes:
  - 10.10.10.9`

	servers, err := ParseServerList([]byte(d), logrus.New())
	assert.Nil(t, err)

	var slice InternelServersSlice
	slice = servers
	sort.Sort(slice)
	fmt.Println(servers)

	var expect []ServerInternal
	ips := []string{"1", "2", "3", "5", "8", "10"}
	for _, v := range ips {
		expect = append(expect, ServerInternal{
			Host:     fmt.Sprintf("10.10.10.%s", v),
			Port:     "22",
			Username: "root",
			Password: "123456",
		})
	}

	assert.Equal(t, expect, servers)
}

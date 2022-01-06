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
	"github.com/stretchr/testify/assert"
	"github.com/weiliang-ms/easyctl/pkg/util/constant"
	"os"
	"runtime"
	"testing"
	"time"
)

func TestGetResult(t *testing.T) {
	os.Setenv(constant.SshNoTimeout, "true")
	// a.测试反序列化失败
	aaa := `
server:
   host: 10.10.10.1
   username: root
   password: 123456
   port: 22
excludes:
 - 192.168.235.132
`
	// b.测试执行结果
	_, err := GetResult([]byte(aaa), nil, "")
	assert.Errorf(t, err, "line 3: cannot unmarshal !!map into []runner.ServerExternal")

	aaa = `
server:
   - host: 10.10.10.1
     username: root
     password: 123456
     port: 22
excludes:
 - 192.168.235.132
`
	re, err := GetResult([]byte(aaa), nil, "")
	assert.Nil(t, err)
	for _, v := range re {
		assert.Equal(t, "10.10.10.1 ssh会话建立失败->dial tcp 10.10.10.1:22: i/o timeout", v.StdErrMsg)
	}
}

func TestRemoteRun(t *testing.T) {
	// a.测试反序列化失败
	aaa := `
server:
   host: 10.10.10.1
   username: root
   password: 123456
   port: 22
excludes:
 - 192.168.235.132
`
	//assert.Nil(t, RemoteRun([]byte(aaa), nil, ""))
	os.Setenv(constant.SshNoTimeout, "true")
	assert.NotEqual(t, nil, RemoteRun(RemoteRunItem{B: []byte(aaa)}))

	aaa = `
server:
   - host: 10.10.10.1
     username: root
     password: 123456
     port: 22
excludes:
 - 192.168.235.132
`
	err := RemoteRun(RemoteRunItem{B: []byte(aaa)})
	_, ok := err.Err.(runtime.Error)
	if ok {
		assert.Equal(t, true, ok)
		fmt.Println(err.Msg)
	}
}

// ssh连接异常error
func TestSFtpConnectSSHError(t *testing.T) {
	sftp, err := SftpConnect("root", "ddd", "1.1.1.1", "22", time.Second)
	assert.Nil(t, sftp)
	assert.Error(t, err, "连接ssh失败 dial tcp 1.1.1.1:22: i/o timeout")
}

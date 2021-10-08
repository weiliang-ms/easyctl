package runner

import (
	"github.com/stretchr/testify/assert"
	"github.com/weiliang-ms/easyctl/pkg/util/constant"
	"os"
	"testing"
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
		assert.Equal(t, "ssh会话建立失败->dial tcp 10.10.10.1:22: i/o timeout", v.StdErrMsg)
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
	assert.NotEqual(t, nil, RemoteRun([]byte(aaa), nil, ""))

	aaa = `
server:
   - host: 10.10.10.1
     username: root
     password: 123456
     port: 22
excludes:
 - 192.168.235.132
`
	err := RemoteRun([]byte(aaa), nil, "")
	assert.Nil(t, err)
}

// ssh连接异常error
func TestSFtpConnectSSHError(t *testing.T) {
	os.Setenv(constant.SshNoTimeout, "true")
	sftp, err := sftpConnect("root", "ddd", "1.1.1.1", "22")
	assert.Nil(t, sftp)
	assert.Error(t, err, "连接ssh失败 dial tcp 1.1.1.1:22: i/o timeout")
}

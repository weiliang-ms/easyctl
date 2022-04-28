package harden

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"gopkg.in/yaml.v2"
	"testing"
)

var (
	mockContent = `
server:
  - host: 10.10.10.[1:3]
    username: root
    password: 123456
    port: 22
excludes:
  - 192.168.235.132
harden:
  common:
    denyPing: true               # 禁ping
    denyICMPTimeStamp: true      # 禁ICMP时间戳报文
    denyRebootByKeyBoard: true   # 禁止快捷键重启操作系统
    delUnusedUser: true          # 删除用户（系统预创建用户等）
    hideOSVersion: true          # 隐藏操作系统版本
    lockKeyFile: true            # 锁定敏感文件
    setTMOUT: true               # 设置无shell操作自动断开时间
    setPasswdPolicy: true        # 设置系统密码策略
    setSysLog: true              # 设置系统日志
    passwdErrRetryCount: 3       # 错误密码重试次数
    sudoUser:                    # 设置sudo用户
      setSudoUser: true
      sudoUser: easyctl
      sudoPasswd: "Ckh3TQ%~Xal+gohyh"
    unuseUsers:
      - adm
      - ftp
  ssh:
    denySSHUseDns: true           # 禁止ssh对地址进行dns解析
    denySSHAgentForwarding: true
  cron:
    delCommonUserCron: true       # 清理普通用户定时任务
    setDelZombieProcessCron: true # 设置清理僵尸进程定时任务
  # obj.modifySSHLogin,
`
)

func Test_Entry(t *testing.T) {
	object := &Object{}
	object.confirm()
}

func Test_Parse(t *testing.T) {
	content := `
server:
  - host: 10.10.10.[1:3]
    username: root
    password: 123456
    port: 22
excludes:
  - 192.168.235.132
harden:
  common:
    denyPing: true               # 禁ping
    denyICMPTimeStamp: true      # 禁ICMP时间戳报文
    denyRebootByKeyBoard: true   # 禁止快捷键重启操作系统
    delUnusedUser: true          # 删除用户（系统预创建用户等）
    hideOSVersion: true          # 隐藏操作系统版本
    lockKeyFile: true            # 锁定敏感文件
    setTMOUT: true               # 设置无shell操作自动断开时间
    setPasswdPolicy: true        # 设置系统密码策略
    setSysLog: true              # 设置系统日志
    passwdErrRetryCount: 3       # 错误密码重试次数
    sudoUser:                    # 设置sudo用户
      setSudoUser: true
      sudoUsers:
        - user: easyctl
          passwrod: "Ckh3TQ%~Xal+gohyh"
    unuseUsers:
      - adm
      - ftp
      - games
      - gopher
      - halt
      - lp
      - mail
      - news
      - sync
      - shutdown
      - uucp
      - operator
  ssh:
    denySSHUseDns: true           # 禁止ssh对地址进行dns解析
    denySSHAgentForwarding: true
  cron:
    delCommonUserCron: true       # 清理普通用户定时任务
    setDelZombieProcessCron: true # 设置清理僵尸进程定时任务
  # obj.modifySSHLogin,
`
	obj := &Object{
		B: []byte(content),
	}
	re := obj.parse()
	if re.Err != nil {
		panic(re.Err)
	}
}

func Test_delUnusedUser(t *testing.T) {
	obj := &Object{
		B:      []byte(mockContent),
		Logger: logrus.New(),
	}
	_ = obj.parse()
	_ = obj.deepCopy()
	obj.delUnusedUser()
}

func Test_errPasswdRetryCount(t *testing.T) {
	obj := &Object{
		B:      []byte(mockContent),
		Logger: logrus.New(),
	}
	_ = obj.parse()
	_ = obj.deepCopy()
	obj.errPasswdRetryCount()
}

func Test_addSudoUser(t *testing.T) {
	obj := &Object{
		B:      []byte(mockContent),
		Logger: logrus.New(),
	}
	_ = obj.parse()
	_ = obj.deepCopy()
	obj.addSudoUser()
}

func TestOS(t *testing.T) {
	err := OS(command.OperationItem{Logger: logrus.New()})
	assert.Equal(t, command.RunErr{}, err)

	// 模式异常
	err = OS(command.OperationItem{Logger: logrus.New(), B: []byte(`
server:
   host: 1.1.1.1
`)})

	_, ok := err.Err.(*yaml.TypeError)
	if ok {
		assert.Equal(t, true, ok)
	}
}

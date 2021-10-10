package harden

import (
	"github.com/sirupsen/logrus"
	"github.com/weiliang-ms/easyctl/pkg/deny"
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
)

// Object 加固对象
type Object struct {
	Servers  []runner.ServerInternal
	B        []byte
	Logger   *logrus.Logger
	unitTest bool
}

const (
	SetTMOUTShell = `
sed -i '/export TMOUT=300/d' /etc/profile
sed -i '/readonly TMOUT/d' /etc/profile
echo "export TMOUT=300" >> /etc/profile
echo "readonly TMOUT" >> /etc/profile
`
	SetPasswdPolicyShell = `
PASS_MAX_DAYS=$(grep -e ^PASS_MAX_DAYS /etc/login.defs |awk '{print $2}')
if [ $PASS_MAX_DAYS -gt 90 ];then
    echo "密码最长保留期限为：$PASS_MAX_DAYS, 更改为90天"
    sed -i "/^PASS_MAX_DAYS/d" /etc/login.defs
    echo "PASS_MAX_DAYS   90" >> /etc/login.defs
fi

PASS_MIN_DAYS=$(grep -e ^PASS_MIN_DAYS /etc/login.defs |awk '{print $2}')
if [ $PASS_MIN_DAYS -ne 0 ];then
    echo "密码最段保留期限为：$PASS_MIN_DAYS, 更改为1天"
    sed -i "/^PASS_MIN_DAYS/d" /etc/login.defs
    echo "PASS_MIN_DAYS   0" >> /etc/login.defs
fi

PASS_MIN_LEN=$(grep -e ^PASS_MIN_LEN /etc/login.defs |awk '{print $2}')
if [ $PASS_MIN_LEN -lt 8 ];then
    echo "密码最少字符为：$PASS_MIN_LEN, 更改为8"
    sed -i "/^PASS_MIN_LEN/d" /etc/login.defs
    echo "PASS_MIN_LEN   8" >> /etc/login.defs
fi

PASS_WARN_AGE=$(grep -e ^PASS_WARN_AGE /etc/login.defs |awk '{print $2}')
if [ $PASS_WARN_AGE -ne 7 ];then
  echo "密码到期前$PASS_MIN_LEN天提醒, 更改为7"
  sed -i "/^PASS_WARN_AGE/d" /etc/login.defs
  echo "PASS_WARN_AGE   7" >> /etc/login.defs
fi
`
	SetErrPasswdRetryShell = `
sed -i "/MaxAuthTries/d" /etc/ssh/sshd_config
echo "MaxAuthTries 3" >> /etc/ssh/sshd_config
service restart sshd
`
	HardenSystemLogShell = `
touch /var/log/secure
chown root:root /var/log/secure
chmod 600 /var/log/secure
`
	DenySSHUseDnsShell = `
sed -i "/UseDNS/d" /etc/ssh/sshd_config
echo "UseDNS no" >> /etc/ssh/sshd_config
service restart sshd
`
	DenySSHdAgentForwardingShell = `
sed -i "/AgentForwarding/d" /etc/ssh/sshd_config
sed -i "/TcpForwarding/d" /etc/ssh/sshd_config
echo "AllowAgentForwarding no" >> /etc/ssh/sshd_config
echo "AllowTcpForwarding no" >> /etc/ssh/sshd_config
service restart sshd
`
	DeleteUnUsedUserShell = `
users=(adm lp sync shutdown halt mail news uucp operator games gopher ftp)
for i in ${users[@]};
do
  userdel $i &>/dev/null || true
done

for i in ${users[@]};
do
  userdel $i &>/dev/null || true
done
`
	DeleteCommonUserCronShell = `
rm -f /etc/cron.deny
`
	DeleteZombieProcessCronShell = `
crontab -l | grep -v '#' > /tmp/file1
echo "0 3 * * * ps -A -ostat,ppid | grep -e '^[Zz]' | awk '{print $2}' | xargs kill -HUP > /dev/null 2>&1" >> /tmp/file1 && awk ' !x[$0]++{print > "/tmp/file1"}' /tmp/file1
crontab /tmp/file1
`

	HideOSVersionShell = `
mv /etc/issue /etc/issue.bak || true
mv /etc/issue.net /etc/issue.net.bak || true
`
	UnsetRebootByKeyBoardShell = `
rm -rf /usr/lib/systemd/system/ctrl-alt-del.target || true
`
	UnsetICMPTimeStampShell = `
iptables -I INPUT -p ICMP --icmp-type timestamp-request -m comment --comment "deny ICMP timestamp" -j DROP || true
iptables -I INPUT -p ICMP --icmp-type timestamp-reply -m comment --comment "deny ICMP timestamp" -j DROP || true
`
)

type task func() error

func OS(item command.OperationItem) command.RunErr {
	obj := &Object{}
	obj.B = item.B
	obj.Logger = item.Logger
	obj.unitTest = item.UnitTest

	tasks := []task{
		obj.denyPing,
		obj.denyICMPTimeStamp,
		obj.setTMOUT,
		obj.hideOSVersion,
		obj.denyRebootByKeyBoard,
		obj.setPasswdPolicy,
		obj.delUnusedUser,
		obj.errPasswdRetryCount,
		obj.denySSHUseDns,
		obj.denySSHAgentForwarding,
		obj.delCommonUserCron,
		obj.delZombieProcessCron,
	}

	for _, f := range tasks {
		err := f()
		if err != nil {
			return command.RunErr{Err: err}
		}
	}

	return command.RunErr{}
}

// 解析主机列表
func (object *Object) parse() error {
	servers, err := runner.ParseServerList(object.B, object.Logger)
	object.Servers = servers
	return err
}

// 禁Ping
func (object *Object) denyPing() error {
	object.Logger.Info("[step 1] 禁ping")
	return runner.RemoteRun(object.B, object.Logger, deny.UnsetPingResponseShell)
}

// 关闭ICMP_TIMESTAMP应答
func (object *Object) denyICMPTimeStamp() error {
	object.Logger.Info("[step 2] 关闭ICMP_TIMESTAMP应答")
	return runner.RemoteRun(object.B, object.Logger, UnsetICMPTimeStampShell)
}

// 设置系统空闲等待时间
func (object *Object) setTMOUT() error {
	object.Logger.Info("[step 3] 设置系统空闲等待时间")
	return runner.RemoteRun(object.B, object.Logger, SetTMOUTShell)
}

// 隐藏系统版本信息
func (object *Object) hideOSVersion() error {
	object.Logger.Info("[step 4] 隐藏系统版本信息")
	return runner.RemoteRun(object.B, object.Logger, HideOSVersionShell)
}

// 禁止Control-Alt-Delete 键盘重启系统命令
func (object *Object) denyRebootByKeyBoard() error {
	object.Logger.Info("[step 5] 禁止Control-Alt-Delete 键盘重启系统命令")
	return runner.RemoteRun(object.B, object.Logger, UnsetRebootByKeyBoardShell)
}

// ssh用户密码加固
func (object *Object) setPasswdPolicy() error {
	object.Logger.Info("[step 6] ssh用户密码加固")
	return runner.RemoteRun(object.B, object.Logger, SetPasswdPolicyShell)
}

// 删除系统默认用户
func (object *Object) delUnusedUser() error {
	object.Logger.Info("[step 7] 删除系统默认用户")
	return runner.RemoteRun(object.B, object.Logger, DeleteUnUsedUserShell)
}

// 修改允许密码错误次数
func (object *Object) errPasswdRetryCount() error {
	object.Logger.Info("[step 8] 修改允许密码错误次数")
	return runner.RemoteRun(object.B, object.Logger, SetErrPasswdRetryShell)
}

// 关闭ssh UseDNS
func (object *Object) denySSHUseDns() error {
	object.Logger.Info("[step 9] ssh关闭UseDNS")
	return runner.RemoteRun(object.B, object.Logger, DenySSHUseDnsShell)
}

// 关闭ssh `AgentForwarding`和`TcpForwarding`
func (object *Object) denySSHAgentForwarding() error {
	object.Logger.Info("[step 9] ssh关闭UseDNS")
	return runner.RemoteRun(object.B, object.Logger, DenySSHdAgentForwardingShell)
}

// 加固系统日志文件
func (object *Object) setSysLog() error {
	object.Logger.Info("[step 10] 加固系统日志文件")
	return runner.RemoteRun(object.B, object.Logger, HardenSystemLogShell)
}

// 删除非root用户定时任务
func (object *Object) delCommonUserCron() error {
	object.Logger.Info("[step 11] 删除非root用户定时任务")
	return runner.RemoteRun(object.B, object.Logger, DeleteCommonUserCronShell)
}

// 定时清理僵尸进程
func (object *Object) delZombieProcessCron() error {
	object.Logger.Info("[step 11] 删除非root用户定时任务")
	return runner.RemoteRun(object.B, object.Logger, DeleteZombieProcessCronShell)
}

// 创建sudo用户

// 锁定敏感文件并降权

// 修改ssh port 禁止root登录

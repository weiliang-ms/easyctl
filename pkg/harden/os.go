package harden

import (
	"github.com/lithammer/dedent"
	"github.com/sirupsen/logrus"
	"github.com/weiliang-ms/easyctl/pkg/deny"
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"github.com/weiliang-ms/easyctl/pkg/util/constant"
	"github.com/weiliang-ms/easyctl/pkg/util/tmplutil"
	"text/template"
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
	ModifySSHLoginShell = `
sed -i "/PermitRootLogin/d" /etc/ssh/sshd_config
sed -i "/Port 22/d" /etc/ssh/sshd_config
echo "Port 22122" >> /etc/ssh/sshd_config
echo "PermitRootLogin no" >> /etc/ssh/sshd_config

setenforce 0
firewall-cmd --zone=public --add-port=22122/tcp --permanent || true
firewall-cmd --zone=public --add-port=22122/tcp --permanent || true
firewall-cmd --reload || true

iptables -I INPUT -p tcp -m state --state NEW -m tcp --dport 22122 -j ACCEPT || true
/etc/rc.d/init.d/iptables save || ture
service iptables restart || ture

service sshd restart
`
	SetErrPasswdRetryShell = `
sed -i "/MaxAuthTries/d" /etc/ssh/sshd_config
echo "MaxAuthTries 3" >> /etc/ssh/sshd_config
service sshd restart
`
	HardenSystemLogShell = `
touch /var/log/secure
chown root:root /var/log/secure
chmod 600 /var/log/secure
`
	DenySSHUseDnsShell = `
sed -i "/UseDNS/d" /etc/ssh/sshd_config
echo "UseDNS no" >> /etc/ssh/sshd_config
service sshd restart
`
	DenySSHdAgentForwardingShell = `
sed -i "/AgentForwarding/d" /etc/ssh/sshd_config
sed -i "/TcpForwarding/d" /etc/ssh/sshd_config
echo "AllowAgentForwarding no" >> /etc/ssh/sshd_config
echo "AllowTcpForwarding no" >> /etc/ssh/sshd_config
service sshd restart
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
	LockKeyFileShell = `
chown root:root /etc/{passwd,shadow,group}
chmod 644 /etc/{passwd,group}
chmod 400 /etc/shadow
chattr +i /etc/services || true
chattr +i /etc/passwd /etc/shadow /etc/group /etc/gshadow /etc/inittab
`
)

var addSudoUserTmpl = template.Must(template.New("addSudoUserTmpl").Parse(dedent.Dedent(`
chattr -i /etc/passwd /etc/shadow /etc/group /etc/gshadow /etc/inittab
useradd -m easyctl &>/dev/null || true
echo {{ .Password }} | passwd --stdin easyctl || true
sed -i '/easyctl/d' /etc/sudoers
echo "easyctl        ALL=(ALL)       NOPASSWD: ALL" >> /etc/sudoers
`)))

type task func() command.RunErr

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
		obj.setSysLog,
		obj.delCommonUserCron,
		obj.delZombieProcessCron,
		obj.addSudoUser,
		obj.lockKeyFile,
		obj.modifySSHLogin,
		obj.print,
	}

	for _, f := range tasks {
		err := f()
		if err.Err != nil {
			return err
		}
	}

	return command.RunErr{}
}

// 禁Ping
func (object *Object) denyPing() command.RunErr {
	object.Logger.Info("[step 1] 禁ping")
	return runner.RemoteRun(object.B, object.Logger, deny.DenyPingShell)
}

// 关闭ICMP_TIMESTAMP应答
func (object *Object) denyICMPTimeStamp() command.RunErr {
	object.Logger.Info("[step 2] 关闭ICMP_TIMESTAMP应答")
	return runner.RemoteRun(object.B, object.Logger, UnsetICMPTimeStampShell)
}

// 设置系统空闲等待时间
func (object *Object) setTMOUT() command.RunErr {
	object.Logger.Info("[step 3] 设置系统空闲等待时间")
	return runner.RemoteRun(object.B, object.Logger, SetTMOUTShell)
}

// 隐藏系统版本信息
func (object *Object) hideOSVersion() command.RunErr {
	object.Logger.Info("[step 4] 隐藏系统版本信息")
	return runner.RemoteRun(object.B, object.Logger, HideOSVersionShell)
}

// 禁止Control-Alt-Delete 键盘重启系统命令
func (object *Object) denyRebootByKeyBoard() command.RunErr {
	object.Logger.Info("[step 5] 禁止Control-Alt-Delete 键盘重启系统命令")
	return runner.RemoteRun(object.B, object.Logger, UnsetRebootByKeyBoardShell)
}

// ssh用户密码加固
func (object *Object) setPasswdPolicy() command.RunErr {
	object.Logger.Info("[step 6] ssh用户密码加固")
	return runner.RemoteRun(object.B, object.Logger, SetPasswdPolicyShell)
}

// 删除系统默认用户
func (object *Object) delUnusedUser() command.RunErr {
	object.Logger.Info("[step 7] 删除系统默认用户")
	return runner.RemoteRun(object.B, object.Logger, DeleteUnUsedUserShell)
}

// 修改允许密码错误次数
func (object *Object) errPasswdRetryCount() command.RunErr {
	object.Logger.Info("[step 8] 修改允许密码错误次数")
	return runner.RemoteRun(object.B, object.Logger, SetErrPasswdRetryShell)
}

// 关闭ssh UseDNS
func (object *Object) denySSHUseDns() command.RunErr {
	object.Logger.Info("[step 9] ssh关闭UseDNS")
	return runner.RemoteRun(object.B, object.Logger, DenySSHUseDnsShell)
}

// 关闭ssh `AgentForwarding`和`TcpForwarding`
func (object *Object) denySSHAgentForwarding() command.RunErr {
	object.Logger.Info("[step 10] ssh关闭AgentForwarding")
	return runner.RemoteRun(object.B, object.Logger, DenySSHdAgentForwardingShell)
}

// 加固系统日志文件
func (object *Object) setSysLog() command.RunErr {
	object.Logger.Info("[step 11] 加固系统日志文件")
	return runner.RemoteRun(object.B, object.Logger, HardenSystemLogShell)
}

// 删除非root用户定时任务
func (object *Object) delCommonUserCron() command.RunErr {
	object.Logger.Info("[step 12] 删除非root用户定时任务")
	return runner.RemoteRun(object.B, object.Logger, DeleteCommonUserCronShell)
}

// 定时清理僵尸进程
func (object *Object) delZombieProcessCron() command.RunErr {
	object.Logger.Info("[step 13] 定时清理僵尸进程")
	return runner.RemoteRun(object.B, object.Logger, DeleteZombieProcessCronShell)
}

// 创建sudo用户
func (object *Object) addSudoUser() command.RunErr {
	// todo: 重构调用内部指令
	object.Logger.Infof("[step 14] 添加sudo用户: easyctl 密码: %s", constant.DefaultPassword)
	cmd, _ := tmplutil.Render(addSudoUserTmpl, tmplutil.TmplRenderData{
		"Password": constant.DefaultPassword,
	})
	return runner.RemoteRun(object.B, object.Logger, cmd)
}

// 锁定敏感文件并降权
func (object *Object) lockKeyFile() command.RunErr {
	object.Logger.Info("[step 15] 锁定敏感文件")
	return runner.RemoteRun(object.B, object.Logger, LockKeyFileShell)
}

// 修改ssh port & 禁止root登录
func (object *Object) modifySSHLogin() command.RunErr {
	object.Logger.Info("[step 16] 调整ssh登录端口为: 22122，禁止root直接登录.")
	return runner.RemoteRun(object.B, object.Logger, ModifySSHLoginShell)
}

func (object *Object) print() command.RunErr {
	object.Logger.Infof("[done] 安全加固完毕，目标主机连方式改为：\n"+
		"ssh端口: 22122\n"+
		"ssh用户: easyctl\n"+
		"ssh密码: %s", constant.DefaultPassword)
	return command.RunErr{}
}

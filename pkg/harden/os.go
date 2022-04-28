package harden

import (
	"bufio"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"github.com/weiliang-ms/easyctl/pkg/deny"
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"github.com/weiliang-ms/easyctl/pkg/util/tmplutil"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
)

// Object 加固对象
type Object struct {
	Servers     []runner.ServerInternal
	B           []byte
	LocalRun    bool
	SkipConfirm bool
	Logger      *logrus.Logger
	unitTest    bool
	HardenItem  ExternalItem
	InternalItem
}

type InternalItem struct {
	DenyPing             bool
	DenyICMPTimeStamp    bool
	DenyRebootByKeyBoard bool
	DelUnusedUser        bool
	HideOSVersion        bool
	LockKeyFile          bool
	SetTMOUT             bool
	SetPasswdPolicy      bool
	SetSysLog            bool
	PasswdErrRetryCount  int
	UnUseUsers           []string

	// ssh
	DenySSHUseDns          bool
	DenySSHAgentForwarding bool
	DenyRootLogin          bool
	ModifyServePort        int

	// sudo user
	SetSudoUser bool
	SudoUser    string
	SudoPasswd  string

	// cron
	DelCommonUserCron       bool
	SetDelZombieProcessCron bool
}

type ExternalItem struct {
	Harden struct {
		Common struct {
			DenyPing             bool     `yaml:"denyPing"`
			DenyICMPTimeStamp    bool     `yaml:"denyICMPTimeStamp"`
			DenyRebootByKeyBoard bool     `yaml:"denyRebootByKeyBoard"`
			DelUnusedUser        bool     `yaml:"delUnusedUser"`
			HideOSVersion        bool     `yaml:"hideOSVersion"`
			LockKeyFile          bool     `yaml:"lockKeyFile"`
			SetTMOUT             bool     `yaml:"setTMOUT"`
			SetPasswdPolicy      bool     `yaml:"setPasswdPolicy"`
			SetSysLog            bool     `yaml:"setSysLog"`
			PasswdErrRetryCount  int      `yaml:"passwdErrRetryCount"`
			SudoUser             SudoUser `yaml:"sudoUser"`
			UnuseUsers           []string `yaml:"unuseUsers"`
		} `yaml:"common"`
		Ssh  SSHItem `yaml:"ssh"`
		Cron struct {
			DelCommonUserCron       bool `yaml:"delCommonUserCron"`
			SetDelZombieProcessCron bool `yaml:"setDelZombieProcessCron"`
		} `yaml:"cron"`
	} `yaml:"harden"`
}

type SudoUser struct {
	SetSudoUser bool   `yaml:"setSudoUser"`
	SudoUser    string `yaml:"sudoUser"`
	SudoPasswd  string `yaml:"sudoPasswd"`
}

type SSHItem struct {
	DenyRootLogin          bool `yaml:"denyRootLogin"`
	DenySSHUseDns          bool `yaml:"denySSHUseDns"`
	DenySSHAgentForwarding bool `yaml:"denySSHAgentForwarding"`
	ModifyServePort        int  `yaml:"modifyServePort"`
}

const (
	denyPing               = "deny ping"
	denyICMPTIMESTAMP      = "deny ICMP_TIMESTAMP"
	setTMOUT               = "set TMOUT"
	hideOSVersion          = "hide system version"
	denyRebootByKeyBoard   = "deny Control-Alt-Delete"
	setUserPasswdPolicy    = "harden user password policy"
	delUnusedUser          = "del unUse user"
	setPasswdRetryCount    = "password retry count"
	denySSHUseDns          = "deny ssh's UseDns"
	denySSHAgentForwarding = "deny ssh's AgentForwarding"
	hardenSysLog           = "harden system logfile"
	delCommonUserCron      = "del unprivileged user's cron"
	delZombieProcessCron   = "set del Zombie Process Cron"
	addSudoUser            = "add sudo user"
	lockKeyFile            = "lock sensitive file"
)

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
	DenySSHRootLoginShell = `
sed -i "/PermitRootLogin/d" /etc/ssh/sshd_config
echo "PermitRootLogin no" >> /etc/ssh/sshd_config
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
chmod -R 700 /etc/rc.d/init.d/*
chmod 644 /etc/{passwd,group}
chmod 400 /etc/shadow
chattr +i /etc/services || true
`
)

type task func() command.RunErr

func OS(item command.OperationItem) command.RunErr {
	obj := &Object{}
	obj.B = item.B
	obj.LocalRun = item.LocalRun
	obj.Logger = item.Logger
	obj.unitTest = item.UnitTest
	obj.SkipConfirm = item.SkipConfirm

	tasks := []task{
		obj.parse,
		obj.deepCopy,
		obj.confirm,
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
		obj.modifySSHPort,
		obj.denySSHRootLogin,
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

// 解析配置
func (object *Object) parse() command.RunErr {
	var item ExternalItem
	err := yaml.Unmarshal(object.B, &item)
	if err != nil {
		return command.RunErr{Err: err}
	}
	// 深拷贝？
	object.HardenItem = item
	return command.RunErr{}
}

func (object *Object) deepCopy() command.RunErr {
	object.DenyPing = object.HardenItem.Harden.Common.DenyPing
	object.DenyICMPTimeStamp = object.HardenItem.Harden.Common.DenyICMPTimeStamp
	object.DenyRebootByKeyBoard = object.HardenItem.Harden.Common.DenyRebootByKeyBoard
	object.DelUnusedUser = object.HardenItem.Harden.Common.DelUnusedUser
	object.HideOSVersion = object.HardenItem.Harden.Common.HideOSVersion
	object.LockKeyFile = object.HardenItem.Harden.Common.LockKeyFile
	object.SetTMOUT = object.HardenItem.Harden.Common.SetTMOUT
	object.SetPasswdPolicy = object.HardenItem.Harden.Common.SetPasswdPolicy
	object.SetSysLog = object.HardenItem.Harden.Common.SetSysLog
	object.PasswdErrRetryCount = object.HardenItem.Harden.Common.PasswdErrRetryCount
	object.UnUseUsers = object.HardenItem.Harden.Common.UnuseUsers

	// ssh
	object.DenySSHUseDns = object.HardenItem.Harden.Ssh.DenySSHUseDns
	object.DenySSHAgentForwarding = object.HardenItem.Harden.Ssh.DenySSHAgentForwarding
	object.DenyRootLogin = object.HardenItem.Harden.Ssh.DenyRootLogin
	object.ModifyServePort = object.HardenItem.Harden.Ssh.ModifyServePort

	// cron
	object.DelCommonUserCron = object.HardenItem.Harden.Cron.DelCommonUserCron
	object.SetDelZombieProcessCron = object.HardenItem.Harden.Cron.SetDelZombieProcessCron

	// sudo user
	object.SetSudoUser = object.HardenItem.Harden.Common.SudoUser.SetSudoUser
	object.SudoUser = object.HardenItem.Harden.Common.SudoUser.SudoUser
	object.SudoPasswd = object.HardenItem.Harden.Common.SudoUser.SudoPasswd

	return command.RunErr{}
}

func (object *Object) confirm() command.RunErr {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"rule", "open|close|value"})
	//table.SetHeader([]string{"加固规则", "是否启用"})
	//table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.SetAlignment(tablewriter.ALIGN_CENTER)

	var data [][]string
	data = append(data, []string{denyPing, convert(object.DenyPing)})
	data = append(data, []string{denyICMPTIMESTAMP, convert(object.DenyICMPTimeStamp)})
	data = append(data, []string{setTMOUT, convert(object.SetTMOUT)})
	data = append(data, []string{hideOSVersion, convert(object.HideOSVersion)})
	data = append(data, []string{denyRebootByKeyBoard, convert(object.DenyRebootByKeyBoard)})
	data = append(data, []string{setUserPasswdPolicy, convert(object.SetPasswdPolicy)})
	data = append(data, []string{delUnusedUser, convert(object.DelUnusedUser)})
	data = append(data, []string{setPasswdRetryCount, convert(object.PasswdErrRetryCount)})
	data = append(data, []string{denySSHUseDns, convert(object.DenySSHUseDns)})
	data = append(data, []string{denySSHAgentForwarding, convert(object.DenySSHAgentForwarding)})
	data = append(data, []string{hardenSysLog, convert(object.SetSysLog)})
	data = append(data, []string{delCommonUserCron, convert(object.DelCommonUserCron)})
	data = append(data, []string{delZombieProcessCron, convert(object.SetDelZombieProcessCron)})
	data = append(data, []string{addSudoUser, convert(object.SetSudoUser)})
	data = append(data, []string{lockKeyFile, convert(object.LockKeyFile)})

	table.AppendBulk(data) // Add Bulk Data
	table.Render()

	if !object.SkipConfirm {
		reader := bufio.NewReader(os.Stdin)
		input, err := confirm(reader)
		if err != nil {
			return command.RunErr{Err: err}
		}
		if input == "no" {
			os.Exit(0)
		}
	}

	object.Logger.Infoln("开始执行加固...")
	return command.RunErr{}
}

func confirm(reader *bufio.Reader) (string, error) {
	for {
		fmt.Printf("Are you sure to harden above rule? [yes/no]: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		input = strings.TrimSpace(input)

		if input != "" && (input == "yes" || input == "no") {
			return input, nil
		}
	}
}

func convert(obj interface{}) string {
	if obj == true {
		return "√"
	} else if obj == false {
		return "x"
	}

	switch value := obj.(type) {
	case int:
		return fmt.Sprintf("%d", value)
	}

	return ""
}

// 禁Ping
func (object *Object) denyPing() command.RunErr {
	object.Logger.Info("[step 1] 禁ping")
	if !object.DenyPing {
		object.Logger.Info("[skip] 未检测到禁ping配置...")
		return command.RunErr{}
	}

	return object.run(deny.DenyPingShell)
}

// 关闭ICMP_TIMESTAMP应答
func (object *Object) denyICMPTimeStamp() command.RunErr {
	object.Logger.Info("[step 2] 关闭ICMP_TIMESTAMP应答")
	if !object.DenyICMPTimeStamp {
		object.Logger.Info("[skip] 未检测到关闭ICMP_TIMESTAMP应答配置...")
		return command.RunErr{}
	}
	return object.run(UnsetICMPTimeStampShell)
}

// 设置系统空闲等待时间
func (object *Object) setTMOUT() command.RunErr {
	object.Logger.Info("[step 3] 设置shell空闲等待时间")
	if !object.SetTMOUT {
		object.Logger.Info("[skip] 未检测到`设置shell空闲等待时间`配置...")
		return command.RunErr{}
	}
	return object.run(SetTMOUTShell)
}

// 隐藏系统版本信息
func (object *Object) hideOSVersion() command.RunErr {
	object.Logger.Info("[step 4] 隐藏系统版本信息")
	if !object.HideOSVersion {
		object.Logger.Info("[skip] 未检测到`隐藏系统版本信息`配置...")
		return command.RunErr{}
	}
	return object.run(HideOSVersionShell)
}

// 禁止Control-Alt-Delete 键盘重启系统命令
func (object *Object) denyRebootByKeyBoard() command.RunErr {
	object.Logger.Info("[step 5] 禁止Control-Alt-Delete 键盘重启系统命令")
	if !object.DenyRebootByKeyBoard {
		object.Logger.Info("[skip] 未检测到`禁止Control-Alt-Delete 键盘重启系统命令`配置...")
		return command.RunErr{}
	}
	return object.run(UnsetRebootByKeyBoardShell)
}

// ssh用户密码加固
func (object *Object) setPasswdPolicy() command.RunErr {
	object.Logger.Info("[step 6] ssh用户密码加固")
	if !object.SetPasswdPolicy {
		object.Logger.Info("[skip] 未检测到`ssh用户密码加固`配置...")
		return command.RunErr{}
	}
	return object.run(SetPasswdPolicyShell)
}

// 删除系统默认用户
func (object *Object) delUnusedUser() command.RunErr {
	object.Logger.Info("[step 7] 删除系统默认用户")
	if !object.DelUnusedUser {
		object.Logger.Info("[skip] 未检测到`删除系统默认用户`配置...")
		return command.RunErr{}
	}

	cmd := tmplutil.RenderPanicErr(delUserListTmpl, tmplutil.TmplRenderData{
		"UserList": object.UnUseUsers,
	})

	return object.run(cmd)
}

// 修改允许密码错误次数
func (object *Object) errPasswdRetryCount() command.RunErr {
	object.Logger.Info("[step 8] 修改允许密码错误次数")
	if object.PasswdErrRetryCount < 0 || object.PasswdErrRetryCount > 3 {
		object.Logger.Info("[skip] 修改允许密码错误次数值非法（合法区间 [0,3]）...")
		return command.RunErr{}
	}

	cmd := tmplutil.RenderPanicErr(SetErrPasswdRetryShellTmpl, tmplutil.TmplRenderData{
		"RetryCount": object.PasswdErrRetryCount,
	})

	return object.run(cmd)
}

// 关闭ssh UseDNS
func (object *Object) denySSHUseDns() command.RunErr {
	object.Logger.Info("[step 9] ssh关闭UseDNS")
	if !object.DenySSHUseDns {
		object.Logger.Info("[skip] 未检测到`ssh关闭UseDNS`配置...")
		return command.RunErr{}
	}
	return object.run(DenySSHUseDnsShell)
}

// 关闭ssh `AgentForwarding`和`TcpForwarding`
func (object *Object) denySSHAgentForwarding() command.RunErr {
	object.Logger.Info("[step 10] ssh关闭AgentForwarding")
	if !object.DenySSHAgentForwarding {
		object.Logger.Info("[skip] 未检测到`ssh关闭AgentForwarding`配置...")
		return command.RunErr{}
	}
	return object.run(DenySSHdAgentForwardingShell)
}

// 加固系统日志文件
func (object *Object) setSysLog() command.RunErr {
	object.Logger.Info("[step 11] 加固系统日志文件")
	if !object.SetSysLog {
		object.Logger.Info("[skip] 未检测到`加固系统日志文件`配置...")
		return command.RunErr{}
	}
	return object.run(HardenSystemLogShell)
}

// 删除非root用户定时任务
func (object *Object) delCommonUserCron() command.RunErr {
	object.Logger.Info("[step 12] 删除非root用户定时任务")
	if !object.DelCommonUserCron {
		object.Logger.Info("[skip] 未检测到`删除非root用户定时任务`配置...")
		return command.RunErr{}
	}
	return object.run(DeleteCommonUserCronShell)
}

// 定时清理僵尸进程
func (object *Object) delZombieProcessCron() command.RunErr {
	object.Logger.Info("[step 13] 定时清理僵尸进程")
	if !object.SetDelZombieProcessCron {
		object.Logger.Info("[skip] 未检测到`定时清理僵尸进程`配置...")
		return command.RunErr{}
	}
	return object.run(DeleteZombieProcessCronShell)
}

// 创建sudo用户
func (object *Object) addSudoUser() command.RunErr {
	// todo: 重构调用内部指令
	if !object.SetSudoUser {
		object.Logger.Info("[skip] 未检测到`添加sudo用户`配置...")
		return command.RunErr{}
	}
	object.Logger.Infof("[step 14] 添加sudo用户: %s 密码: %s",
		object.SudoUser, object.SudoPasswd)

	cmd, _ := tmplutil.Render(addSudoUserTmpl, tmplutil.TmplRenderData{
		"User":     object.SudoUser,
		"Password": object.SudoPasswd,
	})

	return object.run(cmd)
}

// 锁定敏感文件并降权
func (object *Object) lockKeyFile() command.RunErr {
	object.Logger.Info("[step 15] 锁定敏感文件")
	return object.run(LockKeyFileShell)
}

// 修改ssh port & 禁止root登录
func (object *Object) modifySSHPort() command.RunErr {

	ok := object.ModifyServePort != 22 || object.ModifyServePort > 65535 || object.ModifyServePort < 1024

	if !ok {
		object.Logger.Info("[skip] 未检测到`修改ssh端口`配置或ssh端口非法，合法区间为(1024,65535)...")
		return command.RunErr{}
	}

	object.Logger.Infof("[step 16] 调整ssh登录端口为: %d，禁止root直接登录.", object.ModifyServePort)
	cmd := tmplutil.RenderPanicErr(modifySSHLoginShellTmpl, tmplutil.TmplRenderData{
		"Port": object.ModifyServePort,
	})
	return object.run(cmd)
}

func (object *Object) denySSHRootLogin() command.RunErr {
	object.Logger.Info("[step 16] 禁止root直接登录.")
	if !object.DenyRootLogin {
		object.Logger.Info("[skip] 未检测到`禁止root直接登录`配置...")
		return command.RunErr{}
	}

	if object.DenyRootLogin && !object.SetSudoUser {
		object.Logger.Info("[skip] 禁止root登录情况下，必须创建创建sudo用户以免无法登录...")
		return command.RunErr{}
	}
	return object.run(DenySSHRootLoginShell)
}

func (object *Object) print() command.RunErr {
	object.Logger.Infof("[done] 安全加固完毕!")
	return command.RunErr{}
}

func (object *Object) run(cmd string) command.RunErr {

	if object.LocalRun {
		return command.RunErr{Err: runner.LocalRun(cmd, object.Logger).Err}
	}

	return runner.RemoteRun(runner.RemoteRunItem{
		ManifestContent:     object.B,
		Logger:              object.Logger,
		Cmd:                 cmd,
		RecordErrServerList: false,
	})
}

package upgrade

import (
	"github.com/lithammer/dedent"
	"text/template"
)

var (
	sshUpgradeShellTmpl = template.Must(template.New("sshUpgradeShellTmpl").Parse(
		dedent.Dedent(`
echo "启动telnet server..."

systemctl restart telnet.socket
systemctl restart xinetd
echo 'pts/0' >>/etc/securetty
echo 'pts/1' >>/etc/securetty
systemctl restart telnet.socket

mv /etc/ssh/ /etc/ssh-bak

tar zxvf {{ .FilePath }} -C /tmp

cd /tmp/openssh*

./configure --prefix=/usr --sysconfdir=/etc/ssh --with-md5-passwords {{if .OpensslPath }}--with-ssl-dir={{.OpensslPath}}{{end}}

{{if not .OpensslPath }}
./configure --prefix=/usr --sysconfdir=/etc/ssh --with-md5-passwords
{{end}}

make -j $(nproc) && make install

\cp contrib/redhat/sshd.init /etc/init.d/sshd
\chkconfig sshd on

ssh -V

cat > /etc/ssh/sshd_config <<EOF
Protocol 2
SyslogFacility AUTHPRIV
PermitRootLogin yes
PasswordAuthentication yes
ChallengeResponseAuthentication no
PermitRootLogin yes
PubkeyAuthentication yes
UsePAM yes
UseDNS no
AcceptEnv LANG LC_CTYPE LC_NUMERIC LC_TIME LC_COLLATE LC_MONETARY LC_MESSAGES
AcceptEnv LC_PAPER LC_NAME LC_ADDRESS LC_TELEPHONE LC_MEASUREMENT
AcceptEnv LC_IDENTIFICATION LC_ALL LANGUAGE
AcceptEnv XMODIFIERS
AllowTcpForwarding yes
X11Forwarding yes
Subsystem sftp /usr/libexec/openssh/sftp-server
EOF

sed -i "s;Type=notify;#Type=notify;g" /usr/lib/systemd/system/sshd.service
systemctl daemon-reload && systemctl restart sshd

if [ $? -eq 0 ];then
	systemctl disable telnet.socket --now && systemctl disable xinetd --now
fi
   `)))
)

// 下载安装介质->解析列表->检测依赖是否安装->检测yum可用性->尝试安装依赖->

//func (ac *Actuator) OpenSSH() {
//	ac.DependenciesList = []string{
//		"telnet-server",
//		"telnet",
//		"xinetd",
//		"pam-devel",
//		"zlib-devel",
//		"openssl-devel",
//	}
//	ac.download().parseServerList().detect().handoutFile().compileOpenSSHCmd().execute("编译安装openssh", 0)
//}
//
//func (ac *Actuator) compileOpenSSHCmd() *Actuator {
//	content, _ := tmplutil.Render(sshUpgradeShellTmpl, util.Data{
//		"FilePath":    ac.FilePath,
//		"OpensslPath": ac.OpensslDir,
//	})
//	ac.Cmd = content
//	return ac
//}

package upgrade

import (
	"github.com/lithammer/dedent"
	"text/template"
)

var (
	opensslUpgradeShellTmpl = template.Must(template.New("opensslUpgradeShellTmpl").Parse(
		dedent.Dedent(`
tar zxvf {{ .FilePath }} -C /tmp
cd /tmp/openssl-OpenSSL*
./config shared --openssldir=/usr/local/openssl --prefix=/usr/local/openssl
make -j $(nproc) && make install
sed -i '/\/usr\/local\/openssl\/lib/d' /etc/ld.so.conf
echo "/usr/local/openssl/lib" >> /etc/ld.so.conf
ldconfig -v
mv /usr/bin/openssl /usr/bin/openssl.old
ln -s /usr/local/openssl/bin/openssl /usr/bin/openssl
openssl version
   `)))
)

// 下载安装介质->解析列表->检测依赖是否安装->检测yum可用性->尝试安装依赖->

//func (ac *Actuator) Openssl() {
//	ac.DependenciesList = []string{"gcc", "perl"}
//	ac.download().parseServerList().detect().handoutFile().compileOpensslCmd().execute("编译安装openssl", 0)
//}
//
//func (ac *Actuator) compileOpensslCmd() *Actuator {
//	content, _ := util.Render(opensslUpgradeShellTmpl, util.Data{
//		"FilePath": ac.FilePath,
//	})
//	ac.Cmd = content
//	return ac
//}

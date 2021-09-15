package set

import (
	"fmt"
	"github.com/lithammer/dedent"
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"github.com/weiliang-ms/easyctl/pkg/util"
	"log"
	"text/template"
)

var baseRepoConfigTmpl = template.Must(template.New("baseRepoConfigTmpl").Parse(dedent.Dedent(
	`
mkdir -p /etc/yum.repos.d/bak

files=$(ls /etc/yum.repos.d/*.repo 2>/dev/null)

if [ $? -eq 0 ];then
	for i in $files;do
		mv $i /etc/yum.repos.d/bak
	done
fi

cat > /etc/yum.repos.d/base.repo <<EOF
[base]
name=CentOS-\$releasever - Base
baseurl={{ .RepoUrl }}/centos/\$releasever/os/\$basearch/
gpgcheck=1
gpgkey={{ .RepoUrl }}/centos/RPM-GPG-KEY-CentOS-7

[updates]
name=CentOS-\$releasever - Updates
baseurl={{ .RepoUrl }}/centos/\$releasever/updates/\$basearch/
gpgcheck=1
gpgkey={{ .RepoUrl }}/centos/RPM-GPG-KEY-CentOS-7

[extras]
name=CentOS-\$releasever - Extras
baseurl={{ .RepoUrl }}/centos/\$releasever/extras/\$basearch/
gpgcheck=1
gpgkey={{ .RepoUrl }}/centos/RPM-GPG-KEY-CentOS-7

[centosplus]
name=CentOS-\$releasever - Plus
baseurl={{ .RepoUrl }}/centos/\$releasever/centosplus/$basearch/
gpgcheck=1
enabled=0
gpgkey={{ .RepoUrl }}/centos/RPM-GPG-KEY-CentOS-7
EOF
`)))

var localRepoConfigTmpl = template.Must(template.New("localRepoConfigTmpl").Parse(dedent.Dedent(
	`
{{- if .ISOPath }}
mkdir -p /yum
mount -o loop {{ .ISOPath }} /yum

mkdir -p /etc/yum.repos.d/bak

files=$(ls /etc/yum.repos.d/*.repo 2>/dev/null)

if [ $? -eq 0 ];then
	for i in $files;do
		mv $i /etc/yum.repos.d/bak
	done
fi

cat >> /etc/yum.repos.d/c7.repo <<EOF
[c7repo]
name=c7repo
baseurl=file:///yum
enabled=1
gpgcheck=0
gpgkey=file:///etc/pki/rpm-gpg/RPM-GPG-KEY-CentOS-7
EOF
{{- end}}
`)))

func (ac *Actuator) SetYumRepo() {
	ac.parseServerList().download().handoutIso().setRepoCmd().execute(fmt.Sprintf("配置yum仓库地址为：%s", ac.Value))
}

func (ac *Actuator) setRepoCmd() *Actuator {
	if ac.FilePath == "" {
		ac.Cmd = repoSetCmd(ac.Value)
	} else {
		ac.Cmd = localRepoSetCmd(ac.FilePath)
		ac.Value = "file:///yum"
	}
	return ac
}

func repoSetCmd(repoUrl string) string {
	content, err := util.Render(baseRepoConfigTmpl, util.Data{
		"RepoUrl": repoUrl,
	})
	if err != nil {
		log.Fatalf("解析模板失败：%s", err.Error())
	}

	return content
}

func localRepoSetCmd(isoPath string) string {
	content, err := util.Render(localRepoConfigTmpl, util.Data{
		"ISOPath": isoPath,
	})
	if err != nil {
		log.Fatalf("解析模板失败：%s", err.Error())
	}

	return content
}

// 分发iso文件
func (ac *Actuator) handoutIso() *Actuator {
	if len(ac.ServerList.Server) <= 0 {
		return ac
	}
	localIsoPath := ac.FilePath
	ac.FilePath = fmt.Sprintf("/tmp/%s", ac.FilePath)
	for _, v := range ac.ServerList.Server {
		runner.ScpFile(localIsoPath, ac.FilePath, v, 0644)
	}
	return ac
}

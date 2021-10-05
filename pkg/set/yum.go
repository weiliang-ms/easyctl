package set

import (
	"github.com/lithammer/dedent"
	"github.com/weiliang-ms/easyctl/pkg/util/tmplutil"
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

// YUM 实体
type YUM struct {
}

// Config 配置yum源
func (yum YUM) Config(b []byte, debug bool) error {
	return nil
}
func repoSetCmd(repoUrl string) string {
	content, err := tmplutil.Render(baseRepoConfigTmpl, tmplutil.TmplRenderData{
		"RepoUrl": repoUrl,
	})
	if err != nil {
		log.Fatalf("解析模板失败：%s", err.Error())
	}

	return content
}

func localRepoSetCmd(isoPath string) string {
	content, err := tmplutil.Render(localRepoConfigTmpl, tmplutil.TmplRenderData{
		"ISOPath": isoPath,
	})
	if err != nil {
		log.Fatalf("解析模板失败：%s", err.Error())
	}

	return content
}

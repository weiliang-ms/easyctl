package harden

import (
	"github.com/lithammer/dedent"
	"text/template"
)

var delUserListTmpl = template.Must(template.New("delUserListTmpl").Parse(dedent.Dedent(
	`{{- if .UserList }}
{{- range .UserList }}
userdel {{ . }} &>/dev/null || true
{{- end }}
{{- end}}

{{- if .UserList }}
{{- range .UserList }}
userdel {{ . }} &>/dev/null || true
{{- end }}
{{- end}}`)))

var SetErrPasswdRetryShellTmpl = template.Must(template.New("delUserListTmpl").Parse(dedent.Dedent(
	`{{- if .RetryCount }}
sed -i "/MaxAuthTries/d" /etc/ssh/sshd_config
echo "MaxAuthTries {{ .RetryCount }}" >> /etc/ssh/sshd_config
service sshd restart
{{- end }}
`)))

var addSudoUserTmpl = template.Must(template.New("addSudoUserTmpl").Parse(dedent.Dedent(`
{{- if .User }}
chattr -i /etc/passwd /etc/shadow /etc/group /etc/gshadow /etc/inittab
useradd -m {{ .User }} &>/dev/null || true
echo {{ .Password }} | passwd --stdin {{ .User }} || true
sed -i '/{{ .User }}/d' /etc/sudoers
echo "{{ .User }}        ALL=(ALL)       NOPASSWD: ALL" >> /etc/sudoers
{{- end }}
`)))

var modifySSHLoginShellTmpl = template.Must(template.New("modifySSHLoginShellTmpl").Parse(dedent.Dedent(`
{{- if .Port }}

sed -i "/Port 22/d" /etc/ssh/sshd_config
echo "Port {{ .Port }} " >> /etc/ssh/sshd_config

setenforce 0
firewall-cmd --zone=public --add-port={{ .Port }}/tcp --permanent || true
firewall-cmd --reload || true

service sshd restart
{{- end }}
`)))

package add

import (
	"github.com/lithammer/dedent"
	"github.com/weiliang-ms/easyctl/pkg/util"
	"text/template"
)

var addUserTmpl = template.Must(template.New("addUserTmpl").Parse(
	dedent.Dedent(`
{{- if not .NoLogin }}
useradd {{ .User }}
{{- if .Password }}
if [ $? -eq 0 ];then
  echo {{ .Password }} | passwd --stdin {{ .User }}
fi
{{- end}}
{{- end}}

{{- if .NoLogin }}
groupadd {{ .User }}
useradd {{ .User }} -g {{ .User }} -s /sbin/nologin -M
{{- end}}

`)))

func (ac *Actuator) AddUser() {
	ac.parseServerList().addUserCmd().execute("新增用户")
}

func (ac *Actuator) addUserCmd() *Actuator {
	content, _ := util.Render(addUserTmpl, util.Data{
		"NoLogin":  ac.NoLogin,
		"User":     ac.UserName,
		"Password": ac.Password,
	})
	ac.Cmd = content
	return ac
}

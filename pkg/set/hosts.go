package set

import (
	"fmt"
	"github.com/lithammer/dedent"
	"github.com/weiliang-ms/easyctl/pkg/util"
	"os"
	"strings"
	"text/template"
)

type hostResolvInstance struct {
	Address  string
	Hostname string
}

var setHostsShellTmpl = template.Must(template.New("").Parse(
	dedent.Dedent(`
{{- if .HostResolveList }}
sed -i ':a;$!{N;ba};s@# easyctl hosts BEGIN.*# easyctl hosts END@@' /etc/hosts
sed -i '/^$/N;/\n$/N;//D' /etc/hosts

cat >>/etc/hosts<<EOF   
# easyctl hosts BEGIN
{{- range .HostResolveList }}
{{ . }}
{{- end }}
# easyctl hosts END
EOF
{{- end }}
`)))

func (ac *Actuator) HostResolve() {
	ac.parseServerList().setHostResolveCmd().execute("配置hosts解析")
}

func (ac *Actuator) setHostResolveCmd() *Actuator {
	var slice []hostResolvInstance
	if len(ac.ServerList.Server) <= 0 {
		hostname, _ := os.Hostname()
		instance := hostResolvInstance{
			Address:  "127.0.0.1",
			Hostname: hostname,
		}
		slice = append(slice, instance)
	} else {
		for _, s := range ac.ServerList.Server {
			var instance hostResolvInstance
			re := s.RemoteShell("hostname")
			if re.Code == 0 && re.StdOut != "" {
				instance.Hostname = strings.Trim(re.StdOut, "\n")
				instance.Address = s.Host
				slice = append(slice, instance)
			}
		}
	}
	ac.Cmd = ac.packageHostResolve(slice)
	return ac
}

func (ac *Actuator) packageHostResolve(instances []hostResolvInstance) string {
	var list []string
	for _, v := range instances {
		fmt.Println(v.Hostname)
		if v.Hostname != localhost && v.Hostname != "localhost.localdomain" {
			list = append(list, fmt.Sprintf("%s %s", v.Address, v.Hostname))
		}
	}
	content, _ := util.Render(setHostsShellTmpl, util.Data{
		"HostResolveList": list,
	})

	return content
}

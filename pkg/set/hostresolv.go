package set

import (
	"fmt"
	"github.com/lithammer/dedent"
	"github.com/sirupsen/logrus"
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"github.com/weiliang-ms/easyctl/pkg/util/tmplutil"
	"strings"
	"text/template"
)

// IPAddress ip地址
type IPAddress []string

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

// e.g: 1.1.1.1 server-ddd
type getHostResolveFunc func(b []byte, logger *logrus.Logger, cmd string) ([]runner.ShellResult, error)

const GetHostResolveFunc = "getHostResolveFunc"

// HostResolve 配置主机host解析
func HostResolve(item command.OperationItem) command.RunErr {

	resolveFnc, ok := item.OptionFunc[GetHostResolveFunc].(func(b []byte, logger *logrus.Logger, cmd string) ([]runner.ShellResult, error))
	if !ok {
		return command.RunErr{Err: fmt.Errorf("入参：%s 非法", GetHostResolveFunc)}
	}

	results, err := resolveFnc(item.B, item.Logger, "hostname")

	if err != nil {
		return command.RunErr{Err: err}
	}

	// todo: IP地址排序
	hosts := map[string]string{}
	for _, v := range results {
		if v.StdOut != "localhost" {
			hosts[v.StdOut] = v.Host
		}
	}

	var addresses []string
	for k, v := range hosts {
		addresses = append(addresses, strings.TrimSuffix(fmt.Sprintf("%s %s", v, k), "\n"))
	}

	shell, _ := tmplutil.Render(setHostsShellTmpl, tmplutil.TmplRenderData{
		"HostResolveList": addresses,
	})

	return runner.RemoteRun(runner.RemoteRunItem{
		ManifestContent:     item.B,
		Logger:              item.Logger,
		Cmd:                 shell,
		RecordErrServerList: false,
	})

}

func GetHostResolve(b []byte, logger *logrus.Logger, cmd string) ([]runner.ShellResult, error) {
	return runner.GetResult(runner.RemoteRunItem{
		ManifestContent: b,
		Logger:          logger,
		Cmd:             cmd,
	})
}

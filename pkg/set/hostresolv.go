package set

import (
	"fmt"
	"github.com/lithammer/dedent"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"github.com/weiliang-ms/easyctl/pkg/util/tmplutil"
	"strconv"
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

// HostResolve 配置主机host解析
func HostResolve(item command.OperationItem) error {

	results, err := GetResult(item.B, item.Logger, "hostname")
	if err != nil {
		return err
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

	shell, err := tmplutil.Render(setHostsShellTmpl, tmplutil.TmplRenderData{
		"HostResolveList": addresses,
	})

	if err != nil {
		return err
	}

	return Config(item.B, item.Logger, shell)
}

func (addresses IPAddress) Len() int { return len(addresses) }

func (addresses IPAddress) Swap(i, j int) { addresses[i], addresses[j] = addresses[j], addresses[i] }

func (addresses IPAddress) Less(i, j int) bool {

	address1 := strings.Split(addresses[i], ".")
	address2 := strings.Split(addresses[j], ".")

	for k := 0; k < 4; k++ {
		if address1[k] != address2[k] {
			num1, _ := strconv.Atoi(address1[k])
			num2, _ := strconv.Atoi(address2[k])
			return num1 < num2
		}
	}

	return true
}

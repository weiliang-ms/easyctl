package set

import (
	"fmt"
	"github.com/lithammer/dedent"
	"github.com/weiliang-ms/easyctl/pkg/util"
	"log"
	"net"
	"strings"
	"text/template"
)

var setDnsShellTmpl = template.Must(template.New("").Parse(
	dedent.Dedent(`
{{- range .DnsServerList }}
sed -i "/{{ . }}/d" /etc/resolv.conf
echo nameserver {{ . }} >> /etc/resolv.conf
{{- end }}
`)))

func (ac *Actuator) DNS() {
	ac.parseServerList().setDnsCmd().execute(fmt.Sprintf("配置dns地址：%s", ac.Value))
}

func (ac *Actuator) setDnsCmd() *Actuator {
	var servers []string
	if strings.Contains(ac.Value, ",") {
		for _, v := range strings.Split(ac.Value, ",") {
			if net.ParseIP(v) != nil {
				servers = append(servers, v)
			} else {
				log.Printf("dns地址：%s 非法，已忽略...\n", v)
			}
		}
	} else {
		if net.ParseIP(ac.Value) != nil {
			servers = append(servers, ac.Value)
		} else {
			log.Fatalf("dns地址：%s 非法，请检测...\n", ac.Value)
		}
	}

	ac.Cmd, _ = util.Render(setDnsShellTmpl, util.Data{
		"DnsServerList": servers,
	})

	return ac
}

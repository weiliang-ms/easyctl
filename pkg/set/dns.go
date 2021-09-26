package set

import (
	"errors"
	"fmt"
	"github.com/lithammer/dedent"
	"github.com/weiliang-ms/easyctl/pkg/util"
	"gopkg.in/yaml.v2"
	"net"
	"text/template"
)

var setDnsShellTmpl = template.Must(template.New("").Parse(
	dedent.Dedent(`
{{- range .DnsServerList }}
sed -i "/{{ . }}/d" /etc/resolv.conf
echo nameserver {{ . }} >> /etc/resolv.conf
{{- end }}
`)))

type DnsConfig struct {
	DnsList []string `yaml:"dns"`
}

func Dns(b []byte, debug bool) error {
	script, err := AddDnsScript(b, setDnsShellTmpl)
	if err != nil {
		return err
	}
	return Config(b, debug, script)
}

func AddDnsScript(b []byte, tmpl *template.Template) (string, error) {

	p, err := ParseDnsConfig(b)
	if err != nil {
		return "", err
	}

	ok , err := IsValid(p.DnsList)
	if !ok {
		return "", err
	}

	return util.Render(tmpl, map[string]interface{}{
		"DnsServerList": p.DnsList,
	})
}

func ParseDnsConfig(b []byte) (DnsConfig, error) {
	config := DnsConfig{}
	if err := yaml.Unmarshal(b, &config); err != nil {
		return DnsConfig{}, err
	} else {
		return config, nil
	}
}

func IsValid(dnsList []string) (bool, error) {
	for _ , v := range dnsList {
		if ok := net.ParseIP(v);ok == nil {
			return false, errors.New(fmt.Sprintf("%s地址非法", v))
		}
	}
	// todo: 可达性检测
	return true, nil
}
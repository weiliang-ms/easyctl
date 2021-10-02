package set

import (
	"fmt"
	"github.com/lithammer/dedent"
	"github.com/sirupsen/logrus"
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

// DnsConfig 添加dns对象属性
type DnsConfig struct {
	DnsList []string `yaml:"dns"`
}

// Dns 设置Dns
func Dns(b []byte, logger *logrus.Logger) error {
	script, err := AddDnsScript(b, setDnsShellTmpl)
	if err != nil {
		return err
	}
	return Config(b, logger, script)
}

// AddDnsScript 获取添加dns脚本
func AddDnsScript(b []byte, tmpl *template.Template) (string, error) {

	p, err := ParseDnsConfig(b)
	if err != nil {
		return "", err
	}

	ok, err := IsValid(p.DnsList)
	if !ok {
		return "", err
	}

	return util.Render(tmpl, map[string]interface{}{
		"DnsServerList": p.DnsList,
	})
}

// ParseDnsConfig 解析dns配置
func ParseDnsConfig(b []byte) (DnsConfig, error) {
	config := DnsConfig{}
	if err := yaml.Unmarshal(b, &config); err != nil {
		return DnsConfig{}, err
	}
	return config, nil
}

// IsValid 判断Dns合法性
func IsValid(dnsList []string) (bool, error) {
	for _, v := range dnsList {
		if ok := net.ParseIP(v); ok == nil {
			return false, fmt.Errorf("%s地址非法", v)
		}
	}
	// todo: 可达性检测
	return true, nil
}

package clean

import (
	"github.com/lithammer/dedent"
	"github.com/sirupsen/logrus"
	"github.com/weiliang-ms/easyctl/pkg/util"
	"github.com/weiliang-ms/easyctl/pkg/util/slice"
	"gopkg.in/yaml.v2"
	"text/template"
)

var pruneDnsShellTmpl = template.Must(template.New("").Parse(
	dedent.Dedent(`
{{- range .DnsServerList }}
sed -i "/{{ . }}/d" /etc/resolv.conf
{{- end }}

{{if .ALL}}
	{{if .PreserveList}}
		{{- range .PreserveList }}
		while read line
		do
		   if [[ "$line" =~ {{ . }} ]];then
			 echo 过滤...
		   else
			 sed -i "/$line/d" /etc/resolv.conf
		   fi
		done < /etc/resolv.conf
		{{- end }}
	{{else}}
		echo "" > /etc/resolv.conf
	{{end}}
{{end}}
`)))

type DnsCleanerConfig struct {
	CleanDns struct {
		AddressList []string `yaml:"address-list"`
		Excludes    []string `yaml:"excludes"`
	} `yaml:"clean-dns"`
}

func Dns(b []byte, logger *logrus.Logger) error {
	script, err := PruneDnsScript(b, pruneDnsShellTmpl)
	if err != nil {
		return err
	}
	return Config(b, logger, script)
}

func PruneDnsScript(b []byte, tmpl *template.Template) (string, error) {

	config, err := ParseDnsConfig(b)
	if err != nil {
		return "", err
	}

	var cleanList , preserveList []string
	// AddressList非空
	if len(config.CleanDns.AddressList) != 0 {
		cleanList = slice.StringSliceRemove(config.CleanDns.AddressList, config.CleanDns.Excludes)
	}else {
		preserveList = config.CleanDns.Excludes
	}
	return util.Render(tmpl, map[string]interface{}{
		"DnsServerList": cleanList,
		"ALL": len(cleanList)==0,
		"PreserveList": preserveList,
	})
}

func ParseDnsConfig(b []byte) (DnsCleanerConfig, error) {
	config := DnsCleanerConfig{}
	if err := yaml.Unmarshal(b, &config); err != nil {
		return DnsCleanerConfig{}, err
	} else {
		return config, nil
	}
}
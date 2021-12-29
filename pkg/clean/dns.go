package clean

import (
	"github.com/lithammer/dedent"
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"github.com/weiliang-ms/easyctl/pkg/util/slice"
	"github.com/weiliang-ms/easyctl/pkg/util/tmplutil"
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

// DnsCleanerConfig 清理dns实体
type DnsCleanerConfig struct {
	CleanDns struct {
		AddressList []string `yaml:"address-list"`
		Excludes    []string `yaml:"excludes"`
	} `yaml:"clean-dns"`
}

// Dns 清理dns
func Dns(item command.OperationItem) command.RunErr {
	script, err := PruneDnsScript(item.B, pruneDnsShellTmpl)
	if err != nil {
		return command.RunErr{Err: err}
	}
	return runner.RemoteRun(runner.RemoteRunItem{
		B:      item.B,
		Logger: item.Logger,
		Cmd:    script,
	})
}

// PruneDnsScript 清理dns脚本
func PruneDnsScript(b []byte, tmpl *template.Template) (string, error) {

	config, err := ParseDnsConfig(b)
	if err != nil {
		return "", err
	}

	var cleanList, preserveList []string
	// AddressList非空
	if len(config.CleanDns.AddressList) != 0 {
		cleanList = slice.StringSliceRemove(config.CleanDns.AddressList, config.CleanDns.Excludes)
	} else {
		preserveList = config.CleanDns.Excludes
	}
	return tmplutil.Render(tmpl, map[string]interface{}{
		"DnsServerList": cleanList,
		"ALL":           len(cleanList) == 0,
		"PreserveList":  preserveList,
	})
}

// ParseDnsConfig 解析清理dns的配置文件内容
func ParseDnsConfig(b []byte) (DnsCleanerConfig, error) {
	config := DnsCleanerConfig{}
	if err := yaml.Unmarshal(b, &config); err != nil {
		return DnsCleanerConfig{}, err
	}
	return config, nil
}

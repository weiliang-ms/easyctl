package passwordless

import (
	_ "embed"
	"errors"
	"github.com/lithammer/dedent"
	"github.com/weiliang-ms/easyctl/pkg/exec"
	"github.com/weiliang-ms/easyctl/pkg/ssh"
	"github.com/weiliang-ms/easyctl/pkg/util"
	"gopkg.in/yaml.v2"
	"k8s.io/klog"
	"os"
	"text/template"
)

// 本机互信脚本模板
var passwordLessTmpl = template.Must(template.New("passwordLessTmpl").Parse(dedent.Dedent(`
#!/bin/bash
set -e
mkdir -p ~/.ssh
tee ~/.ssh/id_rsa.pub <<EOF
{{ .PublicKey }}
EOF

tee ~/.ssh/authorized_keys <<EOF
{{ .PublicKey }}
EOF

tee ~/.ssh/id_rsa <<EOF
{{ .PrivateKey }}
EOF

chmod 600 ~/.ssh -R
`)))

//go:embed config.yaml
var config []byte

func Config(configFile string, level util.LogLevel) error {
	if configFile == "" {
		klog.Infof("检测到配置文件为空，生成配置文件样例 -> %s", util.ConfigFile)
		_ = os.WriteFile(util.ConfigFile, config, 0666)
	}

	b, err := os.ReadFile(configFile)
	if err != nil {
		return errors.New("读取配置文件失败")
	}
	executorItem, err := parseConfig(b)
	if err != nil {
		return err
	}

	if err := executorItem.Run(true, level); err != nil {
		return err
	}

	return nil
}

// ParseConfig 解析yaml配置
func parseConfig(content []byte) (exec.ExecutorItem, error) {

	executor := exec.ExecutorItem{}
	sl := exec.ServerList{}
	err := yaml.Unmarshal(content, &sl)
	if err != nil {
		return executor, err
	}

	executor.Server = sl.ParseServerList()
	script, err := MakeKeyPairScript(passwordLessTmpl)
	if err != nil {
		return executor, err
	}

	executor.Script = script

	return executor, nil
}

func MakeKeyPairScript(tmpl *template.Template) (string, error) {

	prv, pub, err := ssh.MakeSSHKeyPair()
	if err != nil {
		return "", err
	}

	return util.Render(tmpl, map[string]interface{}{
		"PublicKey":  pub,
		"PrivateKey": prv,
	})
}

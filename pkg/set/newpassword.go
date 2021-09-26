package set

import (
	_ "embed"
	"errors"
	"fmt"
	"github.com/lithammer/dedent"
	"github.com/weiliang-ms/easyctl/pkg/util"
	"gopkg.in/yaml.v2"
	"text/template"
)

// NewPasswordTmpl 本机互信脚本模板
var NewPasswordTmpl = template.Must(template.New("NewPasswordTmpl").Parse(dedent.Dedent(`
#!/bin/bash
set -e
echo "{{ .NewPassword }}" | passwd --stdin root
`)))

type PasswordConfig struct {
	Password string `yaml:"newRootPassword"`
}

func NewPassword(config []byte, debug bool) error {
	script, err := NewPasswordScript(config, NewPasswordTmpl)
	if err != nil {
		return err
	}
	return Config(config, debug, script)
}

func NewPasswordScript(b []byte, tmpl *template.Template) (string, error) {

	p, err := ParseNewPasswordConfig(b)
	if err != nil {
		return "", err
	}

	if p.Password == "" || len(p.Password) < 6 {
		return "", errors.New(fmt.Sprintf("密码长度：%d 不符合标准", len(p.Password)))
	}

	return util.Render(tmpl, map[string]interface{}{
		"NewPassword": p.Password,
	})
}

func ParseNewPasswordConfig(b []byte) (PasswordConfig, error) {
	config := PasswordConfig{}
	if err := yaml.Unmarshal(b, &config); err != nil {
		return PasswordConfig{}, err
	} else {
		return config, nil
	}
}

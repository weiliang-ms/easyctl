package set

import (
	"fmt"
	"github.com/lithammer/dedent"
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"github.com/weiliang-ms/easyctl/pkg/util/tmplutil"
	"gopkg.in/yaml.v2"
	"text/template"
)

// NewPasswordTmpl 本机互信脚本模板
var NewPasswordTmpl = template.Must(template.New("NewPasswordTmpl").Parse(dedent.Dedent(`
#!/bin/bash
set -e
echo "{{ .NewPassword }}" | passwd --stdin root
`)))

// PasswordConfig 更新密码对象属性
type PasswordConfig struct {
	Password string `yaml:"newRootPassword"`
}

//NewPassword 修改用户口令
func NewPassword(item command.OperationItem) command.RunErr {
	script, err := NewPasswordScript(item.B, NewPasswordTmpl)
	if err != nil {
		return command.RunErr{Err: err}
	}
	return command.RunErr{Err: runner.RemoteRun(item.B, item.Logger, script)}
}

// NewPasswordScript 获取修改用户口令脚本
func NewPasswordScript(b []byte, tmpl *template.Template) (string, error) {

	p, err := ParseNewPasswordConfig(b)
	if err != nil {
		return "", err
	}

	if p.Password == "" || len(p.Password) < 6 {
		return "", fmt.Errorf("密码长度：%d 不符合标准", len(p.Password))
	}

	return tmplutil.Render(tmpl, map[string]interface{}{
		"NewPassword": p.Password,
	})
}

// ParseNewPasswordConfig 解析用户口令配置文件内容
func ParseNewPasswordConfig(b []byte) (PasswordConfig, error) {
	config := PasswordConfig{}
	if err := yaml.Unmarshal(b, &config); err != nil {
		return PasswordConfig{}, err
	}

	return config, nil
}

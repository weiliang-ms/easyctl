package add

import (
	"fmt"
	"github.com/lithammer/dedent"
	"github.com/sirupsen/logrus"
	"github.com/weiliang-ms/easyctl/pkg/util"
	"gopkg.in/yaml.v2"
	"text/template"
)

var addUserTmpl = template.Must(template.New("addUserTmpl").Parse(
	dedent.Dedent(`
#!/bin/bash
{{- if not .NoLogin }}
useradd {{ .User }} {{- if not .WorkDir }} -d {{ .WorkDir }} {{- end}}
{{- end }}

{{- if .Password }}
if [ $? -eq 0 ];then
  echo {{ .Password }} | passwd --stdin {{ .User }}
fi
{{- end}}

{{- if .NoLogin }}
groupadd {{ .User }}
useradd {{ .User }} -g {{ .User }} -s /sbin/nologin -M
{{- end}}

`)))

type NewUserConfig struct {
	NewUser struct {
		Name     string `yaml:"name"`
		Nologin  bool   `yaml:"nologin"`
		Password string `yaml:"password"`
		UserDir  string `yaml:"user-dir"`
	} `yaml:"new-user"`
}

func User(b []byte, logger *logrus.Logger) error {
	config, err := ParseNewUserConfig(b, logger)
	if err != nil {
		return err
	}

	if err := config.IsValid(); err != nil {
		return err
	}

	script, err := config.addUserScript()

	return Run(b, logger, script)
}

func ParseNewUserConfig(b []byte, logger *logrus.Logger) (*NewUserConfig, error) {
	config := NewUserConfig{}
	if err := yaml.Unmarshal(b, &config); err != nil {
		return &NewUserConfig{}, err
	} else {
		logger.Debugf("new user结构体: %v", config)
		return &config, nil
	}
}

func (config *NewUserConfig) IsValid() error {
	if config.NewUser.UserDir == "" {
		config.NewUser.Name = fmt.Sprintf("/home/%s", config.NewUser.Name)
	}
	// todo 新增用户名称合法性检测
	// todo 新增用户密码合法性检测

	return nil
}

func (config NewUserConfig) addUserScript() (string, error) {
	return util.Render(addUserTmpl, util.Data{
		"NoLogin":  config.NewUser.Nologin,
		"User":     config.NewUser.Name,
		"Password": config.NewUser.Password,
		"WorkDir":  config.NewUser.UserDir,
	})
}

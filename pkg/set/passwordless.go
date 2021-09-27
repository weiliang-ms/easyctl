package set

import (
	_ "embed"
	"github.com/lithammer/dedent"
	"github.com/sirupsen/logrus"
	"github.com/weiliang-ms/easyctl/pkg/ssh"
	"github.com/weiliang-ms/easyctl/pkg/util"
	"text/template"
)

// PasswordLessTmpl  本机互信脚本模板
var PasswordLessTmpl = template.Must(template.New("PasswordLessTmpl").Parse(dedent.Dedent(`
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

func PasswordLess(config []byte, logger *logrus.Logger) error {
	script, err := MakeKeyPairScript(PasswordLessTmpl)
	if err != nil {
		return err
	}
	return Config(config, logger, script)
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

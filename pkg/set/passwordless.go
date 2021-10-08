package set

import (
	"github.com/lithammer/dedent"
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"github.com/weiliang-ms/easyctl/pkg/ssh"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"github.com/weiliang-ms/easyctl/pkg/util/tmplutil"
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

// PasswordLess 设置主机互信
func PasswordLess(item command.OperationItem) error {
	script, _ := MakeKeyPairScript(PasswordLessTmpl)
	return runner.RemoteRun(item.B, item.Logger, script)
}

// MakeKeyPairScript 生成密钥对
func MakeKeyPairScript(tmpl *template.Template) (string, error) {

	prv, pub, _ := ssh.MakeSSHKeyPair()

	return tmplutil.Render(tmpl, map[string]interface{}{
		"PublicKey":  pub,
		"PrivateKey": prv,
	})
}

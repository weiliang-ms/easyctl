package set

import (
	"github.com/lithammer/dedent"
	"text/template"
)

var setDnsShellTmpl = template.Must(template.New("").Parse(
	dedent.Dedent(`
{{- range .DnsServerList }}
sed -i "/{{ . }}/d" /etc/resolv.conf
echo nameserver {{ . }} >> /etc/resolv.conf
{{- end }}
`)))

type DNS struct {
}

func (dns DNS) Config(b []byte, debug bool) error {
	return nil
}

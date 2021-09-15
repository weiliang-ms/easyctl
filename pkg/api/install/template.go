package install

import (
	"github.com/lithammer/dedent"
	"text/template"
)

var (
	HaproxyScriptTmpl = template.Must(template.New("HaproxyScriptTmpl").Parse(
		dedent.Dedent(`
#!/bin/bash
yum install -y haproxy
`)))

	HaProxyConfigTmpl = template.Must(template.New("HaProxyConfig").Parse(
		dedent.Dedent(`
global
   log /dev/log  local0 warning
   chroot      /var/lib/haproxy
   pidfile     /var/run/haproxy.pid
   maxconn     4000
   user        haproxy
   group       haproxy
   daemon

  stats socket /var/lib/haproxy/stats

defaults
 log global
 option  httplog
 option  dontlognull
       timeout connect 5000
       timeout client 50000
       timeout server 50000

{{- if .Balance }}

{{- range .Balance }}
{{ . }}
{{- end }}

{{- end}}
   `)))
)

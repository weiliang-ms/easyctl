package tmpl

import (
	"github.com/lithammer/dedent"
	"text/template"
)

var RedisCompileTmpl = template.Must(template.New("compileTmpl").Parse(dedent.Dedent(`
#!/bin/bash
set -e
{{- if .PackageName }}
cd /tmp
if [ ! -f {{ .PackageName }} ];then
  echo /tmp/{{ .PackageName }} Not Found.
  exit 1
fi
tar zxvf {{ .PackageName }}
packageName=$(echo {{ .PackageName }}|sed 's#.tar.gz##g')
echo $packageName
cd $packageName
sed -i "s#\$(PREFIX)/bin#%s#g" src/Makefile
make -j $(nproc)
make install
{{- end}}
`)))

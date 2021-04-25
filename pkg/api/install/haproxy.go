package install

import (
	"fmt"
	"github.com/lithammer/dedent"
	"github.com/weiliang-ms/easyctl/asset"
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"github.com/weiliang-ms/easyctl/pkg/util"
	"log"
	"sync"
	"text/template"
)

var (
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

func Haproxy(i runner.Installer) {
	re := runner.ParseServerList(i.ServerListPath, runner.HaProxyServerList{})
	list := re.HA.Attribute.Server
	script, _ := asset.Asset("static/script/install_offline_tmpl.sh")
	i.Cmd = fmt.Sprintf("package_name=haproxy %s", string(script))
	i.FileName = "haproxy.tar.gz"
	if i.Offline {
		offline(i, list)
		configHA(re.HA.Attribute.BalanceList, list)
	}
}

func offline(i runner.Installer, list []runner.Server) {
	var wg sync.WaitGroup
	ch := make(chan runner.ShellResult, len(list))
	// 拷贝文件
	dstPath := fmt.Sprintf("/tmp/%s", i.FileName)

	// 生成本地临时文件
	for _, v := range list {
		runner.ScpFile(i.OfflineFilePath, dstPath, v, 0755)
		log.Println("<- transfer done ...")
	}

	// 并行
	log.Println("-> 批量安装...")
	for _, v := range list {
		wg.Add(1)
		go func(server runner.Server) {
			defer wg.Done()
			re := server.RemoteShell(i.Cmd)
			ch <- re
		}(v)
	}

	wg.Wait()
	close(ch)

	// ch -> slice
	var as []runner.ShellResult
	for target := range ch {
		as = append(as, target)
	}
}

func configHA(balance []runner.Balance, list []runner.Server) {

	var slice []string
	var ports []int
	var openPortCmd string

	for _, v := range balance {

		ports = append(ports, v.Port)

		listen := fmt.Sprintf("listen %s\n   bind *:%d\n   mode tcp\n   option tcplog\n   balance source\n", v.Name, v.Port)
		var server string
		i := 1
		for _, s := range v.Address {
			server += fmt.Sprintf("   server %s-%d %s weight 1\n", v.Name, i, s)
			i++
		}
		slice = append(slice, listen+server)
	}

	config, err := util.Render(HaProxyConfigTmpl, util.Data{
		"Balance": slice,
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, p := range ports {
		openPortCmd += util.OpenPortCmd(p)
	}

	for _, v := range list {
		log.Printf("config harproxy node -> %s", v.Host)
		v.WriteRemoteFile([]byte(config), "/etc/haproxy/haproxy.cfg", 0644)
		log.Printf("open firewalld port -> %s", v.Host)
		v.RemoteShell(openPortCmd)
		log.Printf("boot haproxy...")
		v.RemoteShell("systemctl enable haproxy --now")
		v.RemoteShell("systemctl restart haproxy")
	}

}

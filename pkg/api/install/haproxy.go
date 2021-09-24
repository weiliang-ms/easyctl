package install

//
//import (
//	_ "embed"
//	"fmt"
//	"github.com/modood/table"
//	"github.com/weiliang-ms/easyctl/pkg/runner"
//	"github.com/weiliang-ms/easyctl/pkg/util"
//	"log"
//	"os"
//)
//
//func Haproxy(configPath string) {
//	b, err := os.ReadFile(configPath)
//	if err != nil {
//		panic(err)
//	}
//	h, err := ParseConfig(b, HaproxyMeta{})
//	if err != nil {
//		panic(err)
//	}
//	h.(*HaproxyMeta).prepare().installHA().configHA()
//}
//
//// 预安装阶段：下载文件、解析server列表、检测依赖
//func (h *HaproxyMeta) prepare() *HaproxyMeta {
//	//path, err := Download(h.Haproxy.FilePath)
//	//if err != nil {
//	//	panic(err)
//	//}
//
//	var servers []runner.Server
//	for _, s := range h.Haproxy.Server {
//		for _, v := range ParseServer(h.Haproxy.Excludes, s) {
//			servers = append(servers, v)
//		}
//	}
//
//	//h.Haproxy.FilePath = path
//	h.Haproxy.Server = servers
//
//	// 分发文件
//	//if err := HandFile(h.Haproxy.Server, h.Haproxy.FilePath);err != nil {
//	//	panic(err)
//	//}
//
//	need := map[string]struct {
//		DetectShell  string
//		InstallShell string
//	}{}
//
//	// 检测依赖
//	need["gcc"] = struct {
//		DetectShell  string
//		InstallShell string
//	}{DetectShell: "gcc -v", InstallShell: "yum install -y gcc"}
//
//	for _, s := range h.Haproxy.Server {
//		d := Dependency{
//			Needs:  need,
//			Server: s,
//		}
//
//		d.DetectDep()
//	}
//
//	return h
//}
//
//// 预安装阶段：下载文件、解析server列表、检测依赖
//func (h *HaproxyMeta) installHA() *HaproxyMeta {
//	script, err := util.Render(HaproxyScriptTmpl, util.Data{
//		"FilePath": subFilename(h.Haproxy.FilePath),
//	})
//
//	if err != nil {
//		panic(err)
//	}
//
//	table.OutputA(Execute(h.Haproxy.Server, "安装haproxy", script))
//
//	return h
//}
//
//func (h *HaproxyMeta) configHA() {
//
//	var slice []string
//	var ports []int
//	var openPortCmd string
//
//	for _, v := range h.Haproxy.Balance {
//
//		ports = append(ports, v.ListenPort)
//
//		listen := fmt.Sprintf("listen %s\n   bind *:%d\n   mode tcp\n   option tcplog\n   balance source\n", v.Name, v.ListenPort)
//		var server string
//		i := 1
//		for _, s := range v.Endpoint {
//			server += fmt.Sprintf("   server %s-%d %s weight 1\n", v.Name, i, s)
//			i++
//		}
//		slice = append(slice, listen+server)
//	}
//
//	config, err := util.Render(HaProxyConfigTmpl, util.Data{
//		"Balance": slice,
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	for _, p := range ports {
//		openPortCmd += util.OpenPortCmd(p)
//	}
//
//	for _, v := range h.Haproxy.Server {
//		log.Printf("config harproxy node -> %s", v.Host)
//		v.WriteRemoteFile([]byte(config), "/etc/haproxy/haproxy.cfg", 0644)
//		log.Printf("open firewalld port -> %s", v.Host)
//		v.RemoteShell(openPortCmd)
//		log.Printf("boot haproxy...")
//		v.RemoteShell("systemctl enable haproxy --now")
//		v.RemoteShell("systemctl restart haproxy")
//	}
//
//}

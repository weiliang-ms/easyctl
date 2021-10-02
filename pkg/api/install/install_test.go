package install

//
//import (
//	_ "embed"
//	"fmt"
//	"github.com/modood/table"
//	"github.com/weiliang-ms/easyctl/pkg/runner"
//	"log"
//	"os"
//	"testing"
//)
//
////go:embed asset/haproxy.yaml
//var haproxy []byte
//
//func TestParseConfig(t *testing.T) {
//
//	log.Println("Test parse config file...")
//	h := HaproxyMeta{}
//	r, err := ParseConfig(haproxy, h)
//	if err != nil {
//		panic(err)
//	}
//
//	log.Printf("!---!%+v\n", r)
//
//	log.Println("Test parse servers...")
//
//	var s []runner.Server
//	for _, v := range r.(*HaproxyMeta).Haproxy.Server {
//		log.Printf("开始解析：%v\n", v)
//		for _, server := range ParseServer(r.(*HaproxyMeta).Haproxy.Excludes, v) {
//			fmt.Println(s)
//			s = append(s, server)
//		}
//	}
//
//	log.Printf("server列表为：%+v\n", s)
//}
//
////func TestDownload(t *testing.T) {
////	url := "http://nginx.org/download/nginx-1.20.1.tar.gz"
////	path, err := Download(url)
////	if err != nil {
////		panic(err)
////	}else {
////		log.Printf("清理文件：%s", path)
////		_ = os.RemoveAll(path)
////	}
////}
//
//func TestDependency_DetectDep(t *testing.T) {
//
//	need := map[string]struct {
//		DetectShell  string
//		InstallShell string
//	}{}
//
//	need["gcc"] = struct {
//		DetectShell  string
//		InstallShell string
//	}{DetectShell: "gcc -v", InstallShell: "yum install -y gcc"}
//
//	need["telnet"] = struct {
//		DetectShell  string
//		InstallShell string
//	}{DetectShell: "rpm -qa|grep telnet", InstallShell: "yum install -y telnet"}
//
//	d := Dependency{
//		Needs: need,
//		Server: runner.Server{
//			Host:           "10.79.160.185",
//			Port:           "22",
//			Username:       "root",
//			Password:       "Sybj@2020#",
//			PrivateKeyPath: "",
//		},
//	}
//	re := d.DetectDep()
//	log.Printf("%v\n", re)
//}
//
//func TestHandFile(t *testing.T) {
//	server := runner.Server{
//		Host:           "10.79.160.185",
//		Port:           "22",
//		Username:       "root",
//		Password:       "Sybj@2020#",
//		PrivateKeyPath: "",
//	}
//
//	testFile := "1.txt"
//	//f , _ :=os.Create(testFile)
//	os.WriteFile(testFile, []byte("ddd"), 0666)
//	//defer f.Close()
//
//	err := HandFile([]runner.Server{server}, testFile)
//	if err != nil {
//		panic(err)
//	} else {
//		os.Remove(testFile)
//	}
//
//}
//
//func TestExecute(t *testing.T) {
//
//	for i := 1; i < 41; i++ {
//		server := runner.Server{
//			Host:           fmt.Sprintf("10.10.10.%d", i),
//			Port:           "22",
//			Username:       "root",
//			Password:       "123456",
//			PrivateKeyPath: "",
//		}
//		table.OutputA(Execute([]runner.Server{server}, "调整date", "date"))
//	}
//
//}

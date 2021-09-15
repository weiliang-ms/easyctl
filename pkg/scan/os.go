package scan

import (
	"github.com/weiliang-ms/easyctl/pkg/file"
	"log"
	"net"
	"regexp"
	"strings"
)

type system struct {
	IP       string
	OS       string
	Release  string
	Kernel   string
	Platform string
	CPUCores int
	Time     string
	OpenSSH  OpenSSH
}

type OpenSSH struct {
	Version         string `yaml:"version",json:"version"`
	Port            int    `yaml:"port",json:"port"`
	PermitRootLogin bool   `yaml:"PermitRootLogin",json:"PermitRootLogin"`
}

const cpuInfo = "/proc/cpuinfo"

func OSSecurity() {
	//b, _ := asset.Asset("static/script/scan/scan_os.sh")
	//re := runner.Shell(string(b))
	//fmt.Println(re.StdOut)
}

func (sys system) kernel() system {

	filePath := "/proc/version"
	content := file.ReadAll(filePath)
	reg := regexp.MustCompile("^[0-9]\\.[0-9]*\\.[0-9]*-[0-9]*.*.x86_64")
	for _, v := range strings.Split(content, " ") {
		if reg.MatchString(v) {
			sys.Kernel = v
		}
	}

	return sys
}

func (sys system) cores() system {
	sys.CPUCores = strings.Count(file.ReadAll(cpuInfo), "processor")
	return sys
}

func (sys system) ip() system {
	addrs, err := net.InterfaceAddrs()
	sys.IP = "127.0.0.1"
	if err != nil {
		log.Fatal(err)
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				sys.IP = ipnet.IP.String()
			}
		}
	}
	return sys
}

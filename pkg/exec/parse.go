package exec

import (
	"fmt"
	"github.com/weiliang-ms/easyctl/pkg/util"
	"gopkg.in/yaml.v2"
	"k8s.io/klog"
	"log"
	"net"
	"strconv"
	"strings"
)

// Server 外部序列化
type Server struct {
	Host          string `yaml:"host"`
	Port          string `yaml:"port"`
	Username      string `yaml:"username"`
	Password      string `yaml:"password"`
	PublicKeyPath string `yaml:"publicKeyPath,omitempty"`
}

// Servers 外部解析
type Servers struct {
	Server []Server `yaml:"server"`
}

type ServerList struct {
	Server   []Server `yaml:"server"`
	Excludes []string `yaml:"excludes"`
}

// Executor 外部序列化
type Executor struct {
	Server   []Server `yaml:"server"`
	Excludes []string `yaml:"excludes"`
	Script   string   `yaml:"script"`
}

type ExecutorMeta struct {
	Executor Executor `yaml:"executor"`
}

// ExecutorItem 内部
type ExecutorItem struct {
	Server []Server
	Script string
}

// ParseConfig 解析yaml配置
func parseConfig(content []byte) (ExecutorItem, error) {
	meta := ExecutorMeta{}
	err := yaml.Unmarshal(content, &meta)

	item := ExecutorItem{}

	meta = meta.parseServers()

	item.Script = meta.Executor.Script
	item.Server = meta.Executor.Server

	if err != nil {
		return item, err
	}

	return item, nil
}

// ParseServer 解析server主机列表
func (meta ExecutorMeta) parseServers() ExecutorMeta {
	serverList := []Server{}
	for _, v := range meta.Executor.Server {
		re := v.parseServer(meta.Executor.Excludes)
		if len(re) == 1 {
			serverList = append(serverList, re[1])
		} else {
			for _, s := range re {
				serverList = append(serverList, s)
			}
		}
	}
	meta.Executor.Server = serverList
	return meta
}

// ParseServerList ParseServer 解析server主机列表
// 解析ip地址区间类型，排除excludes数组内的主机
func (sl ServerList) ParseServerList() []Server {
	servers := []Server{}
	for _, v := range sl.Server {
		re := v.parseServer(sl.Excludes)
		if len(re) == 1 {
			servers = append(servers, re[1])
		} else {
			for _, s := range re {
				servers = append(servers, s)
			}
		}
	}
	return servers
}

// 解析ip地址区间类型，排除excludes数组内的主机
func (s Server) parseServer(excludes []string) []Server {
	serverList := []Server{}

	contain := strings.Contains(s.Host, "[") && strings.Contains(s.Host, "]") && strings.Contains(s.Host, ":")

	if address := net.ParseIP(s.Host); address == nil && !contain {
		log.Fatalln("server地址信息非法，无法解析请检查...")
	}

	if strings.Contains(s.Host, "[") {
		log.Println("检测到配置文件中含有IP段，开始解析组装...")
		//192.168.235.
		baseAddress := strings.Split(s.Host, "[")[0]
		klog.Infof("解析到IP子网网段为：%s...\n", baseAddress)

		// 1:3] -> 1:3
		ipRange := strings.Split(strings.Split(s.Host, "[")[1], "]")[0]
		klog.Infof("解析到IP区间为：%s...\n", ipRange)

		// 1:3 -> 1
		begin := strings.Split(ipRange, ":")[0]
		klog.Infof("解析到起始IP为：%s...\n", fmt.Sprintf("%s%s", baseAddress, begin))

		// 1:3 -> 3
		end := strings.Split(ipRange, ":")[1]
		klog.Infof("解析到末尾IP为：%s...\n", fmt.Sprintf("%s%s", baseAddress, end))

		// string -> int
		beginIndex, _ := strconv.Atoi(begin)
		endIndex, _ := strconv.Atoi(end)

		for i := beginIndex; i <= endIndex; i++ {
			server := Server{
				Host:          fmt.Sprintf("%s%d", baseAddress, i),
				Port:          s.Port,
				Username:      s.Username,
				Password:      s.Password,
				PublicKeyPath: s.PublicKeyPath,
			}

			if !util.SliceContain(excludes, server.Host) {
				serverList = append(serverList, server)
			}
		}
	} else {
		serverList = append(serverList, s)
	}

	return serverList
}

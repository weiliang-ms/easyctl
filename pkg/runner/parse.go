package runner

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

// ServerListFilter 解析、过滤server主机列表
// 解析ip地址区间类型，排除excludes数组内的主机
func (serverListInternal ServerListInternal) ServerListFilter() []ServerInternal {
	servers := []ServerInternal{}
	for _, v := range serverListInternal.Servers {
		re := v.ServerFilter(serverListInternal.Excludes)
		if len(re) == 1 {
			servers = append(servers, re[0])
		} else {
			for _, s := range re {
				servers = append(servers, s)
			}
		}
	}
	return servers
}

// ParseServerList ServerList反序列化
func ParseServerList(b []byte) ([]ServerInternal, error) {
	serverList := ServerListExternal{}
	if err := yaml.Unmarshal(b, &serverList); err != nil {
		return []ServerInternal{}, err
	} else {
		serverListInternal := serverListDeepCopy(serverList)
		return serverListInternal.ServerListFilter(), nil
	}
}

// ParseExecutor 执行器反序列化
func ParseExecutor(b []byte) (ExecutorInternal, error) {
	executor := ExecutorExternal{}
	err := yaml.Unmarshal(b, &executor)
	if err != nil {
		return ExecutorInternal{}, err
	}
	// 类型转换
	executorInternal := executorDeepCopy(executor)

	// 解析ip地址段
	servers := []ServerInternal{}
	for _, v := range executorInternal.Servers {
		for _, s := range v.ServerFilter(executor.Excludes) {
			servers = append(servers, s)
		}
	}
	executorInternal.Servers = servers
	return executorInternal, nil
}

// todo:深拷贝
func serverListDeepCopy(serverListExternal ServerListExternal) ServerListInternal {
	return ServerListInternal{
		serversDeepCopy(serverListExternal.Servers),
		serverListExternal.Excludes,
	}
}

// todo:深拷贝
func executorDeepCopy(executorExternal ExecutorExternal) ExecutorInternal {
	return ExecutorInternal{
		serversDeepCopy(executorExternal.Servers),
		executorExternal.Script,
	}
}

// todo:深拷贝
func serversDeepCopy(external []ServerExternal) []ServerInternal {
	internal := []ServerInternal{}
	for _, v := range external {
		internal = append(internal, serverDeepCopy(v))
	}
	return internal
}

// todo:深拷贝
func serverDeepCopy(serverExternal ServerExternal) ServerInternal {
	return ServerInternal{
		serverExternal.Host,
		serverExternal.Port,
		serverExternal.Username,
		serverExternal.Password,
		serverExternal.PrivateKeyPath,
	}
}

// ServerFilter 解析ip地址区间类型，排除excludes数组内的主机
func (server ServerInternal) ServerFilter(excludes []string) []ServerInternal {
	serverList := []ServerInternal{}

	contain := strings.Contains(server.Host, "[") && strings.Contains(server.Host, "]") && strings.Contains(server.Host, ":")

	if address := net.ParseIP(server.Host); address == nil && !contain {
		log.Fatalln("server地址信息非法，无法解析请检查...")
	}

	if strings.Contains(server.Host, "[") {
		log.Println("检测到配置文件中含有IP段，开始解析组装...")
		//192.168.235.
		baseAddress := strings.Split(server.Host, "[")[0]
		klog.Infof("解析到IP子网网段为：%s...\n", baseAddress)

		// 1:3] -> 1:3
		ipRange := strings.Split(strings.Split(server.Host, "[")[1], "]")[0]
		klog.Infof("解析到IP区间为：%s...\n", ipRange)

		// 1:3 -> 1
		begin := strings.Split(ipRange, ":")[0]
		klog.Infof("解析到起始IP为：%s...\n", fmt.Sprintf("%s%s", baseAddress, begin))

		// 1:3 -> 3
		end := strings.Split(ipRange, ":")[1]
		klog.Infof("解析到末尾IP为：%s...\n", fmt.Sprintf("%s%s", baseAddress, end))

		// 区间首尾一直直接返回
		if begin == end {
			return append(serverList, ServerInternal{
				Host:           begin,
				Port:           server.Port,
				Username:       server.Username,
				Password:       server.Password,
				PrivateKeyPath: server.PrivateKeyPath,
			})
		}

		// string -> int
		beginIndex, _ := strconv.Atoi(begin)
		endIndex, _ := strconv.Atoi(end)

		for i := beginIndex; i <= endIndex; i++ {
			server := ServerInternal{
				Host:           fmt.Sprintf("%s%d", baseAddress, i),
				Port:           server.Port,
				Username:       server.Username,
				Password:       server.Password,
				PrivateKeyPath: server.PrivateKeyPath,
			}

			if !util.SliceContain(excludes, server.Host) {
				serverList = append(serverList, server)
			}
		}
	} else {
		serverList = append(serverList, server)
	}

	return serverList
}

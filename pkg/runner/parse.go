package runner

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/weiliang-ms/easyctl/pkg/util"
	"gopkg.in/yaml.v2"
	"net"
	"strconv"
	"strings"
)

type serverFilter struct {
	Servers []ServerInternal
}

// ServerListFilter 解析、过滤server主机列表
// 解析ip地址区间类型，排除excludes数组内的主机
func (serverListInternal ServerListInternal) serverListFilter(logger *logrus.Logger) ([]ServerInternal, error) {
	var servers []ServerInternal
	serverMap := make(map[string]ServerInternal)

	filter := &serverFilter{}
	for _, v := range serverListInternal.Servers {
		if err := v.parseIPRangeServer(filter, logger); err != nil {
			return nil, err
		}
	}

	for _, v := range filter.Servers {
		if !contain(v, serverListInternal.Excludes) {
			serverMap[v.Host] = v
		}
	}

	for _, v := range serverMap {
		servers = append(servers, v)
	}

	return servers, nil
}

// server是否存在于excludes servers列表内
func contain(server ServerInternal, excludeServers []string) bool {
	return util.SliceContain(excludeServers, server.Host)
}

// ParseServerList ServerList反序列化
func ParseServerList(b []byte, logger *logrus.Logger) ([]ServerInternal, error) {
	serverList := ServerListExternal{}
	if err := yaml.Unmarshal(b, &serverList); err != nil {
		return []ServerInternal{}, err
	}
	serverListInternal := serverListDeepCopy(serverList)
	return serverListInternal.serverListFilter(logger)
}

// ParseExecutor 执行器反序列化
func ParseExecutor(b []byte, logger *logrus.Logger) (ExecutorInternal, error) {

	if logger == nil {
		logger = logrus.New()
	}

	executor := ExecutorExternal{}
	err := yaml.Unmarshal(b, &executor)
	if err != nil {
		return ExecutorInternal{}, err
	}
	// 类型转换
	executorInternal := executorDeepCopy(executor)

	// 解析地址段
	filter := &serverFilter{}
	for _, v := range executorInternal.Servers {
		if err := v.parseIPRangeServer(filter, logger); err != nil {
			return ExecutorInternal{}, err
		}
	}

	executorInternal.Servers = filter.Servers
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
		logrus.New(),
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
func (server ServerInternal) parseIPRangeServer(filter *serverFilter, logger *logrus.Logger) error {

	// TODO 正则
	contain := strings.Contains(server.Host, "[") && strings.Contains(server.Host, "]") && strings.Contains(server.Host, ":")

	if address := net.ParseIP(server.Host); address == nil && !contain {
		return fmt.Errorf("server地址: %s 非法，无法解析请检查", server.Host)
	}

	if strings.Contains(server.Host, "[") {
		logger.Println("检测到配置文件中含有IP段，开始解析组装...")
		//192.168.235.
		baseAddress := strings.Split(server.Host, "[")[0]
		logger.Infof("解析到IP子网网段为：%s...", baseAddress)

		// 1:3] -> 1:3
		ipRange := strings.Split(strings.Split(server.Host, "[")[1], "]")[0]
		logger.Infof("解析到IP区间为：%s...", ipRange)

		// 1:3 -> 1
		begin := strings.Split(ipRange, ":")[0]
		logger.Infof("解析到起始IP为：%s...", fmt.Sprintf("%s%s", baseAddress, begin))

		// 1:3 -> 3
		end := strings.Split(ipRange, ":")[1]
		logger.Infof("解析到末尾IP为：%s...", fmt.Sprintf("%s%s", baseAddress, end))

		// 区间首尾一致直接返回
		if begin == end {
			filter.Servers = append(filter.Servers, ServerInternal{
				Host:           fmt.Sprintf("%s%s", baseAddress, begin),
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
			filter.Servers = append(filter.Servers, server)
		}
	} else {
		filter.Servers = append(filter.Servers, server)
	}

	return nil
}

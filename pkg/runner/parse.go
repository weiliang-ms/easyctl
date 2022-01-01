/*
	MIT License

Copyright (c) 2020 xzx.weiliang

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

*/
package runner

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/weiliang-ms/easyctl/pkg/util/slice"
	strings2 "github.com/weiliang-ms/easyctl/pkg/util/strings"
	"gopkg.in/yaml.v2"
	"io"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type serverFilter struct {
	Servers []ServerInternal
}

// ip地址range
type addressInterval struct {
	BeginIndex int
	EndIndex   int
	Cidr       string
}

var validSplitChar = []string{"..", "-", ":"}

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
		} else {
			logger.Infof("排除ip: %s", v.Host)
		}
	}

	for _, v := range serverMap {
		servers = append(servers, v)
	}

	return servers, nil
}

// server是否存在于excludes servers列表内
func contain(server ServerInternal, excludeServers []string) bool {
	return slice.StringSliceContain(excludeServers, server.Host)
}

// ParseServerList ServerList反序列化
func ParseServerList(b []byte, logger *logrus.Logger) ([]ServerInternal, error) {
	if logger == nil {
		logger = logrus.New()
	}

	serverList := ServerListExternal{}
	if err := yaml.Unmarshal(b, &serverList); err != nil {
		return []ServerInternal{}, err
	}
	serverListInternal := serverListDeepCopy(serverList)

	defer logger.Info("解析server列表完毕!")
	return serverListInternal.serverListFilter(logger)
}

// ParseExecutor 执行器反序列化
func ParseExecutor(b []byte, logger *logrus.Logger) (ExecutorInternal, error) {

	if logger == nil {
		logger = logrus.New()
	}

	executor := ExecutorExternal{}

	logger.Debugf("执行器反序列化内容为: %s", string(b))
	err := yaml.Unmarshal(b, &executor)
	if err != nil {
		return ExecutorInternal{}, err
	}

	logger.Debugf("执行器配置为: %#v", executor)

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

	// 格式化输出结构体
	bs, _ := json.Marshal(executorInternal.Servers)
	var out bytes.Buffer
	_ = json.Indent(&out, bs, "", "\t")
	logger.Debugf("主机列表：\n%v", out.String())

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
		false,
	}
}

// todo:深拷贝
func serversDeepCopy(external []ServerExternal) []ServerInternal {
	var internal []ServerInternal
	for _, v := range external {
		internal = append(internal, serverDeepCopy(v))
	}
	return internal
}

// todo:深拷贝
func serverDeepCopy(serverExternal ServerExternal) ServerInternal {
	return ServerInternal{
		serverExternal.Host,
		serverExternal.Hosts,
		serverExternal.Port,
		serverExternal.Username,
		serverExternal.Password,
		serverExternal.PrivateKeyPath,
	}
}

// ServerFilter 解析ip地址区间类型，排除excludes数组内的主机
func (server ServerInternal) parseIPRangeServer(filter *serverFilter, logger *logrus.Logger) error {

	logger.Info("分析host是否为地址段/地址列表")
	
	if server.Hosts != nil {
		err := server.parseHosts(filter, logger)
		if err != nil {
			return err
		}
	}
	if len(server.Host) > 0 {
		err := server.parseHost(filter, logger)
		if err != nil {
			return err
		}
	}
	return nil
}

func (server ServerInternal) parseHosts(filter *serverFilter, logger *logrus.Logger) error {
	for _, host := range server.Hosts {
		err := ServerInternal{
			Host:           host,
			Port:           server.Port,
			Username:       server.Username,
			Password:       server.Password,
			PrivateKeyPath: server.PrivateKeyPath,
		}.parseHost(filter, logger)
		if err != nil {
			return err
		}
	}
	server.Hosts = nil
	return nil
}

func (server ServerInternal) parseHost(filter *serverFilter, logger *logrus.Logger) error {
	// host入参为文件类型(主机列表)
	if _, err := os.Stat(server.Host); err == nil {

		logger.Info("解析到Host类型为文件列表")
		file, _ := os.OpenFile(server.Host, os.O_RDWR, 0666)
		defer file.Close()

		buf := bufio.NewReader(file)
		for {
			line, err := buf.ReadString('\n')
			line = strings.TrimSpace(line)

			if ok := net.ParseIP(line); ok != nil {
				filter.Servers = append(filter.Servers, ServerInternal{
					Host:           line,
					Port:           server.Port,
					Username:       server.Username,
					Password:       server.Password,
					PrivateKeyPath: server.PrivateKeyPath,
				})
			}

			if err != nil {
				if err == io.EOF {
					logger.Debug("File read ok!")
					return nil
				}
			}
		}
	}

	if ok := net.ParseIP(server.Host); ok != nil {
		filter.Servers = append(filter.Servers, server)
		return nil
	}

	logger.Infoln("检测到配置文件中可能含有IP地址区间，开始解析组装...")
	flysnowRegexp := regexp.MustCompile("^((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})(\\.((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})){2}\\.")
	cidr := flysnowRegexp.FindString(server.Host)
	if cidr == "" {
		return fmt.Errorf("%s 地址区间非法", server.Host)
	}

	logger.Infof("截取到IP地址区间: %s0/24", cidr)

	rangeStr := strings.TrimPrefix(server.Host, cidr)
	logger.Infof("区间为: %s", rangeStr)

	logger.Infoln("开始组装地址区间类型server")
	for _, v := range packageIPRange(server, getInterval(cidr, rangeStr, logger)) {
		filter.Servers = append(filter.Servers, v)
	}

	logger.Infoln("地址区间类型server组装完毕！")
	return nil
}

func packageIPRange(server ServerInternal, interval addressInterval) []ServerInternal {
	var servers []ServerInternal
	if interval.BeginIndex == 0 && interval.EndIndex == 0 {
		return servers
	}

	if interval.BeginIndex > 255 || interval.EndIndex > 255 {
		return servers
	}

	for i := interval.BeginIndex; i <= interval.EndIndex; i++ {
		s := ServerInternal{
			Host:           fmt.Sprintf("%s%d", interval.Cidr, i),
			Port:           server.Port,
			Username:       server.Username,
			Password:       server.Password,
			PrivateKeyPath: server.PrivateKeyPath,
		}
		servers = append(servers, s)
	}
	return servers
}

func getInterval(cidr, rangeStr string, logger *logrus.Logger) addressInterval {
	interval := addressInterval{Cidr: cidr}

	// trim [1?2] -> 1?2
	if strings.Contains(rangeStr, "[") {
		rangeStr = strings.TrimPrefix(rangeStr, "[")
	}
	if strings.Contains(rangeStr, "]") {
		rangeStr = strings.TrimSuffix(rangeStr, "]")
	}

	result, err := strings2.SplitIfContain(rangeStr, validSplitChar)
	if err != nil {
		return interval
	}

	if len(result) != 2 {
		logger.Info(result)
		return interval
	}

	logger.Infof("解析到起始IP为：%s...", fmt.Sprintf("%s%s", cidr, result[0]))
	logger.Infof("解析到末尾IP为：%s...", fmt.Sprintf("%s%s", cidr, result[1]))

	interval.BeginIndex, _ = strconv.Atoi(result[0])
	interval.EndIndex, _ = strconv.Atoi(result[1])
	interval.Cidr = cidr

	return interval
}

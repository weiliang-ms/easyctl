package scan

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"github.com/weiliang-ms/easyctl/pkg/util/format"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type OSInfo struct {
	BaseOSInfo
	CPUInfo
	MemoryInfo
}

type OSInfoSlice []OSInfo

type BaseOSInfo struct {
	Address  string // ip地址
	Hostname string // 主机名
	KernelV  string // 内核版本
	OSV      string // 操作系统版本
}

type CPUInfo struct {
	CPUThreadCount string // cpu线程数
	CPUClockSpeed  string // cpu主频
	CPUModeNum     string // CPU版本号
	CPULoadAverage string // CPU版本号
}

type MemoryInfo struct {
	MemUsed       float64 // 内存使用量
	MemUsePercent float64 // 内存使用占比
	MemTotal      float64 // 内存总量
}

// OS 扫描系统信息
func OS(item command.OperationItem) command.RunErr {
	servers, err := runner.ParseServerList(item.B, item.Logger)
	if err != nil {
		return command.RunErr{Err: err, Msg: "解析异常"}
	}

	var result OSInfoSlice

	ch := make(chan OSInfo, len(servers))
	wg := sync.WaitGroup{}

	for _, v := range servers {
		go func(s runner.ServerInternal) {
			wg.Add(1)
			ch <- osInfo(s, item.Logger)
			defer wg.Done()
		}(v)
	}

	go func() {
		for v := range ch {
			result = append(result, v)
		}
	}()

	wg.Wait()

	// 排序
	sort.Sort(result)
	out, err := format.Object(result)
	if err != nil {
		panic(err)
	}
	item.Logger.Infof("系统信息：\n%v", out.String())

	return command.RunErr{}
}

func osInfo(s runner.ServerInternal, logger *logrus.Logger) OSInfo {

	var osInfo OSInfo
	var cpuInfo CPUInfo
	baseInfo := BaseOSInfo{Address: s.Host}

	if re := s.ReturnRunResult(runner.RunItem{
		Logger: logger,
		Cmd:    "hostname",
	}); re.Err != nil {
		panic(re.Err)
	} else {
		baseInfo.Hostname = strings.TrimSuffix(re.StdOut, "\n")
	}

	if re := s.ReturnRunResult(runner.RunItem{
		Logger: logger,
		Cmd:    "uname -r",
	}); re.Err != nil {
		panic(re.Err)
	} else {
		baseInfo.KernelV = strings.TrimSuffix(re.StdOut, "\n")
	}

	if re := s.ReturnRunResult(runner.RunItem{
		Logger: logger,
		Cmd:    "cat /etc/system-release",
	}); re.Err != nil {
		panic(re.Err)
	} else {
		baseInfo.OSV = strings.TrimSuffix(re.StdOut, "\n")
	}

	if re := s.ReturnRunResult(runner.RunItem{
		Logger: logger,
		Cmd:    "cat /proc/cpuinfo",
	}); re.Err != nil {
		panic(re.Err)
	} else {
		cpuInfo = NewCPUInfoItem(re.StdOut)
	}

	if re := s.ReturnRunResult(runner.RunItem{
		Logger: logger,
		Cmd:    "cat /proc/loadavg|awk '{print $1,$2,$3}'",
	}); re.Err != nil {
		panic(re.Err)
	} else {
		cpuInfo.CPULoadAverage = strings.TrimSpace(re.StdOut)
	}

	osInfo.BaseOSInfo = baseInfo
	osInfo.CPUInfo = cpuInfo
	return osInfo
}

func (re OSInfoSlice) Len() int { return len(re) }

func (re OSInfoSlice) Swap(i, j int) {
	re[i], re[j] = re[j], re[i]
}

func (re OSInfoSlice) Less(i, j int) bool {
	address1 := strings.Split(re[i].Address, ".")
	address2 := strings.Split(re[j].Address, ".")

	for k := 0; k < 4; k++ {
		if address1[k] != address2[k] {
			num1, _ := strconv.Atoi(address1[k])
			num2, _ := strconv.Atoi(address2[k])
			return num1 < num2
		}
	}

	return true
}

func NewCPUInfoItem(content string) CPUInfo {
	var c CPUInfo
	var count int
	for _, v := range strings.Split(content, "\n") {

		if strings.Contains(v, "processor") {
			count++
		}

		reg := regexp.MustCompile("^model name")
		if reg.MatchString(v) && c.CPUModeNum == "" && c.CPUClockSpeed == "" {
			c.CPUModeNum = strings.TrimSpace(strings.Split(strings.Split(v, ":")[1], "@")[0])
			c.CPUClockSpeed = strings.TrimSpace(strings.Split(strings.Split(v, ":")[1], "@")[1])
		}
	}

	c.CPUThreadCount = fmt.Sprintf("%d", count)
	return c
}

func NewMemInfoItem(content string) MemoryInfo {
	var c MemoryInfo
	var total, available float64

	for _, v := range strings.Split(content, "\n") {

		if strings.Contains(v, "MemTotal") {
			s := strings.Split(v, " ")[1]
			total, _ = strconv.ParseFloat(s, 64)
		}

		if strings.Contains(v, "MemAvailable") {
			s := strings.Split(v, " ")[1]
			available, _ = strconv.ParseFloat(s, 64)
		}

	}

	used := total - available
	c.MemTotal = total / (1024 * 1024)
	c.MemUsed = used / (1024 * 1024)
	c.MemUsePercent = used / total * 100
	return c
}

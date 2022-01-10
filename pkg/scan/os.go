package scan

import (
	"bytes"
	strings2 "github.com/weiliang-ms/easyctl/pkg/util/strings"
	"io"
	"os"
	//
	_ "embed"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"github.com/weiliang-ms/easyctl/pkg/util/format"
	slice2 "github.com/weiliang-ms/easyctl/pkg/util/slice"
	"github.com/xuri/excelize/v2"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type OSInfo struct {
	BaseOSInfo
	CPUInfo
	DiskInfo
	MemoryInfo
}

const (
	MountHighUsedValue      = 90
	PrintHostnameShell      = "hostname"
	PrintKernelVersionShell = "uname -r"
	PrintOSVersionShell     = "cat /etc/system-release"
	PrintCPUInfoShell       = "cat /proc/cpuinfo"
	PrintMemInfoShell       = "cat /proc/meminfo"
	PrintCPULoadavgShell    = "cat /proc/loadavg|awk '{print $1,$2,$3}'"
	PrintMountInfoShell     = "df -h|grep -v Filesystem"
)

var UnitTest bool

type OSInfoSlice []OSInfo

type BaseOSInfo struct {
	Address  string // ip地址
	Hostname string // 主机名
	KernelV  string // 内核版本
	OSV      string // 操作系统版本
}

type CPUInfo struct {
	CPUThreadCount int    // cpu线程数
	CPUClockSpeed  string // cpu主频
	CPUModeNum     string // CPU版本号
	CPULoadAverage string // CPU版本号
}

type MemoryInfo struct {
	MemUsed       float64 // 内存使用量
	MemUsePercent float64 // 内存使用占比
	MemTotal      float64 // 内存总量
}

type DiskInfoMeta struct {
	Filesystem   string
	Size         string
	Used         string
	Avail        string
	UsedPercent  string
	MountedPoint string
}

type DiskInfo struct {
	RootMountUsedPercent      int    // 根路径使用率
	HighUsedPercentMountPoint string // 90%使用率挂载点
	OSDriveLetter             string // 系统分区
}

// OS 扫描系统信息
func OS(item command.OperationItem) command.RunErr {
	servers, err := runner.ParseServerList(item.B, item.Logger)

	if err != nil {
		return command.RunErr{Err: err, Msg: "解析异常"}
	}

	serversOut := format.ObjectToJson(servers)
	item.Logger.Debugf("列表信息：%s", &serversOut)

	var result OSInfoSlice

	ch := make(chan OSInfo, len(servers))
	errCh := make(chan error, len(servers))
	wg := sync.WaitGroup{}
	wg.Add(len(servers))

	for _, v := range servers {
		go func(s runner.ServerInternal) {
			re, scanErr := osInfo(s, item.Logger)
			if scanErr != nil {
				errCh <- fmt.Errorf("[%s] 扫描异常: %s", s.Host, scanErr)
			}
			ch <- re
			defer wg.Done()
		}(v)
	}

	wg.Wait()

	close(ch)
	close(errCh)

	for v := range ch {
		result = append(result, v)
	}

	// 排序
	sort.Sort(result)
	out := format.ObjectToJson(result)

	item.Logger.Infof("系统信息：\n%v", out.String())

	for v := range errCh {
		item.Logger.Error(v)
	}

	return command.RunErr{Err: SaveAsExcel(result)}
}

// 获取操作系统信息
func osInfo(s runner.ServerInternal, logger *logrus.Logger) (osInfo OSInfo, err error) {

	defer HandleTestErr(UnitTest, &err)

	var cpuInfo CPUInfo
	var memInfo MemoryInfo
	var diskInfo DiskInfo

	baseInfo := BaseOSInfo{Address: s.Host}

	if re := s.ReturnRunResult(runner.RunItem{
		Logger: logger,
		Cmd:    PrintHostnameShell,
	}); re.Err != nil && !UnitTest {
		return osInfo, err
	} else {
		hostname := strings.TrimSuffix(re.StdOut, "\n")
		logger.Debugf("[%s] 主机名为：%s", s.Host, hostname)
		baseInfo.Hostname = hostname
	}

	if re := s.ReturnRunResult(runner.RunItem{
		Logger: logger,
		Cmd:    PrintKernelVersionShell,
	}); re.Err != nil && !UnitTest {
		return osInfo, err
	} else {
		baseInfo.KernelV = strings.TrimSuffix(re.StdOut, "\n")
	}

	if re := s.ReturnRunResult(runner.RunItem{
		Logger: logger,
		Cmd:    PrintOSVersionShell,
	}); re.Err != nil && !UnitTest {
		return osInfo, err
	} else {
		baseInfo.OSV = strings.TrimSuffix(re.StdOut, "\n")
	}

	if re := s.ReturnRunResult(runner.RunItem{
		Logger: logger,
		Cmd:    PrintCPUInfoShell,
	}); re.Err != nil && !UnitTest {
		return osInfo, err
	} else {
		cpuInfo = NewCPUInfoItem(re.StdOut)
	}

	if re := s.ReturnRunResult(runner.RunItem{
		Logger: logger,
		Cmd:    PrintCPULoadavgShell,
	}); re.Err != nil && !UnitTest {
		return osInfo, err
	} else {
		cpuInfo.CPULoadAverage = strings.TrimSpace(re.StdOut)
	}

	if re := s.ReturnRunResult(runner.RunItem{
		Logger: logger,
		Cmd:    PrintMemInfoShell,
	}); re.Err != nil && !UnitTest {
		return osInfo, err
	} else {
		memInfo = NewMemInfoItem(strings.TrimSpace(re.StdOut))
	}

	if re := s.ReturnRunResult(runner.RunItem{
		Logger: logger,
		Cmd:    PrintMountInfoShell,
	}); re.Err != nil && !UnitTest {
		return osInfo, err
	} else {
		diskInfo = NewDiskInfoItem(strings.TrimSpace(re.StdOut))
	}

	return OSInfo{
		BaseOSInfo: baseInfo,
		CPUInfo:    cpuInfo,
		DiskInfo:   diskInfo,
		MemoryInfo: memInfo,
	}, nil
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
			if re := strings.Split(strings.Split(v, ":")[1], "@"); len(re) > 1 {
				c.CPUClockSpeed = strings.TrimSpace(re[1])
			}
		}
	}

	c.CPUThreadCount = count
	return c
}

func NewMemInfoItem(content string) MemoryInfo {
	var c MemoryInfo
	var total, available float64

	for _, v := range strings.Split(content, "\n") {

		if strings.Contains(v, "MemTotal") {
			slice := strings.Split(v, " ")
			total, _ = strconv.ParseFloat(slice[len(slice)-2], 64)
		}

		if strings.Contains(v, "MemAvailable") {
			slice := strings.Split(v, " ")
			available, _ = strconv.ParseFloat(slice[len(slice)-2], 64)
		}

	}

	used := total - available

	floatTotal, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", total/(1024*1024)), 64)
	floatUsed, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", used/(1024*1024)), 64)
	percent, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", (floatUsed/floatTotal)*100), 64)
	c.MemUsePercent = percent
	c.MemTotal = floatTotal
	c.MemUsed = floatUsed

	return c
}

func NewDiskInfoItem(content string) DiskInfo {
	var metas []DiskInfoMeta
	var disk DiskInfo

	for _, v := range strings.Split(content, "\n") {
		if v != "" {
			metas = append(metas, NewDiskInfoMetaItem(v))
		}
	}

	for _, v := range metas {

		percent, _ := strconv.Atoi(strings.TrimSuffix(v.UsedPercent, "%"))

		if v.MountedPoint == "/" {
			disk.RootMountUsedPercent = percent
			disk.OSDriveLetter = strings2.TrimNumSuffix(v.Filesystem)
		}

		if percent > MountHighUsedValue {
			disk.HighUsedPercentMountPoint += fmt.Sprintf("%s,", v.MountedPoint)
		}
	}

	disk.HighUsedPercentMountPoint = strings.TrimSuffix(disk.HighUsedPercentMountPoint, ",")

	return disk
}

/*
	数据结构
	/dev/vda1        40G  3.4G   34G   9% /
*/
func NewDiskInfoMetaItem(content string) (diskInfoMeta DiskInfoMeta) {

	slice := slice2.StringSliceFilter(strings.Split(content, " "), "")

	if len(slice) != 6 {
		return
	}

	diskInfoMeta.Filesystem = slice[0]
	diskInfoMeta.Size = slice[1]
	diskInfoMeta.Used = slice[2]
	diskInfoMeta.Avail = slice[3]
	diskInfoMeta.UsedPercent = slice[4]
	diskInfoMeta.MountedPoint = slice[5]

	return
}

//go:embed system-tmpl.xlsx
var excelTmpl []byte

func SaveAsExcel(data []OSInfo) error {

	sheet := "Sheet1"

	excel, err := os.Create("system.xlsx")
	defer excel.Close()

	if err != nil {
		return err
	}

	if _, err = io.Copy(excel, bytes.NewReader(excelTmpl)); err != nil {
		return err
	}

	f, err := excelize.OpenFile("system.xlsx")
	if err != nil {
		return err
	}

	row := 3

	maps := make(map[string]interface{})
	for _, v := range data {
		maps[fmt.Sprintf("A%d", row)] = v.Address
		maps[fmt.Sprintf("B%d", row)] = v.Hostname
		maps[fmt.Sprintf("C%d", row)] = v.OSV
		maps[fmt.Sprintf("D%d", row)] = v.KernelV
		maps[fmt.Sprintf("E%d", row)] = v.CPUThreadCount
		maps[fmt.Sprintf("F%d", row)] = v.CPUClockSpeed
		maps[fmt.Sprintf("G%d", row)] = v.CPUModeNum
		maps[fmt.Sprintf("H%d", row)] = v.CPULoadAverage
		maps[fmt.Sprintf("I%d", row)] = v.MemTotal
		maps[fmt.Sprintf("J%d", row)] = v.MemUsed
		maps[fmt.Sprintf("K%d", row)] = v.MemUsePercent
		maps[fmt.Sprintf("L%d", row)] = v.RootMountUsedPercent
		maps[fmt.Sprintf("M%d", row)] = v.OSDriveLetter
		maps[fmt.Sprintf("N%d", row)] = v.HighUsedPercentMountPoint
		row++
	}

	for k, v := range maps {
		f.SetSheetRow(sheet, k, &[]interface{}{v})
	}

	return f.Save()
}

func HandleTestErr(unitTest bool, err *error) {

	if v := recover(); v != nil {
		if unitTest {
			*err = nil
		} else {
			panic(v)
		}
	}
}

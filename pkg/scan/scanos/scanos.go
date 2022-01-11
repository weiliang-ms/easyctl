package scanos

import (
	"fmt"
	slice2 "github.com/weiliang-ms/easyctl/pkg/util/slice"
	strings2 "github.com/weiliang-ms/easyctl/pkg/util/strings"
	"regexp"
	"strconv"
	"strings"
)

const MountHighUsedValue = 90

type MetaInfo struct {
	BaseOSInfo
	CPUInfo
	DiskInfo
	MemoryInfo
}

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

type MetaInfoSlice []MetaInfo

func (re MetaInfoSlice) Len() int { return len(re) }

func (re MetaInfoSlice) Swap(i, j int) {
	re[i], re[j] = re[j], re[i]
}

func (re MetaInfoSlice) Less(i, j int) bool {
	address1 := strings.Split(re[i].BaseOSInfo.Address, ".")
	address2 := strings.Split(re[j].BaseOSInfo.Address, ".")

	result := true

	for k := 0; k < 4; k++ {
		if len(address1) == 4 && len(address2) == 4 && address1[k] != address2[k] {
			num1, _ := strconv.Atoi(address1[k])
			num2, _ := strconv.Atoi(address2[k])
			result = num1 < num2
			break
		}
	}

	return result
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

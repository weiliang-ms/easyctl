package scan

import (
	"bytes"
	"github.com/weiliang-ms/easyctl/pkg/scan/scanos"
	"io"
	"os"
	//
	_ "embed"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"github.com/weiliang-ms/easyctl/pkg/util/format"
	"github.com/xuri/excelize/v2"
	"sort"
	"sync"
)

const (
	PrintHostnameShell      = "hostname"
	PrintKernelVersionShell = "uname -r"
	PrintOSVersionShell     = "cat /etc/system-release"
	PrintCPUInfoShell       = "cat /proc/cpuinfo"
	PrintMemInfoShell       = "cat /proc/meminfo"
	PrintCPULoadavgShell    = "cat /proc/loadavg|awk '{print $1,$2,$3}'"
	PrintMountInfoShell     = "df -h|grep -v Filesystem"
	preserveFileName        = "system.xlsx"
)

type osScanner struct {
	Logger           *logrus.Logger
	HandlerInterface HandleOSInterface
}

//go:embed scanos/system-tmpl.xlsx
var excelTmpl []byte

// OS 扫描系统信息
func OS(item command.OperationItem) command.RunErr {
	servers, err := runner.ParseServerList(item.B, item.Logger)

	if err != nil {
		return command.RunErr{Err: err, Msg: "解析异常"}
	}

	serversOut := format.ObjectToJson(servers)
	item.Logger.Debugf("列表信息：%s", &serversOut)

	scanner := osScanner{
		Logger: item.Logger,
	}

	scanner.HandlerInterface = getHandlerInterface(item.Interface)

	result, err := scanner.ParallelGetOSInfo(servers)
	if err != nil {
		return command.RunErr{Err: err}
	}

	return command.RunErr{Err: SaveAsExcel(result, preserveFileName, excelTmpl)}
}
func (scanner osScanner) ParallelGetOSInfo(servers []runner.ServerInternal) (result scanos.MetaInfoSlice, err error) {
	wg := sync.WaitGroup{}
	wg.Add(len(servers))

	ch := make(chan scanos.MetaInfo, len(servers))
	errCh := make(chan error, len(servers))

	f := func(s runner.ServerInternal, osScanner osScanner) (osInfo scanos.MetaInfo, err error) {

		// 1. get hostname and ip address
		osInfo.Address = s.Host
		hostname, err := osScanner.HandlerInterface.GetHostName(s, osScanner.Logger)
		if err != nil {
			return osInfo, err
		}
		osInfo.Hostname = hostname

		// 2. get kernel version
		kernelV, err := osScanner.HandlerInterface.GetKernelVersion(s, osScanner.Logger)
		if err != nil {
			return osInfo, err
		}
		osInfo.KernelV = kernelV

		// 3. get system version
		systemV, err := osScanner.HandlerInterface.GetSystemVersion(s, osScanner.Logger)
		if err != nil {
			return osInfo, err
		}
		osInfo.OSV = systemV

		// 4. get cpu info
		cpuContent, err := osScanner.HandlerInterface.GetCPUInfo(s, osScanner.Logger)
		if err != nil {
			return osInfo, err
		}
		osInfo.CPUInfo = scanos.NewCPUInfoItem(cpuContent)

		// 5. get cpu load average
		cpuLoadAverage, err := osScanner.HandlerInterface.GetCPULoadAverage(s, osScanner.Logger)
		if err != nil {
			return osInfo, err
		}
		osInfo.CPULoadAverage = cpuLoadAverage

		// 6. get Mem info
		memContent, err := osScanner.HandlerInterface.GetMemoryInfo(s, osScanner.Logger)
		if err != nil {
			return osInfo, err
		}
		osInfo.MemoryInfo = scanos.NewMemInfoItem(memContent)

		// 7. get mount info
		diskContent, err := osScanner.HandlerInterface.GetMountPointInfo(s, osScanner.Logger)
		if err != nil {
			return osInfo, err
		}
		osInfo.DiskInfo = scanos.NewDiskInfoItem(diskContent)

		return osInfo, nil
	}

	for _, v := range servers {
		go func(s runner.ServerInternal) {
			re, scanErr := f(s, scanner)
			if scanErr != nil {
				fmt.Println(scanErr)
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

	scanner.Logger.Infof("系统信息：\n%v", out.String())

	for v := range errCh {
		return nil, v
	}

	return result, nil

}

func SaveAsExcel(data []scanos.MetaInfo, path string, tmplBytes []byte) error {

	sheet := "Sheet1"

	excel, err := os.Create(path)
	defer excel.Close()

	if err != nil {
		return err
	}

	// todo: err case
	//if _, err = io.Copy(excel, bytes.NewReader(tmplBytes)); err != nil {
	//	return err
	//}
	_, _ = io.Copy(excel, bytes.NewReader(tmplBytes))

	// todo: err case
	//f, err := excelize.OpenFile(preserveFileName)
	//if err != nil {
	//	return err
	//}
	f, _ := excelize.OpenFile(preserveFileName)

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

func getHandlerInterface(i interface{}) HandleOSInterface {
	handlerInterface, _ := i.(HandleOSInterface)
	if handlerInterface == nil {
		return new(OsExecutor)
	}
	return handlerInterface
}

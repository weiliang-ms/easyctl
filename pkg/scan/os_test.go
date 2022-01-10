package scan

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"github.com/weiliang-ms/easyctl/pkg/scan/mocks"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"os"
	"sort"
	"testing"
)

var (
	mockHostName = "nodeA"
	mockKernelV  = "3.10.0-693.2.2.el7.x86_64"
	mockOSV      = "CentOS Linux release 7.4.1708 (Core)"
	mockLogger   = logrus.New()

	mockCPUModeNum     = "Intel(R) Xeon(R) Silver 4214 CPU"
	mockCPUThreadCount = 4
	mockCPUClockSpeed  = "2.20GHz"

	mockCPULoadAverage = "3.0 2.0 1.0"

	mockMemUsed       = 0.66
	mockMemUsePercent = 8.64
	mockMemTotal      = 7.64

	mockRootMountUsedPercent       = 98
	mockHighUsedPercentMountPoints = "/,/dev/shm"
	mockErr                        = fmt.Errorf("mock error")

	mockCPUInfo = `
processor       : 0
vendor_id       : GenuineIntel
cpu family      : 6
model           : 85
model name      : Intel(R) Xeon(R) Silver 4214 CPU @ 2.20GHz
stepping        : 7
microcode       : 0x1
cpu MHz         : 2194.842
cache size      : 16896 KB
physical id     : 0
siblings        : 4
core id         : 0
cpu cores       : 2
apicid          : 0
initial apicid  : 0
fpu             : yes
fpu_exception   : yes
cpuid level     : 13
wp              : yes
flags           : fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush mmx fxsr sse sse2 ss ht syscall nx pdpe1gb rdtscp lm constant_tsc rep_good nopl eagerfpu pni pclmulqdq ssse3 fma cx16 pcid sse4_1 sse4_2 x2apic movbe popcnt tsc_deadline_timer aes xsave avx f16c rdrand hypervisor lahf_lm abm 3dnowprefetch fsgsbase tsc_adjust bmi1 hle avx2 smep bmi2 erms invpcid rtm mpx avx512f avx512dq rdseed adx smap avx512cd avx512bw avx512vl xsaveopt xsavec xgetbv1 arat
bogomips        : 4389.68
clflush size    : 64
cache_alignment : 64
address sizes   : 46 bits physical, 48 bits virtual
power management:

processor       : 1
vendor_id       : GenuineIntel
cpu family      : 6
model           : 85
model name      : Intel(R) Xeon(R) Silver 4214 CPU @ 2.20GHz
stepping        : 7
microcode       : 0x1
cpu MHz         : 2194.842
cache size      : 16896 KB
physical id     : 0
siblings        : 4
core id         : 0
cpu cores       : 2
apicid          : 1
initial apicid  : 1
fpu             : yes
fpu_exception   : yes
cpuid level     : 13
wp              : yes
flags           : fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush mmx fxsr sse sse2 ss ht syscall nx pdpe1gb rdtscp lm constant_tsc rep_good nopl eagerfpu pni pclmulqdq ssse3 fma cx16 pcid sse4_1 sse4_2 x2apic movbe popcnt tsc_deadline_timer aes xsave avx f16c rdrand hypervisor lahf_lm abm 3dnowprefetch fsgsbase tsc_adjust bmi1 hle avx2 smep bmi2 erms invpcid rtm mpx avx512f avx512dq rdseed adx smap avx512cd avx512bw avx512vl xsaveopt xsavec xgetbv1 arat
bogomips        : 4389.68
clflush size    : 64
cache_alignment : 64
address sizes   : 46 bits physical, 48 bits virtual
power management:

processor       : 2
vendor_id       : GenuineIntel
cpu family      : 6
model           : 85
model name      : Intel(R) Xeon(R) Silver 4214 CPU @ 2.20GHz
stepping        : 7
microcode       : 0x1
cpu MHz         : 2194.842
cache size      : 16896 KB
physical id     : 0
siblings        : 4
core id         : 1
cpu cores       : 2
apicid          : 2
initial apicid  : 2
fpu             : yes
fpu_exception   : yes
cpuid level     : 13
wp              : yes
flags           : fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush mmx fxsr sse sse2 ss ht syscall nx pdpe1gb rdtscp lm constant_tsc rep_good nopl eagerfpu pni pclmulqdq ssse3 fma cx16 pcid sse4_1 sse4_2 x2apic movbe popcnt tsc_deadline_timer aes xsave avx f16c rdrand hypervisor lahf_lm abm 3dnowprefetch fsgsbase tsc_adjust bmi1 hle avx2 smep bmi2 erms invpcid rtm mpx avx512f avx512dq rdseed adx smap avx512cd avx512bw avx512vl xsaveopt xsavec xgetbv1 arat
bogomips        : 4389.68
clflush size    : 64
cache_alignment : 64
address sizes   : 46 bits physical, 48 bits virtual
power management:

processor       : 3
vendor_id       : GenuineIntel
cpu family      : 6
model           : 85
model name      : Intel(R) Xeon(R) Silver 4214 CPU @ 2.20GHz
stepping        : 7
microcode       : 0x1
cpu MHz         : 2194.842
cache size      : 16896 KB
physical id     : 0
siblings        : 4
core id         : 1
cpu cores       : 2
apicid          : 3
initial apicid  : 3
fpu             : yes
fpu_exception   : yes
cpuid level     : 13
wp              : yes
flags           : fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush mmx fxsr sse sse2 ss ht syscall nx pdpe1gb rdtscp lm constant_tsc rep_good nopl eagerfpu pni pclmulqdq ssse3 fma cx16 pcid sse4_1 sse4_2 x2apic movbe popcnt tsc_deadline_timer aes xsave avx f16c rdrand hypervisor lahf_lm abm 3dnowprefetch fsgsbase tsc_adjust bmi1 hle avx2 smep bmi2 erms invpcid rtm mpx avx512f avx512dq rdseed adx smap avx512cd avx512bw avx512vl xsaveopt xsavec xgetbv1 arat
bogomips        : 4389.68
clflush size    : 64
cache_alignment : 64
address sizes   : 46 bits physical, 48 bits virtual
power management:`
	mockMemoryInfo = `
MemTotal:        8010192 kB
MemFree:         6411000 kB
MemAvailable:    7313528 kB
Buffers:          152044 kB
Cached:           943412 kB
SwapCached:            0 kB
Active:           802048 kB
Inactive:         583496 kB
Active(anon):     292088 kB
Inactive(anon):      736 kB
Active(file):     509960 kB
Inactive(file):   582760 kB
Unevictable:           0 kB
Mlocked:               0 kB
SwapTotal:             0 kB
SwapFree:              0 kB
Dirty:               124 kB
Writeback:             0 kB
AnonPages:        290096 kB
Mapped:           118764 kB
Shmem:              2736 kB
Slab:              97704 kB
SReclaimable:      62628 kB
SUnreclaim:        35076 kB
KernelStack:        4080 kB
PageTables:        26616 kB
NFS_Unstable:          0 kB
Bounce:                0 kB
WritebackTmp:          0 kB
CommitLimit:     4005096 kB
Committed_AS:    1072496 kB
VmallocTotal:   34359738367 kB
VmallocUsed:       21436 kB
VmallocChunk:   34359707388 kB
HardwareCorrupted:     0 kB
AnonHugePages:     94208 kB
HugePages_Total:       0
HugePages_Free:        0
HugePages_Rsvd:        0
HugePages_Surp:        0
Hugepagesize:       2048 kB
DirectMap4k:       61312 kB
DirectMap2M:     3084288 kB
DirectMap1G:     7340032 kB`
	mocMountPointInfo = `
/dev/vda1        40G  3.4G   34G   98% /
devtmpfs        3.9G     0  3.9G   0% /dev
tmpfs           3.9G     0  3.9G   91% /dev/shm`
)

func Test_GetOSInfo_Mock(t *testing.T) {
	mockInterface := &mocks.HandleOSInterface{}

	r := osScanner{
		Logger:           mockLogger,
		HandlerInterface: mockInterface,
	}

	s := runner.ServerInternal{
		Host:     "1.1.1.1",
		Port:     "22",
		UserName: "root",
		Password: "123456",
	}

	mockInterface.On("GetHostName", s, mockLogger).Return(mockHostName, nil)
	mockInterface.On("GetKernelVersion", s, mockLogger).Return(mockKernelV, nil)
	mockInterface.On("GetSystemVersion", s, mockLogger).Return(mockOSV, nil)
	mockInterface.On("GetCPUInfo", s, mockLogger).Return(mockCPUInfo, nil)
	mockInterface.On("GetCPULoadAverage", s, mockLogger).Return(mockCPULoadAverage, nil)
	mockInterface.On("GetMemoryInfo", s, mockLogger).Return(mockMemoryInfo, nil)
	mockInterface.On("GetMountPointInfo", s, mockLogger).Return(mocMountPointInfo, nil)

	re, err := r.GetOSInfo(s)
	require.Nil(t, err)
	require.Equal(t, mockHostName, re.Hostname)
	require.Equal(t, mockKernelV, re.KernelV)
	require.Equal(t, mockOSV, re.OSV)
	require.Equal(t, mockCPUModeNum, re.CPUModeNum)
	require.Equal(t, mockCPUThreadCount, re.CPUThreadCount)
	require.Equal(t, mockCPUClockSpeed, re.CPUClockSpeed)
	require.Equal(t, mockCPULoadAverage, re.CPULoadAverage)
	require.Equal(t, mockMemUsed, re.MemUsed)
	require.Equal(t, mockMemUsePercent, re.MemUsePercent)
	require.Equal(t, mockMemTotal, re.MemTotal)
}

func Test_GetOSInfo_GetHostNameErr_Mock(t *testing.T) {
	mockInterface := &mocks.HandleOSInterface{}

	r := osScanner{
		Logger:           mockLogger,
		HandlerInterface: mockInterface,
	}
	var s runner.ServerInternal

	mockInterface.On("GetHostName", s, mockLogger).Return("", mockErr)
	_, err := r.GetOSInfo(s)
	require.NotNil(t, err)
}

func Test_GetOSInfo_GetKernelVersionErr_Mock(t *testing.T) {
	mockInterface := &mocks.HandleOSInterface{}

	r := osScanner{
		Logger:           mockLogger,
		HandlerInterface: mockInterface,
	}
	var s runner.ServerInternal

	mockInterface.On("GetHostName", s, mockLogger).Return(mockHostName, nil)
	mockInterface.On("GetKernelVersion", s, mockLogger).Return("", mockErr)
	_, err := r.GetOSInfo(s)
	require.NotNil(t, err)
}

func Test_GetOSInfo_GetSystemVersionErr_Mock(t *testing.T) {
	mockInterface := &mocks.HandleOSInterface{}

	r := osScanner{
		Logger:           mockLogger,
		HandlerInterface: mockInterface,
	}
	var s runner.ServerInternal

	mockInterface.On("GetHostName", s, mockLogger).Return(mockHostName, nil)
	mockInterface.On("GetKernelVersion", s, mockLogger).Return(mockKernelV, nil)
	mockInterface.On("GetSystemVersion", s, mockLogger).Return("", mockErr)
	_, err := r.GetOSInfo(s)
	require.NotNil(t, err)
}

func Test_GetOSInfo_GetCPUInfoErr_Mock(t *testing.T) {
	mockInterface := &mocks.HandleOSInterface{}

	r := osScanner{
		Logger:           mockLogger,
		HandlerInterface: mockInterface,
	}
	var s runner.ServerInternal

	mockInterface.On("GetHostName", s, mockLogger).Return(mockHostName, nil)
	mockInterface.On("GetKernelVersion", s, mockLogger).Return(mockKernelV, nil)
	mockInterface.On("GetSystemVersion", s, mockLogger).Return(mockOSV, nil)
	mockInterface.On("GetCPUInfo", s, mockLogger).Return("", mockErr)
	_, err := r.GetOSInfo(s)
	require.NotNil(t, err)
}

func Test_GetOSInfo_GetCPULoadAverageErr_Mock(t *testing.T) {
	mockInterface := &mocks.HandleOSInterface{}

	r := osScanner{
		Logger:           mockLogger,
		HandlerInterface: mockInterface,
	}
	var s runner.ServerInternal

	mockInterface.On("GetHostName", s, mockLogger).Return(mockHostName, nil)
	mockInterface.On("GetKernelVersion", s, mockLogger).Return(mockKernelV, nil)
	mockInterface.On("GetSystemVersion", s, mockLogger).Return(mockOSV, nil)
	mockInterface.On("GetCPUInfo", s, mockLogger).Return(mockCPUInfo, nil)
	mockInterface.On("GetCPULoadAverage", s, mockLogger).Return("", mockErr)
	_, err := r.GetOSInfo(s)
	require.NotNil(t, err)
}

func Test_GetOSInfo_GetMemoryInfoErr_Mock(t *testing.T) {
	mockInterface := &mocks.HandleOSInterface{}

	r := osScanner{
		Logger:           mockLogger,
		HandlerInterface: mockInterface,
	}
	var s runner.ServerInternal

	mockInterface.On("GetHostName", s, mockLogger).Return(mockHostName, nil)
	mockInterface.On("GetKernelVersion", s, mockLogger).Return(mockKernelV, nil)
	mockInterface.On("GetSystemVersion", s, mockLogger).Return(mockOSV, nil)
	mockInterface.On("GetCPUInfo", s, mockLogger).Return(mockCPUInfo, nil)
	mockInterface.On("GetCPULoadAverage", s, mockLogger).Return(mockCPULoadAverage, nil)
	mockInterface.On("GetMemoryInfo", s, mockLogger).Return("", mockErr)
	_, err := r.GetOSInfo(s)
	require.NotNil(t, err)
}

func Test_GetOSInfo_GetMountPointInfoErr_Mock(t *testing.T) {
	mockInterface := &mocks.HandleOSInterface{}

	r := osScanner{
		Logger:           mockLogger,
		HandlerInterface: mockInterface,
	}
	var s runner.ServerInternal

	mockInterface.On("GetHostName", s, mockLogger).Return(mockHostName, nil)
	mockInterface.On("GetKernelVersion", s, mockLogger).Return(mockKernelV, nil)
	mockInterface.On("GetSystemVersion", s, mockLogger).Return(mockOSV, nil)
	mockInterface.On("GetCPUInfo", s, mockLogger).Return(mockCPUInfo, nil)
	mockInterface.On("GetCPULoadAverage", s, mockLogger).Return(mockCPULoadAverage, nil)
	mockInterface.On("GetMemoryInfo", s, mockLogger).Return(mockMemoryInfo, nil)
	mockInterface.On("GetMountPointInfo", s, mockLogger).Return("", mockErr)
	_, err := r.GetOSInfo(s)
	require.NotNil(t, err)
}

func TestNewCPUInfoItem(t *testing.T) {
	c := NewCPUInfoItem(mockCPUInfo)

	require.Equal(t, mockCPUThreadCount, c.CPUThreadCount)
	require.Equal(t, mockCPUModeNum, c.CPUModeNum)
	require.Equal(t, mockCPUClockSpeed, c.CPUClockSpeed)

}

func TestNewMemInfoItem(t *testing.T) {
	m := NewMemInfoItem(mockMemoryInfo)
	require.Equal(t, mockMemUsed, m.MemUsed)
	require.Equal(t, mockMemUsePercent, m.MemUsePercent)
	require.Equal(t, mockMemTotal, m.MemTotal)
}

func TestNewDiskInfoMeta(t *testing.T) {
	content := "/dev/vda1        40G  3.4G   34G   9% /"
	d := NewDiskInfoMetaItem(content)

	assert.Equal(t, "/dev/vda1", d.Filesystem)
	require.Equal(t, "40G", d.Size)
	require.Equal(t, "3.4G", d.Used)
	require.Equal(t, "34G", d.Avail)
	require.Equal(t, "9%", d.UsedPercent)
	require.Equal(t, "/", d.MountedPoint)

	r := NewDiskInfoMetaItem("/dev/vda1        40G  3.4G   34G   9%")
	require.Equal(t, DiskInfoMeta{}, r)
}

func TestNewDiskInfoItem(t *testing.T) {

	d := NewDiskInfoItem(mocMountPointInfo)
	require.Equal(t, mockRootMountUsedPercent, d.RootMountUsedPercent)
	require.Equal(t, mockHighUsedPercentMountPoints, d.HighUsedPercentMountPoint)
}

func TestSaveAsExcel(t *testing.T) {

	defer os.Remove(preserveFileName)

	data := []OSInfo{{
		BaseOSInfo: BaseOSInfo{
			Address:  "192.168.11.1",
			Hostname: "node1",
			KernelV:  "5.15.11-1.el7.elrepo.x86_64",
			OSV:      "CentOS Linux release 7.9.2009 (Core)",
		},
	}, {
		BaseOSInfo: BaseOSInfo{
			Address:  "192.168.11.2",
			Hostname: "node2",
			KernelV:  "5.15.11-1.el7.elrepo.x86_64",
			OSV:      "CentOS Linux release 7.9.2009 (Core)",
		},
	}}

	if err := SaveAsExcel(data); err != nil {
		panic(err)
	}
}

func Test_OS_Mock(t *testing.T) {
	defer os.Remove(preserveFileName)

	content := `
server:
excludes:
 - 192.168.235.132
`
	mockInterface := &mocks.HandleOSInterface{}
	OS(command.OperationItem{
		B:          []byte(content),
		Logger:     logrus.New(),
		OptionFunc: nil,
		Interface:  mockInterface,
		UnitTest:   false,
		Mock:       false,
		Local:      false,
	})
}

func Test_Sort(t *testing.T) {

	var slice OSInfoSlice

	slice = append(slice, OSInfo{BaseOSInfo: BaseOSInfo{Address: "192.168.1.200"}})
	slice = append(slice, OSInfo{BaseOSInfo: BaseOSInfo{Address: "192.168.3.200"}})
	slice = append(slice, OSInfo{BaseOSInfo: BaseOSInfo{Address: "192.168.2.200"}})
	slice = append(slice, OSInfo{BaseOSInfo: BaseOSInfo{Address: "192.168.2.100"}})

	sort.Sort(slice)
}

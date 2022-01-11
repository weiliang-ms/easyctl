package scan

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"github.com/weiliang-ms/easyctl/pkg/scan/mocks"
	"github.com/weiliang-ms/easyctl/pkg/scan/scanos"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"os"
	"testing"
)

var (
	MockHostName = "nodeA"
	MockKernelV  = "3.10.0-693.2.2.el7.x86_64"
	MockOSV      = "CentOS Linux release 7.4.1708 (Core)"

	MockCPUModeNum     = "Intel(R) Xeon(R) Silver 4214 CPU"
	MockCPUThreadCount = 4
	MockCPUClockSpeed  = "2.20GHz"

	MockCPULoadAverage = "3.0 2.0 1.0"

	MockMemUsed       = 0.66
	MockMemUsePercent = 8.64
	MockMemTotal      = 7.64

	MockRootMountUsedPercent       = 98
	MockHighUsedPercentMountPoints = "/,/dev/shm"

	MockCPUInfo = `
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
	MockMemoryInfo = `
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
	MockMountPointInfo = `
/dev/vda1        40G  3.4G   34G   98% /
devtmpfs        3.9G     0  3.9G   0% /dev
tmpfs           3.9G     0  3.9G   91% /dev/shm`
)

var (
	mockServerA = runner.ServerInternal{
		Host:     "192.168.1.1",
		Port:     "22",
		UserName: "root",
		Password: "1",
	}
	mockServerB = runner.ServerInternal{
		Host:     "192.168.1.2",
		Port:     "22",
		UserName: "root",
		Password: "1",
	}
	mockServers = []runner.ServerInternal{
		mockServerA, mockServerB,
	}

	mockServerAHostName = "nodeA"
	mockServerBHostName = "nodeB"

	mockScanOSMetaInfoA = scanos.MetaInfo{BaseOSInfo: scanos.BaseOSInfo{
		Address: mockServerA.Host,
	}}

	mockScanOSMetaInfoB = scanos.MetaInfo{BaseOSInfo: scanos.BaseOSInfo{
		Address: mockServerB.Host,
	}}

	mockLogger = logrus.New()
	mockErr    = fmt.Errorf("mock error")
)

func Test_ParallelGetOSInfo_Mock(t *testing.T) {
	mockInterface := &mocks.HandleOSInterface{}

	r := osScanner{
		Logger:           mockLogger,
		HandlerInterface: mockInterface,
	}

	mockInterface.On("GetOSInfo", mockServerA, mockLogger).Return(mockScanOSMetaInfoA, nil)
	mockInterface.On("GetOSInfo", mockServerB, mockLogger).Return(mockScanOSMetaInfoB, nil)

	mockInterface.On("GetHostName", mockServerA, mockLogger).Return(mockServerAHostName, nil)
	mockInterface.On("GetHostName", mockServerB, mockLogger).Return(mockServerBHostName, nil)

	mockInterface.On("GetKernelVersion", mockServerA, mockLogger).Return(MockKernelV, nil)
	mockInterface.On("GetKernelVersion", mockServerB, mockLogger).Return(MockKernelV, nil)

	mockInterface.On("GetSystemVersion", mockServerA, mockLogger).Return(MockOSV, nil)
	mockInterface.On("GetSystemVersion", mockServerB, mockLogger).Return(MockOSV, nil)

	mockInterface.On("GetCPUInfo", mockServerA, mockLogger).Return(MockCPUInfo, nil)
	mockInterface.On("GetCPUInfo", mockServerB, mockLogger).Return(MockCPUInfo, nil)

	mockInterface.On("GetCPULoadAverage", mockServerA, mockLogger).Return(MockCPULoadAverage, nil)
	mockInterface.On("GetCPULoadAverage", mockServerB, mockLogger).Return(MockCPULoadAverage, nil)

	mockInterface.On("GetMemoryInfo", mockServerA, mockLogger).Return(MockMemoryInfo, nil)
	mockInterface.On("GetMemoryInfo", mockServerB, mockLogger).Return(MockMemoryInfo, nil)

	mockInterface.On("GetMountPointInfo", mockServerA, mockLogger).Return(MockMountPointInfo, nil)
	mockInterface.On("GetMountPointInfo", mockServerB, mockLogger).Return(MockMountPointInfo, nil)

	result, err := r.ParallelGetOSInfo(mockServers)
	require.Nil(t, err)
	require.Equal(t, len(mockServers), len(result))
}

func Test_ParallelGetOSInfo_Err_Mock(t *testing.T) {
	mockInterface := &mocks.HandleOSInterface{}

	r := osScanner{
		Logger:           mockLogger,
		HandlerInterface: mockInterface,
	}

	mockInterface.On("GetOSInfo", mockServerA, mockLogger).Return(scanos.MetaInfo{}, mockErr)
	mockInterface.On("GetOSInfo", mockServerB, mockLogger).Return(scanos.MetaInfo{}, mockErr)

	mockInterface.On("GetHostName", mockServerA, mockLogger).Return("", mockErr)
	mockInterface.On("GetHostName", mockServerB, mockLogger).Return("", mockErr)

	result, err := r.ParallelGetOSInfo(mockServers)
	require.NotNil(t, err)
	var expect scanos.MetaInfoSlice
	require.Equal(t, expect, result)
}

func Test_OS_Mock(t *testing.T) {
	defer os.Remove(preserveFileName)

	content := `
server:
excludes:
- 192.168.235.132
`
	mockInterface := &mocks.HandleOSInterface{}
	s := runner.ServerInternal{
		Host:     "1.1.1.1",
		Port:     "22",
		UserName: "root",
		Password: "123456",
	}

	mockInterface.On("GetHostName", s, mockLogger).Return(MockHostName, nil)
	mockInterface.On("GetKernelVersion", s, mockLogger).Return(MockKernelV, nil)
	mockInterface.On("GetSystemVersion", s, mockLogger).Return(MockOSV, nil)
	mockInterface.On("GetCPUInfo", s, mockLogger).Return(MockCPUInfo, nil)
	mockInterface.On("GetCPULoadAverage", s, mockLogger).Return(MockCPULoadAverage, nil)
	mockInterface.On("GetMemoryInfo", s, mockLogger).Return(MockMemoryInfo, nil)
	mockInterface.On("GetMountPointInfo", s, mockLogger).Return(MockMountPointInfo, nil)

	err := OS(command.OperationItem{
		B:          []byte(content),
		Logger:     logrus.New(),
		OptionFunc: nil,
		Interface:  mockInterface,
		UnitTest:   false,
		Mock:       false,
		Local:      false,
	})

	require.NotNil(t, err)
}

func Test_OS_ParseErr_Mock(t *testing.T) {
	defer os.Remove(preserveFileName)

	content := `
server:
  - host: ddd
  - name: ddd
excludes:
- 192.168.235.132
`
	mockInterface := &mocks.HandleOSInterface{}

	err := OS(command.OperationItem{
		B:          []byte(content),
		Logger:     logrus.New(),
		OptionFunc: nil,
		Interface:  mockInterface,
		UnitTest:   false,
		Mock:       false,
		Local:      false,
	})

	require.NotNil(t, err)
}

func Test_OS_GetHostNameErr_Mock(t *testing.T) {
	defer os.Remove(preserveFileName)

	content := `
server:
 - host: "1.1.1.1"
   port: "22"
   username: root
   password: 1
 - host: "1.1.1.2"
   port: "22"
   username: root
   password: 1
excludes:
- 192.168.235.132
`
	mockInterface := &mocks.HandleOSInterface{}
	s1 := runner.ServerInternal{
		Host:     "1.1.1.1",
		Port:     "22",
		UserName: "root",
		Password: "1",
	}

	s2 := runner.ServerInternal{
		Host:     "1.1.1.2",
		Port:     "22",
		UserName: "root",
		Password: "1",
	}

	mockInterface.On("GetHostName", s1, mockLogger).Return("", mockErr)
	mockInterface.On("GetHostName", s2, mockLogger).Return("", mockErr)

	err := OS(command.OperationItem{
		B:          []byte(content),
		Logger:     mockLogger,
		OptionFunc: nil,
		Interface:  mockInterface,
		UnitTest:   false,
		Mock:       false,
		Local:      false,
	})

	require.NotNil(t, err)
}

func Test_OS_GetKernelVersionErr_Mock(t *testing.T) {
	defer os.Remove(preserveFileName)

	content := `
server:
 - host: "1.1.1.1"
   port: "22"
   username: root
   password: 1
 - host: "1.1.1.2"
   port: "22"
   username: root
   password: 1
excludes:
- 192.168.235.132
`
	mockInterface := &mocks.HandleOSInterface{}
	s1 := runner.ServerInternal{
		Host:     "1.1.1.1",
		Port:     "22",
		UserName: "root",
		Password: "1",
	}

	s2 := runner.ServerInternal{
		Host:     "1.1.1.2",
		Port:     "22",
		UserName: "root",
		Password: "1",
	}

	mockInterface.On("GetHostName", s1, mockLogger).Return(mockServerAHostName, nil)
	mockInterface.On("GetHostName", s2, mockLogger).Return(mockServerBHostName, nil)
	mockInterface.On("GetKernelVersion", s1, mockLogger).Return("", mockErr)
	mockInterface.On("GetKernelVersion", s2, mockLogger).Return("", mockErr)

	err := OS(command.OperationItem{
		B:          []byte(content),
		Logger:     mockLogger,
		OptionFunc: nil,
		Interface:  mockInterface,
		UnitTest:   false,
		Mock:       false,
		Local:      false,
	})

	require.NotNil(t, err)
}

func Test_OS_GetSystemVersionErr_Mock(t *testing.T) {
	defer os.Remove(preserveFileName)

	content := `
server:
 - host: "1.1.1.1"
   port: "22"
   username: root
   password: 1
 - host: "1.1.1.2"
   port: "22"
   username: root
   password: 1
excludes:
- 192.168.235.132
`
	mockInterface := &mocks.HandleOSInterface{}
	s1 := runner.ServerInternal{
		Host:     "1.1.1.1",
		Port:     "22",
		UserName: "root",
		Password: "1",
	}

	s2 := runner.ServerInternal{
		Host:     "1.1.1.2",
		Port:     "22",
		UserName: "root",
		Password: "1",
	}

	mockInterface.On("GetHostName", s1, mockLogger).Return(mockServerAHostName, nil)
	mockInterface.On("GetHostName", s2, mockLogger).Return(mockServerBHostName, nil)
	mockInterface.On("GetKernelVersion", s1, mockLogger).Return(MockKernelV, nil)
	mockInterface.On("GetKernelVersion", s2, mockLogger).Return(MockKernelV, nil)
	mockInterface.On("GetSystemVersion", s1, mockLogger).Return("", mockErr)
	mockInterface.On("GetSystemVersion", s2, mockLogger).Return("", mockErr)

	err := OS(command.OperationItem{
		B:          []byte(content),
		Logger:     mockLogger,
		OptionFunc: nil,
		Interface:  mockInterface,
		UnitTest:   false,
		Mock:       false,
		Local:      false,
	})

	require.NotNil(t, err)
}

func Test_OS_GetCPUInfoErr_Mock(t *testing.T) {
	defer os.Remove(preserveFileName)

	content := `
server:
 - host: "1.1.1.1"
   port: "22"
   username: root
   password: 1
 - host: "1.1.1.2"
   port: "22"
   username: root
   password: 1
excludes:
- 192.168.235.132
`
	mockInterface := &mocks.HandleOSInterface{}
	s1 := runner.ServerInternal{
		Host:     "1.1.1.1",
		Port:     "22",
		UserName: "root",
		Password: "1",
	}

	s2 := runner.ServerInternal{
		Host:     "1.1.1.2",
		Port:     "22",
		UserName: "root",
		Password: "1",
	}

	mockInterface.On("GetHostName", s1, mockLogger).Return(mockServerAHostName, nil)
	mockInterface.On("GetHostName", s2, mockLogger).Return(mockServerBHostName, nil)
	mockInterface.On("GetKernelVersion", s1, mockLogger).Return(MockKernelV, nil)
	mockInterface.On("GetKernelVersion", s2, mockLogger).Return(MockKernelV, nil)
	mockInterface.On("GetSystemVersion", s1, mockLogger).Return(MockOSV, nil)
	mockInterface.On("GetSystemVersion", s2, mockLogger).Return(MockOSV, nil)
	mockInterface.On("GetCPUInfo", s1, mockLogger).Return("", mockErr)
	mockInterface.On("GetCPUInfo", s2, mockLogger).Return("", mockErr)

	err := OS(command.OperationItem{
		B:          []byte(content),
		Logger:     mockLogger,
		OptionFunc: nil,
		Interface:  mockInterface,
		UnitTest:   false,
		Mock:       false,
		Local:      false,
	})

	require.NotNil(t, err)
}

func Test_OS_GetCPULoadAverageErr_Mock(t *testing.T) {
	defer os.Remove(preserveFileName)

	content := `
server:
 - host: "1.1.1.1"
   port: "22"
   username: root
   password: 1
 - host: "1.1.1.2"
   port: "22"
   username: root
   password: 1
excludes:
- 192.168.235.132
`
	mockInterface := &mocks.HandleOSInterface{}
	s1 := runner.ServerInternal{
		Host:     "1.1.1.1",
		Port:     "22",
		UserName: "root",
		Password: "1",
	}

	s2 := runner.ServerInternal{
		Host:     "1.1.1.2",
		Port:     "22",
		UserName: "root",
		Password: "1",
	}

	mockInterface.On("GetHostName", s1, mockLogger).Return(mockServerAHostName, nil)
	mockInterface.On("GetHostName", s2, mockLogger).Return(mockServerBHostName, nil)
	mockInterface.On("GetKernelVersion", s1, mockLogger).Return(MockKernelV, nil)
	mockInterface.On("GetKernelVersion", s2, mockLogger).Return(MockKernelV, nil)
	mockInterface.On("GetSystemVersion", s1, mockLogger).Return(MockOSV, nil)
	mockInterface.On("GetSystemVersion", s2, mockLogger).Return(MockOSV, nil)
	mockInterface.On("GetCPUInfo", s1, mockLogger).Return(MockCPUInfo, nil)
	mockInterface.On("GetCPUInfo", s2, mockLogger).Return(MockCPUInfo, nil)
	mockInterface.On("GetCPULoadAverage", s1, mockLogger).Return("", mockErr)
	mockInterface.On("GetCPULoadAverage", s2, mockLogger).Return("", mockErr)
	//mockInterface.On("GetMemoryInfo", s1, mockLogger).Return(MockMemoryInfo, nil)
	//mockInterface.On("GetMountPointInfo", s1, mockLogger).Return(MockMountPointInfo, nil)

	err := OS(command.OperationItem{
		B:          []byte(content),
		Logger:     mockLogger,
		OptionFunc: nil,
		Interface:  mockInterface,
		UnitTest:   false,
		Mock:       false,
		Local:      false,
	})

	require.NotNil(t, err)
}

func Test_OS_GetMemoryInfoErr_Mock(t *testing.T) {
	defer os.Remove(preserveFileName)

	content := `
server:
 - host: "1.1.1.1"
   port: "22"
   username: root
   password: 1
 - host: "1.1.1.2"
   port: "22"
   username: root
   password: 1
excludes:
- 192.168.235.132
`
	mockInterface := &mocks.HandleOSInterface{}
	s1 := runner.ServerInternal{
		Host:     "1.1.1.1",
		Port:     "22",
		UserName: "root",
		Password: "1",
	}

	s2 := runner.ServerInternal{
		Host:     "1.1.1.2",
		Port:     "22",
		UserName: "root",
		Password: "1",
	}

	mockInterface.On("GetHostName", s1, mockLogger).Return(mockServerAHostName, nil)
	mockInterface.On("GetHostName", s2, mockLogger).Return(mockServerBHostName, nil)
	mockInterface.On("GetKernelVersion", s1, mockLogger).Return(MockKernelV, nil)
	mockInterface.On("GetKernelVersion", s2, mockLogger).Return(MockKernelV, nil)
	mockInterface.On("GetSystemVersion", s1, mockLogger).Return(MockOSV, nil)
	mockInterface.On("GetSystemVersion", s2, mockLogger).Return(MockOSV, nil)
	mockInterface.On("GetCPUInfo", s1, mockLogger).Return(MockCPUInfo, nil)
	mockInterface.On("GetCPUInfo", s2, mockLogger).Return(MockCPUInfo, nil)
	mockInterface.On("GetCPULoadAverage", s1, mockLogger).Return(MockCPULoadAverage, nil)
	mockInterface.On("GetCPULoadAverage", s2, mockLogger).Return(MockCPULoadAverage, nil)
	mockInterface.On("GetMemoryInfo", s1, mockLogger).Return("", mockErr)
	mockInterface.On("GetMemoryInfo", s2, mockLogger).Return("", mockErr)

	err := OS(command.OperationItem{
		B:          []byte(content),
		Logger:     mockLogger,
		OptionFunc: nil,
		Interface:  mockInterface,
		UnitTest:   false,
		Mock:       false,
		Local:      false,
	})

	require.NotNil(t, err)
}

func Test_OS_GetMountPointInfoErr_Mock(t *testing.T) {
	defer os.Remove(preserveFileName)

	content := `
server:
- host: "1.1.1.1"
  port: "22"
  username: root
  password: 1
- host: "1.1.1.2"
  port: "22"
  username: root
  password: 1
excludes:
- 192.168.235.132
`
	mockInterface := &mocks.HandleOSInterface{}
	s1 := runner.ServerInternal{
		Host:     "1.1.1.1",
		Port:     "22",
		UserName: "root",
		Password: "1",
	}

	s2 := runner.ServerInternal{
		Host:     "1.1.1.2",
		Port:     "22",
		UserName: "root",
		Password: "1",
	}

	mockInterface.On("GetHostName", s1, mockLogger).Return(mockServerAHostName, nil)
	mockInterface.On("GetHostName", s2, mockLogger).Return(mockServerBHostName, nil)
	mockInterface.On("GetKernelVersion", s1, mockLogger).Return(MockKernelV, nil)
	mockInterface.On("GetKernelVersion", s2, mockLogger).Return(MockKernelV, nil)
	mockInterface.On("GetSystemVersion", s1, mockLogger).Return(MockOSV, nil)
	mockInterface.On("GetSystemVersion", s2, mockLogger).Return(MockOSV, nil)
	mockInterface.On("GetCPUInfo", s1, mockLogger).Return(MockCPUInfo, nil)
	mockInterface.On("GetCPUInfo", s2, mockLogger).Return(MockCPUInfo, nil)
	mockInterface.On("GetCPULoadAverage", s1, mockLogger).Return(MockCPULoadAverage, nil)
	mockInterface.On("GetCPULoadAverage", s2, mockLogger).Return(MockCPULoadAverage, nil)
	mockInterface.On("GetMemoryInfo", s1, mockLogger).Return(MockMemoryInfo, nil)
	mockInterface.On("GetMemoryInfo", s2, mockLogger).Return(MockMemoryInfo, nil)
	mockInterface.On("GetMountPointInfo", s1, mockLogger).Return("", mockErr)
	mockInterface.On("GetMountPointInfo", s2, mockLogger).Return("", mockErr)

	err := OS(command.OperationItem{
		B:          []byte(content),
		Logger:     mockLogger,
		OptionFunc: nil,
		Interface:  mockInterface,
		UnitTest:   false,
		Mock:       false,
		Local:      false,
	})

	require.NotNil(t, err)
}

func Test_getHandlerInterface(t *testing.T) {
	var h HandleOSInterface
	r := getHandlerInterface(h)
	require.Equal(t, new(OsExecutor), r)

	h2 := &mocks.HandleOSInterface{}
	r2 := getHandlerInterface(h2)
	require.Equal(t, h2, r2)
}

func TestSaveAsExcel(t *testing.T) {

	defer os.Remove(preserveFileName)

	data := []scanos.MetaInfo{{
		BaseOSInfo: scanos.BaseOSInfo{
			Address:  "192.168.11.1",
			Hostname: "node1",
			KernelV:  "5.15.11-1.el7.elrepo.x86_64",
			OSV:      "CentOS Linux release 7.9.2009 (Core)",
		},
	}, {
		BaseOSInfo: scanos.BaseOSInfo{
			Address:  "192.168.11.2",
			Hostname: "node2",
			KernelV:  "5.15.11-1.el7.elrepo.x86_64",
			OSV:      "CentOS Linux release 7.9.2009 (Core)",
		},
	}}

	if err := SaveAsExcel(data, preserveFileName, excelTmpl); err != nil {
		panic(err)
	}
}

func TestSaveAsExcel_ErrCase(t *testing.T) {

	defer os.Remove(preserveFileName)

	invalidPath := "ddd/ddd/ddd.xlsx"
	os.Create(invalidPath)

	data := []scanos.MetaInfo{{
		BaseOSInfo: scanos.BaseOSInfo{
			Address:  "192.168.11.1",
			Hostname: "node1",
			KernelV:  "5.15.11-1.el7.elrepo.x86_64",
			OSV:      "CentOS Linux release 7.9.2009 (Core)",
		},
	}, {
		BaseOSInfo: scanos.BaseOSInfo{
			Address:  "192.168.11.2",
			Hostname: "node2",
			KernelV:  "5.15.11-1.el7.elrepo.x86_64",
			OSV:      "CentOS Linux release 7.9.2009 (Core)",
		},
	}}

	require.NotNil(t, SaveAsExcel(data, invalidPath, excelTmpl))
}

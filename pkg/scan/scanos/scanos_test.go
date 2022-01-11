package scanos

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sort"
	"testing"
)

func Test_Sort(t *testing.T) {

	var slice MetaInfoSlice

	slice = append(slice, MetaInfo{BaseOSInfo: BaseOSInfo{Address: "192.168.1.200"}})
	slice = append(slice, MetaInfo{BaseOSInfo: BaseOSInfo{Address: "192.168.3.200"}})
	slice = append(slice, MetaInfo{BaseOSInfo: BaseOSInfo{Address: "192.168.2.200"}})
	slice = append(slice, MetaInfo{BaseOSInfo: BaseOSInfo{Address: "192.168.2.100"}})

	sort.Sort(slice)
}

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
	MocMountPointInfo = `
/dev/vda1        40G  3.4G   34G   98% /
devtmpfs        3.9G     0  3.9G   0% /dev
tmpfs           3.9G     0  3.9G   91% /dev/shm`
)

func TestNewCPUInfoItem(t *testing.T) {
	c := NewCPUInfoItem(MockCPUInfo)

	require.Equal(t, MockCPUThreadCount, c.CPUThreadCount)
	require.Equal(t, MockCPUModeNum, c.CPUModeNum)
	require.Equal(t, MockCPUClockSpeed, c.CPUClockSpeed)

}

func TestNewMemInfoItem(t *testing.T) {
	m := NewMemInfoItem(MockMemoryInfo)
	require.Equal(t, MockMemUsed, m.MemUsed)
	require.Equal(t, MockMemUsePercent, m.MemUsePercent)
	require.Equal(t, MockMemTotal, m.MemTotal)
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

	d := NewDiskInfoItem(MocMountPointInfo)
	require.Equal(t, MockRootMountUsedPercent, d.RootMountUsedPercent)
	require.Equal(t, MockHighUsedPercentMountPoints, d.HighUsedPercentMountPoint)
}

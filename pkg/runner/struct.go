package runner

import (
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"sync"
)

// ServerExternal server序列化对象
type ServerExternal struct {
	Host           string `yaml:"host"`
	Port           string `yaml:"port"`
	Username       string `yaml:"username"`
	Password       string `yaml:"password"`
	PrivateKeyPath string `yaml:"privateKeyPath"`
}

// ServerInternal server内部对象
type ServerInternal struct {
	Host           string
	Port           string
	Username       string
	Password       string
	PrivateKeyPath string
}

// ServerListExternal server列表序列化对象
type ServerListExternal struct {
	Servers  []ServerExternal `yaml:"server"`
	Excludes []string         `yaml:"excludes"`
}

// ServerListInternal server列表内部对象
type ServerListInternal struct {
	Servers  []ServerInternal
	Excludes []string
}

// ExecutorExternal 执行器序列化对象
type ExecutorExternal struct {
	Servers  []ServerExternal `yaml:"server"`
	Excludes []string         `yaml:"excludes"`
	Script   string           `yaml:"script"`
	logrus.Logger
}

// ExecutorInternal 执行器内部对象
type ExecutorInternal struct {
	Servers        []ServerInternal
	Script         string
	Logger         *logrus.Logger
	OutPutRealTime bool
	RunOnServer    ServerInternal
	*sync.Mutex
}

// ShellResult shell执行结果
type ShellResult struct {
	Host      string `table:"主机地址"`
	Cmd       string `table:"执行语句"`
	Code      int    `table:"退出码"`
	Status    string `table:"执行状态"`
	StdOut    string `table:"执行结果"`
	Output    string `table:"标准输出"`
	StdErrMsg string
	Err       error
}

// ShellResultSlice shell执行结果切片
type ShellResultSlice []ShellResult

func (re ShellResultSlice) Len() int { return len(re) }

func (re ShellResultSlice) Swap(i, j int) { re[i], re[j] = re[j], re[i] }

func (re ShellResultSlice) Less(i, j int) bool {
	address1 := strings.Split(re[i].Host, ".")
	address2 := strings.Split(re[j].Host, ".")

	for k := 0; k < 4; k++ {
		if address1[k] != address2[k] {
			num1, _ := strconv.Atoi(address1[k])
			num2, _ := strconv.Atoi(address2[k])
			return num1 < num2
		}
	}

	return true
}

// InternelServersSlice 带有排序的server列表
type InternelServersSlice []ServerInternal

func (servers InternelServersSlice) Len() int { return len(servers) }

func (servers InternelServersSlice) Swap(i, j int) { servers[i], servers[j] = servers[j], servers[i] }

func (servers InternelServersSlice) Less(i, j int) bool {
	address1 := strings.Split(servers[i].Host, ".")
	address2 := strings.Split(servers[j].Host, ".")

	for k := 0; k < 4; k++ {
		if address1[k] != address2[k] {
			num1, _ := strconv.Atoi(address1[k])
			num2, _ := strconv.Atoi(address2[k])
			return num1 < num2
		}
	}

	return true
}

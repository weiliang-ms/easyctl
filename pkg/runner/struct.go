package runner

import (
	"strconv"
	"strings"
)

type ServerExternal struct {
	Host           string `yaml:"host"`
	Port           string `yaml:"port"`
	Username       string `yaml:"username"`
	Password       string `yaml:"password"`
	PrivateKeyPath string `yaml:"privateKeyPath"`
}

type ServerInternal struct {
	Host           string
	Port           string
	Username       string
	Password       string
	PrivateKeyPath string
}

type ServerListExternal struct {
	Servers  []ServerExternal `yaml:"server"`
	Excludes []string         `yaml:"excludes"`
}

type ServerListInternal struct {
	Servers  []ServerInternal
	Excludes []string
}

type ExecutorExternal struct {
	Servers  []ServerExternal `yaml:"server"`
	Excludes []string         `yaml:"excludes"`
	Script   string           `yaml:"script"`
}

type ExecutorInternal struct {
	Servers []ServerInternal
	Script  string
}

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

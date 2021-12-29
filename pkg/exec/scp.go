package exec

import (
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
	"sync"
	"time"
)

type scpItem struct {
	sync.RWMutex
	SrcDir       string
	SrcName      string
	DstDir       string
	DstName      string
	Mode         string
	TransferType TransferType
	StartAt      time.Time
	StopAt       time.Time
	Server       runner.ServerInternal
}

type TransferType string

const (
	TransferFile TransferType = "file"
	TransferDir  TransferType = "directory"
)

type scpObject struct {
	Scp struct {
		Src  string `yaml:"src"`
		Dst  string `yaml:"dst"`
		Mode string `yaml:"mode"`
	} `yaml:"scp"`
}

// Scp 执行指令
func Scp(item command.OperationItem) command.RunErr {
	executor, err := runner.ParseExecutor(item.B, item.Logger)
	if err != nil {
		return command.RunErr{Err: err}
	}
	return runner.RemoteRun(runner.RemoteRunItem{
		B:                   item.B,
		Logger:              item.Logger,
		Cmd:                 executor.Script,
		RecordErrServerList: false,
	})
}

func parseScpItem(b []byte) (scpItem, command.RunErr) {
	object := scpObject{}
	item := scpItem{}
	err := yaml.Unmarshal(b, &object)

	if err != nil {
		return scpItem{}, command.RunErr{Err: err, Msg: "解析错误"}
	}

	item.Mode = object.Scp.Mode

	f, err := os.Stat(object.Scp.Src)
	if err != nil {
		return scpItem{}, command.RunErr{Err: err, Msg: "源文件/目录不存在"}
	}

	// 拷贝目录
	item.SrcDir = object.Scp.Src
	item.TransferType = TransferDir

	// 拷贝文件
	if !f.IsDir() {
		item.TransferType = TransferFile
		item.SrcName = f.Name()
		dirName := strings.Split(object.Scp.Src, f.Name())
		item.SrcDir = dirName[0]
		item.DstName = object.Scp.Dst
	}

	return item, command.RunErr{}
}

// 合法性检测
func (object scpObject) valid() command.RunErr {

	// 目标文件合法性

	return command.RunErr{}
}

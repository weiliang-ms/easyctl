package add

import (
	"github.com/modood/table"
	"github.com/weiliang-ms/easyctl/pkg/runner"
)

const (
	localhost = "localhost"
	success   = "success"
	fail      = "fail"
)

// 执行器执行结果
type ActuatorExecRe struct {
	Host     string `table:"主机IP"`
	Action   string `table:"操作内容"`
	Status   string `table:"执行状态"`
	ErrorMsg string `table:"错误信息"`
}

// 动作
type Action string

type Actuator struct {
	ServerListFile string
	Cmd            string
	UserName       string
	Password       string
	NoLogin        bool
	ServerList     runner.CommonServerList
	Err            error
}

func (ac *Actuator) parseServerList() *Actuator {
	if ac.ServerListFile == "" {
		return ac
	}
	serverList := runner.ParseCommonServerList(ac.ServerListFile)
	ac.ServerList = serverList
	return ac
}

func (ac *Actuator) execute(action string) {

	var result []ActuatorExecRe
	if len(ac.ServerList.Server) <= 0 {
		re := runner.Shell(ac.Cmd)
		result = append(result, packageRe(localhost, action, re.StdErr, re.ExitCode))
	} else {
		for _, v := range ac.ServerList.Server {
			re := v.RemoteShell(ac.Cmd)
			result = append(result, packageRe(v.Host, action, re.StdErrMsg, re.Code))
		}
	}

	table.OutputA(result)
}

func packageRe(host, action, errMsg string, code int) ActuatorExecRe {

	var status string

	if code == 0 {
		status = success
	} else {
		status = fail
	}
	return ActuatorExecRe{
		Host:     host,
		Action:   action,
		Status:   status,
		ErrorMsg: errMsg,
	}
}

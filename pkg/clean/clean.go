package clean

import (
	"github.com/modood/table"
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"github.com/weiliang-ms/easyctl/pkg/util"
)

// 执行器执行结果
type ActuatorExecRe struct {
	Host     string `table:"主机IP"`
	Action   string `table:"执行操作"`
	Status   string `table:"执行状态"`
	ErrorMsg string `table:"错误信息"`
}

type Actuator struct {
	ServerListFile string
	Cmd            string
	Value          string
	ServerList     runner.CommonServerList
	Err            error
	ExitCode       int
}

func (ac *Actuator) parseServerList() *Actuator {
	if ac.ServerListFile == "" {
		return &Actuator{
			ServerList: runner.CommonServerList{},
		}
	}
	serverList := runner.ParseCommonServerList(ac.ServerListFile)
	ac.ServerList = serverList
	return ac
}

func (ac *Actuator) execute(action string, ignoreErrCode int) {

	var result []ActuatorExecRe
	if len(ac.ServerList.Server) <= 0 {
		re := runner.Shell(ac.Cmd)
		ac.ExitCode = re.ExitCode
		result = append(result, packageRe(util.Localhost, action, re.StdErr, ignoreErrCode, re.ExitCode))
	} else {
		for _, v := range ac.ServerList.Server {
			re := v.RemoteShell(ac.Cmd)
			result = append(result, packageRe(v.Host, action, re.StdErrMsg, ignoreErrCode, re.Code))
		}
	}

	table.OutputA(result)
}

func packageRe(host, action, errMsg string, ignoreErrCode int, code int) ActuatorExecRe {

	var status string

	if code == 0 || code == ignoreErrCode {
		status = util.Success
	} else {
		status = util.Fail
	}
	return ActuatorExecRe{
		Host:     host,
		Action:   action,
		Status:   status,
		ErrorMsg: errMsg,
	}
}

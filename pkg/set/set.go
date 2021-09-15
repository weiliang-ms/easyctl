package set

import (
	"fmt"
	"github.com/modood/table"
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"github.com/weiliang-ms/easyctl/pkg/util"
	"log"
	"net/url"
	"os"
	"strings"
)

const (
	localhost = "localhost"
	success   = "success"
	fail      = "fail"
)

// ActuatorExecRe 执行器执行结果
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
	FilePath       string
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

func (ac *Actuator) download() *Actuator {

	if ac.FilePath == "" {
		return ac
	}
	// 判断是否为 url 类型
	_, err := url.ParseRequestURI(ac.FilePath)
	requestUrl := ac.FilePath
	if err == nil {
		strSlice := strings.Split(ac.FilePath, "/")
		if len(strSlice) > 0 {
			ac.FilePath = strSlice[len(strSlice)-1]
		}
		log.Println(fmt.Sprintf("检测到文件类型为url，尝试下载: %s...", requestUrl))
		err := util.DownloadFile(requestUrl, ac.FilePath)
		if err != nil {
			log.Fatalln(err.Error())
		}
	} else {
		_, err := os.Stat(ac.FilePath)
		if err != nil {
			log.Fatalf(fmt.Sprintf("检测到：%s异常 -> %s", ac.FilePath, err.Error()))
		}
	}

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

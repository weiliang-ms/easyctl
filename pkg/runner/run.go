/*
	MIT License

Copyright (c) 2020 xzx.weiliang

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

*/
package runner

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/sftp"
	"github.com/sirupsen/logrus"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"github.com/weiliang-ms/easyctl/pkg/util/constant"
	"golang.org/x/crypto/ssh"
	"net"
	"os"
	"sort"
	"time"
)

type RemoteRunItem struct {
	B                   []byte
	Logger              *logrus.Logger
	Cmd                 string
	RecordErrServerList bool
}

type RunItem struct {
	Logger *logrus.Logger
	Cmd    string
}

func sftpConnect(user, password, host string, port string) (sftpClient *sftp.Client, err error) { //参数: 远程服务器用户名, 密码, ip, 端口
	auth := make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	var timeout time.Duration

	if os.Getenv(constant.SshNoTimeout) == "true" {
		timeout = 1
	} else {
		timeout = 5
	}
	clientConfig := &ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: timeout * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	addr := host + ":" + port
	sshClient, err := ssh.Dial("tcp", addr, clientConfig) //连接ssh
	if err != nil {
		return nil, fmt.Errorf("连接ssh失败 %s", err)
	}

	if sftpClient, err = sftp.NewClient(sshClient); err != nil { //创建客户端
		return nil, fmt.Errorf("创建客户端失败 %s", err)
	}

	return sftpClient, nil
}

// RemoteRun 远程执行输出结果
func RemoteRun(item RemoteRunItem) command.RunErr {

	var f *os.File

	results, err := GetResult(item.B, item.Logger, item.Cmd)
	if err != nil {
		return command.RunErr{Err: err}
	}
	var data [][]string

	if item.RecordErrServerList {
		f, err = os.Create("error-server-list.txt")
		if err != nil {
			panic(err)
		}
		defer f.Close()
	}

	var errServerList []ShellResult

	for _, v := range results {
		if v.Err != nil && item.RecordErrServerList {
			f.Write([]byte(fmt.Sprintf("%s\n", v.Host)))
		}
		data = append(data, []string{v.Host, v.Cmd, fmt.Sprintf("%d", v.Code), v.Status, v.StdOut, v.StdErrMsg})
	}

	if len(errServerList) > 0 {
		return command.RunErr{Err: errServerList[0].Err, Msg: errServerList[0].StdErrMsg}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"IP ADDRESS", "cmd", "exit code", "result", "output", "exception"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	//table.SetRowLine(true)
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.AppendBulk(data) // Add Bulk Data
	table.Render()

	return command.RunErr{}
}

// GetResult 远程执行，获取结果
func GetResult(b []byte, logger *logrus.Logger, cmd string) ([]ShellResult, error) {

	servers, err := ParseServerList(b, logger)
	if err != nil {
		return []ShellResult{}, err
	}

	// 组装执行器,执行命令
	executor := ExecutorInternal{Servers: servers, Script: cmd, Logger: logger}
	ch := executor.ParallelRun()

	// 打包执行结果
	var results []ShellResult

	for re := range ch {
		var result ShellResult
		_ = mapstructure.Decode(re, &result)
		results = append(results, result)
	}

	sort.Sort(ShellResultSlice(results))

	return results, nil
}

package runner

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/sftp"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"net"
	"os"
	"sort"
	"time"
)

func sftpConnect(user, password, host string, port string) (sftpClient *sftp.Client, err error) { //参数: 远程服务器用户名, 密码, ip, 端口
	auth := make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	clientConfig := &ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: 30 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	addr := host + ":" + port
	sshClient, err := ssh.Dial("tcp", addr, clientConfig) //连接ssh
	if err != nil {
		fmt.Println("连接ssh失败", err)
		return
	}

	if sftpClient, err = sftp.NewClient(sshClient); err != nil { //创建客户端
		fmt.Println("创建客户端失败", err)
		return
	}

	return
}

// RemoteRun 远程执行输出结果
func RemoteRun(b []byte, logger *logrus.Logger, cmd string) error {

	results, err := GetResult(b, logger, cmd)
	if err != nil {
		return err
	}
	var data [][]string

	for _, v := range results {
		data = append(data, []string{v.Host, v.Cmd, fmt.Sprintf("%d", v.Code), v.Status, v.StdOut, v.StdErrMsg})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"IP ADDRESS", "cmd", "exit code", "result", "output", "exception"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	//table.SetRowLine(true)
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.AppendBulk(data) // Add Bulk Data
	table.Render()

	return nil
}

// GetResult 远程执行，获取结果
func GetResult(b []byte, logger *logrus.Logger, cmd string) ([]ShellResult, error) {

	servers, err := ParseServerList(b, nil)
	if err != nil {
		return []ShellResult{}, err
	}

	executor := ExecutorInternal{Servers: servers, Script: cmd, Logger: logger}

	ch := executor.ParallelRun()

	var results []ShellResult

	for re := range ch {
		var result ShellResult
		_ = mapstructure.Decode(re, &result)
		results = append(results, result)
	}

	// todo: ip地址排序
	sort.Sort(ShellResultSlice(results))

	return results, nil
}

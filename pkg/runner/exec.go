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
	"errors"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/weiliang-ms/easyctl/pkg/util/constant"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
)

// Item 执行run指令的对象
type Item struct {
	Server         ServerInternal
	Cmd            string
	Logger         *logrus.Logger
	OutputRealTime bool
}

type WindowsErr struct {
	Errors string
}

// LocalRun 本地执行
func LocalRun(shell string, logger *logrus.Logger) ShellResult {

	result := ShellResult{}

	if logger == nil {
		logger = logrus.New()
	}

	logger.Debugf("执行指令: %s", shell)
	cmd := exec.Command("/bin/bash", "-c", shell)
	b, err := cmd.CombinedOutput()

	if err != nil {
		return ShellResult{Err: err}
	}

	logger.Debug(string(b))

	return result
}

// ParallelRun 并发执行
func (executor ExecutorInternal) ParallelRun() chan ShellResult {

	if executor.Logger == nil {
		executor.Logger = logrus.New()
	}

	executor.Logger.Infoln("开始并行执行命令...")
	wg := sync.WaitGroup{}
	ch := make(chan ShellResult, len(executor.Servers))

	// todo: 屏蔽rm -rf
	// 判断入参为文件还是shell

	if _, err := os.Stat(executor.Script); err == nil {
		b, _ := os.ReadFile(executor.Script)
		executor.Script = string(b)
	}

	for _, v := range executor.Servers {
		wg.Add(1)
		go func(s ServerInternal) {
			ch <- executor.runOnNode(s)
			defer wg.Done()
		}(v)
	}

	wg.Wait()
	close(ch)
	return ch
}

func (executor ExecutorInternal) runOnNode(server ServerInternal) (re ShellResult) {
	var shell string
	defer handleErr(&re.Err)

	if executor.Logger.Level == logrus.DebugLevel {
		shell = executor.Script
	} else {
		shell = "shell content"
	}

	// 截取cmd output 长度
	var subCmd string
	if len(executor.Script) > 10 {
		subCmd = "******"
	} else {
		subCmd = executor.Script
	}

	executor.Logger.Infof("[%s] 开始执行指令 -> %s", server.Host, shell)
	session, err := server.sshConnect()
	defer session.Close()

	if err != nil {
		errMsg := fmt.Sprintf("ssh会话建立失败->%s", err.(*net.OpError).Error())
		return ShellResult{
			Host:      server.Host,
			Err:       err,
			Cmd:       executor.Script,
			Status:    constant.Fail,
			Code:      -1,
			StdErrMsg: errMsg,
		}
	}

	if executor.OutPutRealTime == true {
		session.Stdout = os.Stdout
		err := session.Run(executor.Script)
		if err != nil {
			return ShellResult{
				Host:      server.Host,
				Err:       errors.New(err.(*ssh.ExitError).String()),
				Cmd:       executor.Script,
				Status:    constant.Fail,
				Code:      err.(*ssh.ExitError).ExitStatus(),
				StdErrMsg: err.(*ssh.ExitError).String()}
		}

	} else {
		out, err := session.Output(executor.Script)
		if err != nil {
			return ShellResult{
				Host:      server.Host,
				Err:       errors.New(err.(*ssh.ExitError).String()),
				Cmd:       executor.Script,
				Status:    constant.Fail,
				Code:      err.(*ssh.ExitError).ExitStatus(),
				StdErrMsg: err.(*ssh.ExitError).String()}
		}

		executor.Logger.Infof("<- %s执行命令成功...", server.Host)
		executor.Logger.Debugf("%s -> 返回值: %s", server.Host, string(out))

		var subOut string

		if len(string(out)) > 20 {
			subOut = string(out)[:20]
		} else {
			subOut = string(out)
		}

		return ShellResult{Host: server.Host, StdOut: subOut,
			Cmd: strings.TrimPrefix(subCmd, "\n"), Status: constant.Success}
	}

	return ShellResult{}
}

// ReturnRunResult 获取执行结果
func (server ServerInternal) ReturnRunResult(cmd string) ShellResult {
	log.Printf("<- %s开始执行命令...\n", server.Host)
	session, err := server.sshConnect()
	if err != nil {
		return ShellResult{Err: fmt.Errorf("%s建立ssh会话失败 -> %s", server.Host, err.Error())}
	}

	combo, err := session.CombinedOutput(cmd)
	if err != nil {
		//klog.Fatal("远程执行cmd 失败",err)
		return ShellResult{Err: fmt.Errorf("%s执行失败, %s", server.Host, combo)}
	}
	log.Printf("<- %s执行命令成功，返回结果 => %s...\n", server.Host, string(combo))
	defer session.Close()

	return ShellResult{StdOut: string(combo)}
}

func (server ServerInternal) sshConnect() (*ssh.Session, error) {
	s := server.completeDefault()
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		session      *ssh.Session
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)

	if server.Password == "" || server.PrivateKeyPath != "" {
		a, err := publicKeyAuthFunc(server.PrivateKeyPath)
		if err != nil {
			return session, err
		}
		auth = append(auth, a)
	} else {
		auth = append(auth, ssh.Password(s.Password))
	}

	hostKeyCallbk := func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		return nil
	}

	var timeout time.Duration

	if os.Getenv(constant.SshNoTimeout) == "true" {
		timeout = 1
	} else {
		timeout = 3
	}

	clientConfig = &ssh.ClientConfig{
		User:            s.Username,
		Auth:            auth,
		Timeout:         timeout * time.Second,
		HostKeyCallback: hostKeyCallbk,
	}

	// connet to ssh
	addr = fmt.Sprintf("%s:%s", server.Host, server.Port)

	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	// create session
	if session, err = client.NewSession(); err != nil {
		return nil, err
	}

	return session, nil
}

func publicKeyAuthFunc(kPath string) (ssh.AuthMethod, error) {
	keyPath, err := homedir.Expand(kPath)
	if err != nil {
		return nil, fmt.Errorf("find key's home dir failed %s", err)
	}
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("ssh key file read failed %s", err)
	}
	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("ssh key signer failed %s", err)
	}
	return ssh.PublicKeys(signer), nil
}

// todo: 利用反射断言等赋默认值
func (server ServerInternal) completeDefault() ServerInternal {
	if server.Port == "" {
		server.Port = "22"
	}

	if server.Username == "" {
		server.Username = constant.Root
	}

	if server.PrivateKeyPath == "" {
		server.PrivateKeyPath = constant.RsaPrvPath
	}

	return server
}

func handleErr(err *error) {

	if v := recover(); v != nil {
		if e, ok := v.(runtime.Error); ok {
			*err = e
		} else {
			panic(v)
		}
	}
}

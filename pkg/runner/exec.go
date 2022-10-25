/*
	MIT License

# Copyright (c) 2020 xzx.weiliang

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
	"bytes"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/weiliang-ms/easyctl/pkg/util/constant"
	"github.com/weiliang-ms/easyctl/pkg/util/log"
	"github.com/weiliang-ms/easyctl/pkg/util/value"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
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
	logger = log.SetDefault(logger)

	var cmd *exec.Cmd

	logger.Debugf("[shell] local执行指令: \n%s", shell)

	switch runtime.GOOS {
	case "windows":
		return ShellResult{Err: fmt.Errorf("不支持windows平台")}
	case "linux":
		cmd = exec.Command("/bin/bash", "-c", shell)
	default:
		cmd = exec.Command("/bin/bash", "-c", shell)
	}

	b, err := cmd.CombinedOutput()

	if err != nil {
		logger.Debugf("[shell] 执行错误：%s", err)
		return ShellResult{Err: err}
	}

	logger.Debug(string(b))
	return result
}

// ParallelRun 并发执行
// todo: 添加缓冲队列，避免过多goroutine
func (executor ExecutorInternal) ParallelRun(timeout time.Duration) chan ShellResult {

	if executor.RunShellFunc == nil {
		executor.RunShellFunc = RunOnNode
	}

	executor.Logger = log.SetDefault(executor.Logger)
	executor.Logger.Info("开始并行执行命令...")
	wg := sync.WaitGroup{}
	ch := make(chan ShellResult, len(executor.Servers))

	// todo: 屏蔽rm -rf
	// 判断入参为文件还是shell

	if executor.Script == "rm -rf /" {
		panic(fmt.Errorf("执行命令: %s 存在重大安全隐患", executor.Script))
	}

	if _, err := os.Stat(executor.Script); err == nil {
		b, _ := os.ReadFile(executor.Script)
		executor.Script = string(b)
		executor.Logger.Debug("解析出执行指令为脚本")
	}

	for _, v := range executor.Servers {
		wg.Add(1)
		go func(s ServerInternal) {
			ch <- executor.RunShellFunc(executor.Script, s, timeout, executor.Logger)
			defer wg.Done()
		}(v)
	}

	wg.Wait()
	close(ch)
	return ch
}

func RunOnNode(shell string, server ServerInternal, timeout time.Duration, logger *logrus.Logger) (re ShellResult) {
	defer handleErr(&re.Err)

	logger.Infof("[%s] 开始执行指令 -> start", server.Host)
	logger.Debugf("\n# 指令开始\n%s\n# 指令结束\n", shell)
	session, err := server.SSHConnect(timeout)
	defer session.Close()

	if err != nil {
		errMsg := fmt.Sprintf("%s ssh会话建立失败->%s", server.Host, err.Error())
		return ShellResult{
			Host:      server.Host,
			Err:       err,
			Cmd:       shell,
			Status:    constant.Fail,
			Code:      -1,
			StdErrMsg: errMsg,
		}
	}

	// 是否实时输出
	//if executor.OutPutRealTime == true {
	//	session.Stdout = os.Stdout
	//	var errOut bytes.Buffer
	//	session.Stderr = &errOut
	//
	//	if err := session.Run(shell); err != nil {
	//		code := err.(*ssh.ExitError).ExitStatus()
	//		return ShellResult{
	//			Host:      server.Host,
	//			Err:       err,
	//			Cmd:       shell,
	//			Status:    constant.Fail,
	//			Code:      code,
	//			StdErrMsg: fmt.Sprintf("[%s] 执行失败, %s", server.Host, string(errOut.Bytes()))}
	//	}
	//
	//	return ShellResult{}
	//}

	var out, errOut bytes.Buffer
	session.Stdout = &out
	session.Stderr = &errOut

	if err := session.Run(shell); err != nil {
		code := err.(*ssh.ExitError).ExitStatus()
		session.Close()
		return ShellResult{
			Host:      server.Host,
			Err:       err,
			Cmd:       shell,
			Status:    constant.Fail,
			Code:      code,
			StdErrMsg: fmt.Sprintf("[%s] 执行失败, %s", server.Host, string(errOut.Bytes()))}
	}

	logger.Infof("[%s] 执行命令成功 <- end", server.Host)
	logger.Debugf("[%s] 执行结果 => %s...\n", server.Host, string(out.Bytes()))

	var subOut string

	if len(string(out.Bytes())) > 20 {
		subOut = string(out.Bytes())[:20]
	} else {
		subOut = string(out.Bytes())
	}

	defer session.Close()
	return ShellResult{Host: server.Host, StdOut: subOut, Cmd: strings.TrimPrefix(shell, "\n"), Status: constant.Success}
}

func rootMuxShell(w io.Writer, r, e io.Reader, rootPassword string) (chan<- string, <-chan string) {
	in := make(chan string, 1)
	out := make(chan string, 1)
	var wg sync.WaitGroup
	wg.Add(1) //for the shell itself
	go func() {
		for cmd := range in {
			wg.Add(1)
			w.Write([]byte(cmd + "\n"))
			wg.Wait()
		}
	}()
	go func() {
		// here i try to grep sudo from stderr, but not work
		var (
			buf [65 * 1024]byte
			t   int
		)
		for {
			n, err := e.Read(buf[t:])
			if err != nil && err.Error() != "EOF" {
				//fmt.Println(err)
			}
			if s := string(buf[t:]); strings.Contains(s, "su - root -c") {
				w.Write([]byte(fmt.Sprintf("%s\n", rootPassword)))
			} else {
			}
			t += n
		}
	}()
	go func() {
		var (
			buf [65 * 1024]byte
			t   int
		)
		for {
			n, err := r.Read(buf[t:])
			if err != nil {
				//fmt.Println(err.Error())
				close(in)
				close(out)
				return
			}
			if s := string(buf[t:]); strings.Contains(s, "Password:") {
				w.Write([]byte(fmt.Sprintf("%s\n", rootPassword)))
			} else {
			}
			t += n
			if buf[t-2] == '$' { //assuming the $PS1 == 'sh-4.3$ '
				out <- string(buf[:t])
				t = 0
				//out <- ""
				wg.Done()
			}
		}
	}()
	return in, out
}

func RunOnNodeWithChangeToRoot(shell string, server ServerInternal, timeout time.Duration, logger *logrus.Logger) (re ShellResult) {
	defer handleErr(&re.Err)

	logger.Infof("[%s] 开始执行指令 -> start", server.Host)
	logger.Debugf("\n# 指令开始\n%s\n# 指令结束\n", shell)
	session, err := server.SSHConnect(timeout)
	defer session.Close()

	if err != nil {
		errMsg := fmt.Sprintf("%s ssh会话建立失败->%s", server.Host, err.Error())
		return ShellResult{
			Host:      server.Host,
			Err:       err,
			Cmd:       shell,
			Status:    constant.Fail,
			Code:      -1,
			StdErrMsg: errMsg,
		}
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		return ShellResult{
			Host:      server.Host,
			Err:       err,
			Cmd:       shell,
			Status:    constant.Fail,
			Code:      -1,
			StdErrMsg: "",
		}
	}

	w, err := session.StdinPipe()
	if err != nil {
		panic(err)
	}
	r, err := session.StdoutPipe()
	if err != nil {
		panic(err)
	}

	e, err := session.StderrPipe()
	if err != nil {
		panic(err)
	}
	in, out := rootMuxShell(w, r, e, server.RootPassword)
	if err := session.Shell(); err != nil {
		return ShellResult{
			Host:      server.Host,
			Err:       err,
			Cmd:       shell,
			Status:    constant.Fail,
			Code:      -1,
			StdErrMsg: "",
		}
	}
	<-out //ignore the shell output

	in <- fmt.Sprintf("su - root -c \"%s\"", shell)
	outPut := <-out

	strings.Trim(outPut, "Password:\n")

	arr := strings.Split(outPut, "\n")
	var result string
	for _, v := range arr {
		if !strings.HasPrefix(v, "Password:") && !strings.HasSuffix(v, "]$ ") {
			result += v
		}
	}

	in <- "exit"
	session.Wait()

	logger.Infof("[%s] 执行命令成功 <- end", server.Host)
	logger.Debugf("[%s] 执行结果 => \n%s...\n", server.Host, result)

	return ShellResult{Host: server.Host, StdOut: result, Cmd: strings.TrimPrefix(shell, "\n"), Status: constant.Success}
}

// ReturnRunResult 获取执行结果
func (server ServerInternal) ReturnRunResult(item RemoteRunItem) ShellResult {
	item.Logger.Infof("-> %s开始执行命令...", server.Host)
	item.Logger.Debug(item.Cmd)
	session, err := server.SSHConnect(item.SSHTimeout)
	if err != nil {
		return ShellResult{Err: fmt.Errorf("%s建立ssh会话失败 -> %s", server.Host, err.Error())}
	}

	var out, errOut bytes.Buffer
	session.Stdout = &out
	session.Stderr = &errOut

	if err := session.Run(item.Cmd); err != nil {
		code := err.(*ssh.ExitError).ExitStatus()
		item.Logger.Debugf("执行命令失败, %s -> %s", err, string(errOut.Bytes()))
		return ShellResult{Code: code, Err: err, StdErrMsg: fmt.Sprintf("%s执行失败, %s", server.Host, string(errOut.Bytes()))}
	}

	item.Logger.Infof("[%s] 执行命令成功 <- end", server.Host)
	item.Logger.Debugf("[%s] 执行结果 => \n%s...\n", server.Host, string(out.Bytes()))
	defer session.Close()
	return ShellResult{StdOut: string(out.Bytes())}
}

func (server ServerInternal) SSHConnect(timeout time.Duration) (*ssh.Session, error) {
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

	clientConfig = &ssh.ClientConfig{
		User:            s.UserName,
		Auth:            auth,
		Timeout:         timeout,
		HostKeyCallback: hostKeyCallbk,
	}

	// connect to ssh
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
	// ~/.ssh/id_rsa
	keyPath, err := homedir.Expand(kPath)
	// /root/.ssh/id_rsa
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

func (server ServerInternal) completeDefault() ServerInternal {

	// ignore error
	_ = value.SetStructDefaultValue(&server, "Port", "22")
	_ = value.SetStructDefaultValue(&server, "UserName", constant.Root)
	_ = value.SetStructDefaultValue(&server, "PrivateKeyPath", constant.RsaPrvPath)

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

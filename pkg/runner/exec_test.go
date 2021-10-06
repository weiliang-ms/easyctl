package runner

import (
	"github.com/stretchr/testify/assert"
	"github.com/weiliang-ms/easyctl/pkg/util/constant"
	"os"
	"runtime"
	"testing"
)

func TestLocalRun(t *testing.T) {
	// test error cmd
	result := LocalRun("ddd", nil)
	assert.NotNil(t, result.Err)

	// test success cmd
	if runtime.GOOS == "linux" {
		result = LocalRun("ls", nil)
		assert.Nil(t, result.Err)
	}
}

func TestParallelRun(t *testing.T) {

	server := ServerInternal{
		Host:     "10.10.10.10",
		Port:     "22",
		Password: "123",
		Username: "root",
	}

	script := "ls"

	// test shell
	executor := ExecutorInternal{
		Servers: []ServerInternal{server},
		Script:  script,
	}

	re := executor.ParallelRun()
	for v := range re {
		assert.NotNil(t, v.Err)
		assert.Equal(t, script, v.Cmd)
	}

	// test shell-script file
	script = "1.sh"
	f, _ := os.Create(script)
	f.Write([]byte("pwd"))
	f.Close()
	executor = ExecutorInternal{
		Servers: []ServerInternal{server},
		Script:  script,
	}

	re = executor.ParallelRun()
	for v := range re {
		assert.NotNil(t, v.Err)
		assert.Equal(t, "pwd", v.Cmd)
	}
	os.Remove(script)
}

func TestCompleteDefault(t *testing.T) {
	server := ServerInternal{
		Password: "123",
	}
	s := server.completeDefault()
	assert.Equal(t, "22", s.Port)
	assert.Equal(t, constant.Root, s.Username)
	assert.Equal(t, constant.RsaPrvPath, s.PrivateKeyPath)
}

//// ParallelRun 并发执行
//func (executor ExecutorInternal) ParallelRun() chan ShellResult {
//
//	executor.Logger.Infoln("开始并行执行命令...")
//	wg := sync.WaitGroup{}
//	ch := make(chan ShellResult, len(executor.Servers))
//
//	// todo: 屏蔽rm -rf
//	// 判断入参为文件还是shell
//
//	if _, err := os.Stat(executor.Script); err == nil {
//		b, _ := os.ReadFile(executor.Script)
//		executor.Script = string(b)
//	}
//
//	for _, v := range executor.Servers {
//		wg.Add(1)
//		go func(s ServerInternal) {
//			ch <- executor.runOnNode(s)
//			defer wg.Done()
//		}(v)
//	}
//
//	wg.Wait()
//	close(ch)
//	return ch
//}
//
//func (executor ExecutorInternal) runOnNode(server ServerInternal) (re ShellResult) {
//	var shell string
//	defer handleErr(&re.Err)
//
//	if executor.Logger.Level == logrus.DebugLevel {
//		shell = executor.Script
//	} else {
//		shell = "shell content"
//	}
//
//	// 截取cmd output 长度
//	var subCmd string
//	if len(executor.Script) > 10 {
//		subCmd = "******"
//	} else {
//		subCmd = executor.Script
//	}
//
//	executor.Logger.Infof("[%s] 开始执行指令 -> %s", server.Host, shell)
//	session, err := server.sshConnect()
//	defer session.Close()
//
//	if err != nil {
//		errMsg := fmt.Sprintf("ssh会话建立失败->%s", err.(*net.OpError).Error())
//		return ShellResult{
//			Host:      server.Host,
//			Err:       err,
//			Cmd:       executor.Script,
//			Status:    constant.Fail,
//			Code:      -1,
//			StdErrMsg: errMsg,
//		}
//	}
//
//	if executor.OutPutRealTime == true {
//		session.Stdout = os.Stdout
//		err := session.Run(executor.Script)
//		if err != nil {
//			return ShellResult{
//				Host:      server.Host,
//				Err:       errors.New(err.(*ssh.ExitError).String()),
//				Cmd:       executor.Script,
//				Status:    constant.Fail,
//				Code:      err.(*ssh.ExitError).ExitStatus(),
//				StdErrMsg: err.(*ssh.ExitError).String()}
//		}
//
//	} else {
//		out, err := session.Output(executor.Script)
//		if err != nil {
//			return ShellResult{
//				Host:      server.Host,
//				Err:       errors.New(err.(*ssh.ExitError).String()),
//				Cmd:       executor.Script,
//				Status:    constant.Fail,
//				Code:      err.(*ssh.ExitError).ExitStatus(),
//				StdErrMsg: err.(*ssh.ExitError).String()}
//		}
//
//		executor.Logger.Infof("<- %s执行命令成功...", server.Host)
//		executor.Logger.Debugf("%s -> 返回值: %s", server.Host, string(out))
//
//		var subOut string
//
//		if len(string(out)) > 20 {
//			subOut = string(out)[:20]
//		} else {
//			subOut = string(out)
//		}
//
//		return ShellResult{Host: server.Host, StdOut: subOut,
//			Cmd: strings.TrimPrefix(subCmd, "\n"), Status: constant.Success}
//	}
//
//	return ShellResult{}
//}
//
//// ReturnRunResult 获取执行结果
//func (server ServerInternal) ReturnRunResult(cmd string) ShellResult {
//	log.Printf("<- %s开始执行命令...\n", server.Host)
//	session, err := server.sshConnect()
//	if err != nil {
//		return ShellResult{Err: fmt.Errorf("%s建立ssh会话失败 -> %s", server.Host, err.Error())}
//	}
//
//	combo, err := session.CombinedOutput(cmd)
//	if err != nil {
//		//klog.Fatal("远程执行cmd 失败",err)
//		return ShellResult{Err: fmt.Errorf("%s执行失败, %s", server.Host, combo)}
//	}
//	log.Printf("<- %s执行命令成功，返回结果 => %s...\n", server.Host, string(combo))
//	defer session.Close()
//
//	return ShellResult{StdOut: string(combo)}
//}
//
//func (server ServerInternal) sshConnect() (*ssh.Session, error) {
//	s := server.completeDefault()
//	var (
//		auth         []ssh.AuthMethod
//		addr         string
//		clientConfig *ssh.ClientConfig
//		client       *ssh.Client
//		session      *ssh.Session
//		err          error
//	)
//	// get auth method
//	auth = make([]ssh.AuthMethod, 0)
//
//	if server.Password == "" || server.PrivateKeyPath != "" {
//		auth = append(auth, publicKeyAuthFunc(server.PrivateKeyPath))
//	} else {
//		auth = append(auth, ssh.Password(s.Password))
//	}
//
//	hostKeyCallbk := func(hostname string, remote net.Addr, key ssh.PublicKey) error {
//		return nil
//	}
//
//	clientConfig = &ssh.ClientConfig{
//		User:            s.Username,
//		Auth:            auth,
//		Timeout:         3 * time.Second,
//		HostKeyCallback: hostKeyCallbk,
//	}
//
//	// connet to ssh
//	addr = fmt.Sprintf("%s:%s", server.Host, server.Port)
//
//	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
//		return nil, err
//	}
//
//	// create session
//	if session, err = client.NewSession(); err != nil {
//		return nil, err
//	}
//
//	return session, nil
//}
//
//func publicKeyAuthFunc(kPath string) ssh.AuthMethod {
//	keyPath, err := homedir.Expand(kPath)
//	if err != nil {
//		klog.Fatal("find key's home dir failed", err)
//	}
//	key, err := ioutil.ReadFile(keyPath)
//	if err != nil {
//		klog.Fatal("ssh key file read failed", err)
//	}
//	// Create the Signer for this private key.
//	signer, err := ssh.ParsePrivateKey(key)
//	if err != nil {
//		klog.Fatal("ssh key signer failed", err)
//	}
//	return ssh.PublicKeys(signer)
//}

func TestPublicKeyAuthFunc(t *testing.T) {
	_, err := publicKeyAuthFunc("/ss/ss")
	assert.Errorf(t, err, "ssh key file read failed open /ss/ss: The system cannot find the path specified.")

	// test not read
	_, err = publicKeyAuthFunc("~/.ssh/1.pub")
	assert.NotNil(t, err)

}

//
//func handleErr(err *error) {
//
//	if v := recover(); v != nil {
//		if e, ok := v.(runtime.Error); ok {
//			*err = e
//		} else {
//			panic(v)
//		}
//	}
//}

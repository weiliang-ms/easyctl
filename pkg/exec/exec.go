package exec

import (
	_ "embed"
	"errors"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/weiliang-ms/easyctl/pkg/util"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"k8s.io/klog"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

//go:embed executor.yaml
var config []byte

func Run(configFile string, level util.LogLevel) error {
	if configFile == "" {
		klog.Infof("检测到配置文件为空，生成配置文件样例 -> %s", util.ConfigFile)
		os.WriteFile(util.ConfigFile, config, 0666)
	}

	b, err := os.ReadFile(configFile)
	if err != nil {
		return errors.New("读取配置文件失败")
	}
	exec, err := parseConfig(b)
	if err != nil {
		return err
	}

	if err := exec.Run(true, level); err != nil {
		return err
	}

	return nil
}

func (exec ExecutorItem) Run(parallel bool, level util.LogLevel) error {

	if !parallel {
		if err := exec.run(level); err != nil {
			return err
		}
	} else {
		if err := exec.parallelRun(level); err != nil {
			return err
		}
	}

	return nil
}

func (exec ExecutorItem) run(level util.LogLevel) error {

	for _, v := range exec.Server {
		if err := runOnNode(v, exec.Script, level); err != nil {
			return err
		}
	}
	return nil
}

func (exec ExecutorItem) parallelRun(level util.LogLevel) error {

	klog.Infoln("开始并行执行命令...")
	wg := sync.WaitGroup{}
	wg.Add(len(exec.Server))
	ch := make(chan error, len(exec.Server))

	var script string

	if _, err := os.Stat(exec.Script); err != nil {
		script = exec.Script
	} else {
		b, _ := os.ReadFile(exec.Script)
		script = string(b)
	}

	for _, v := range exec.Server {
		go func(s Server) {
			err := runOnNode(s, script, level)
			ch <- err
			//close(ch)
			defer wg.Done()
		}(v)
	}

	wg.Wait()

	for v := range ch {
		if v != nil {
			return v
		}
	}

	return nil
}

func runOnNode(s Server, cmd string, level util.LogLevel) error {
	//session , err := session(s)
	var shell string

	if level == util.Debug {
		shell = cmd
	}

	klog.Infof("[%s] 开始执行指令 -> %s\n", s.Host, shell)
	session, err := s.sshConnect()
	if err != nil {
		return err
	}

	combo, err := session.CombinedOutput(cmd)
	if err != nil {
		//klog.Fatal("远程执行cmd 失败",err)
		return errors.New(fmt.Sprintf("%s执行失败, %s", s.Host, combo))
	}
	log.Printf("<- %s执行命令成功...\n", s.Host)
	if string(combo) != "" && level == util.Debug {
		fmt.Printf("<- [%s] 命令输出: ->\n\n%s\n", s.Host, string(combo))
	}

	defer session.Close()

	return err
}

func (s Server) sshConnect() (*ssh.Session, error) {
	server := s.completeDefault()
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
	auth = append(auth, ssh.Password(server.Password))

	hostKeyCallbk := func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		return nil
	}

	clientConfig = &ssh.ClientConfig{
		User: server.Username,
		Auth: auth,
		// Timeout:             30 * time.Second,
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

func session(server Server) (*ssh.Session, error) {

	server = server.completeDefault()

	//创建sshp登陆配置
	config := &ssh.ClientConfig{
		Timeout:         time.Second, //ssh 连接time out 时间一秒钟, 如果ssh验证错误 会在一秒内返回
		User:            server.Username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //这个可以， 但是不够安全
		//HostKeyCallback: hostKeyCallBackFunc(h.Host),
	}
	//if sshType == "password" {
	config.Auth = []ssh.AuthMethod{ssh.Password(server.Password)}
	//} else {
	//	config.Auth = []ssh.AuthMethod{publicKeyAuthFunc(sshKeyPath)}
	//}

	//dial 获取ssh client
	addr := fmt.Sprintf("%s:%s", server.Host, server.Port)
	sshClient, err := ssh.Dial("tcp", addr, config)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s创建ssh client 失败, %s", server.Host, err.Error()))
	}
	defer sshClient.Close()

	//创建ssh-session
	session, err := sshClient.NewSession()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s创建ssh session 失败, %s", server.Host, err.Error()))
	}

	return session, nil

}

func publicKeyAuthFunc(kPath string) ssh.AuthMethod {
	keyPath, err := homedir.Expand(kPath)
	if err != nil {
		klog.Fatal("find key's home dir failed", err)
	}
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		klog.Fatal("ssh key file read failed", err)
	}
	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		klog.Fatal("ssh key signer failed", err)
	}
	return ssh.PublicKeys(signer)
}

func (s Server) completeDefault() Server {
	if s.Port == "" {
		s.Port = "22"
	}

	if s.Username == "" {
		s.Username = "root"
	}

	if s.PublicKeyPath == "" {
		s.PublicKeyPath = "~/.ssh/id_rsa.pub"
	}

	return s
}

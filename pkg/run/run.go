package run

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/sftp"
	"github.com/vbauerster/mpb/v6"
	"github.com/vbauerster/mpb/v6/decor"
	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v2"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"time"
)

type ExecResult struct {
	ExitCode int
	StdErr   string
	StdOut   string
}

type ServerList struct {
	Server []Server `yaml:"server,flow"`
}
type Server struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type KeepaliveServerList struct {
	Vip       string   `yaml:"vip"`
	Interface string   `yaml:"interface"`
	Server    []Server `yaml:"server,flow"`
}

type ShellResult struct {
	Host   string
	Cmd    string
	Code   int
	Status string
}

func ParseServerList(yamlPath string) ServerList {

	var serverList ServerList
	if f, err := os.Open(yamlPath); err != nil {
		fmt.Println("open yaml...")
		log.Fatal(err)
	} else {
		decodeErr := yaml.NewDecoder(f).Decode(&serverList)
		if decodeErr != nil {
			fmt.Println("decode failed...")
			log.Fatal(decodeErr)
		}
	}
	_, err := json.Marshal(serverList)
	if err != nil {
		fmt.Println("marshal failed...")
		log.Fatal(err)
	}

	return serverList
}

func ParseKeepaliveList(yamlPath string) KeepaliveServerList {

	var serverList KeepaliveServerList
	if f, err := os.Open(yamlPath); err != nil {
		log.Println("open yaml...")
		log.Fatal(err)
	} else {
		decodeErr := yaml.NewDecoder(f).Decode(&serverList)
		if decodeErr != nil {
			log.Println("decode failed...")
			log.Fatal(decodeErr)
		}
	}

	data, err := json.Marshal(serverList)

	fmt.Println(string(data))
	if err != nil {
		log.Println("marshal failed...")
		log.Fatal(err)
	}

	//log.Printf("%v",serverList)
	return serverList
}

// 远程写文件
func RemoteWriteFile(srcPath string, dstPath string, instance Server) {
	// init sftp
	sftp, err := sftpConnect(instance.Username, instance.Password, instance.Host, instance.Port)
	if err != nil {
		fmt.Println(err.Error())
	}

	log.Printf("-> transfer %s to %s", srcPath, instance.Host)
	dstFile, err := sftp.Create(dstPath)
	sftp.Chmod(dstPath, 0755)

	if err != nil {
		fmt.Println(err.Error())
	}

	f, _ := os.Open(srcPath)
	ff, _ := os.Stat(srcPath)

	total := ff.Size()
	reader := io.LimitReader(f, total)

	p := mpb.New(
		mpb.WithWidth(60),
		mpb.WithRefreshRate(180*time.Millisecond),
	)

	bar := p.Add(total,
		mpb.NewBarFiller("[=>-|"),
		mpb.PrependDecorators(
			decor.CountersKibiByte("% .2f / % .2f"),
		),
		mpb.AppendDecorators(
			decor.EwmaETA(decor.ET_STYLE_GO, 90),
			decor.Name(" ] "),
			decor.EwmaSpeed(decor.UnitKiB, "% .2f", 60),
		),
	)

	// create proxy reader
	proxyReader := bar.ProxyReader(reader)
	io.Copy(dstFile, proxyReader)

	p.Wait()

	defer f.Close()
	defer proxyReader.Close()
	defer dstFile.Close()

}

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

func (server Server) RemoteShellReturnStd(cmd string) string {
	session, conErr := server.sshConnect()
	if conErr != nil {
		log.Fatal(conErr)
	}

	defer session.Close()

	combo, runErr := session.CombinedOutput(cmd)

	if runErr != nil {
		log.Fatal(runErr.Error())
	}
	return string(combo)
}

func (server Server) RemoteShell(cmd string) ShellResult {

	var resulft ShellResult
	resulft.Cmd = cmd
	resulft.Host = server.Host

	log.Printf("-> [%s] shell => %s", server.Host, cmd)
	session, conErr := server.sshConnect()
	if conErr != nil {
		log.Fatal(conErr)
	}

	defer session.Close()

	combo, runErr := session.CombinedOutput(cmd)

	if runErr != nil {
		e, _ := runErr.(*ssh.ExitError)
		resulft.Code = e.ExitStatus()
		log.Print(runErr.Error())
	}

	log.Printf("<- call back %s\n %s", server.Host, string(combo))

	if resulft.Code == 0 {
		resulft.Status = "success"
	} else {
		resulft.Status = "failed"
	}

	return resulft
}

func (server *Server) sshConnect() (*ssh.Session, error) {
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

func Shell(command string) (re ExecResult) {
	//函数返回一个*Cmd，用于使用给出的参数执行name指定的程序
	cmd := exec.Command("/bin/bash", "-c", command)
	//log.Printf("%s 执行语句：%s\n", PrintCyan(constant.Shell), command)
	//读取io.Writer类型的cmd.Stdout，再通过bytes.Buffer(缓冲byte类型的缓冲器)将byte类型转化为string类型(out.String():这是bytes类型提供的接口)

	stderr, _ := cmd.StderrPipe()
	stdout, _ := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		log.Fatal(err.Error())
	}

	// 标准输出
	logScan := bufio.NewScanner(stdout)
	go func() {
		for logScan.Scan() {
			log.Println(logScan.Text())
			re.StdOut = logScan.Text()
		}
	}()

	// 标准错误输出
	errBuf := bytes.NewBufferString("")
	scan := bufio.NewScanner(stderr)
	for scan.Scan() {
		s := scan.Text()
		log.Println("build error: ", s)
		errBuf.WriteString(s)
		errBuf.WriteString("\n")
		re.StdErr = logScan.Text()
	}

	// 等待命令执行完
	cmd.Wait()

	re.ExitCode = cmd.ProcessState.ExitCode()

	return re

}

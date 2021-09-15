package runner

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pkg/sftp"
	"github.com/vbauerster/mpb/v6"
	"github.com/vbauerster/mpb/v6/decor"
	"github.com/weiliang-ms/easyctl/pkg/util"
	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type ServerList struct {
	Common    CommonServerList
	Harbor    HarborServerList
	HA        HaProxyServerList
	Keepalive KeepaliveServerList
	Docker    DockerServerList
	Compose   DockerComposeServerList
}

type ExecResult struct {
	ExitCode int
	StdErr   string
	StdOut   string
}

type CommonServerList struct {
	Server  []Server `yaml:"server,flow"`
	Exclude []string `yaml:"exclude,flow"`
}

type Server struct {
	Host           string `yaml:"host"`
	Port           string `yaml:"port"`
	Username       string `yaml:"username"`
	Password       string `yaml:"password"`
	PrivateKeyPath string `yaml:"privateKeyPath,omitempty" json:"privateKeyPath,omitempty"`
}

type KeepaliveServerList struct {
	Attribute Keepalive `yaml:"keepalive"`
}

type Keepalive struct {
	Vip       string   `yaml:"vip"`
	Interface string   `yaml:"interface"`
	Server    []Server `yaml:"server,flow"`
}

type HaProxyServerList struct {
	Attribute HaProxy `yaml:"haproxy"`
}

type DockerServerList struct {
	Attribute Docker `yaml:"docker"`
}

type Docker struct {
	Servers []Server `yaml:"server,flow"`
}

type DockerComposeServerList struct {
	Attribute DockerCompose `yaml:"docker-compose"`
}

type DockerCompose struct {
	Server []Server `yaml:"server,flow"`
}

type HaProxy struct {
	Server      []Server  `yaml:"server,flow"`
	BalanceList []Balance `yaml:"balance,flow"`
}

type Balance struct {
	Name    string   `yaml:"name"`
	Port    int      `yaml:"port"`
	Address []string `yaml:"address"`
}

type HarborServerList struct {
	Attribute Harbor `yaml:"harbor"`
}

type Harbor struct {
	Server        []Server      `yaml:"server,flow"`
	Project       HarborProject `yaml:"project"`
	Password      string        `yaml:"password"`
	DataDir       string        `yaml:"dataDir"`
	Domain        string        `yaml:"domain"`
	HttpPort      string        `yaml:"http-port"`
	ResolvAddress string        `yaml:"resolve-ip"`
	OpenGC        bool          `yaml:"openGC"`
	ReserveNum    int8          `yaml:"reserveNum"`
}

type HarborProject struct {
	Private []string `yaml:"private"`
	Public  []string `yaml:"public"`
}

type ShellResult struct {
	Host      string
	Cmd       string
	Code      int
	Status    string
	Stderr    error
	StdOut    string // 标准输出
	StdErrMsg string
	Output    string // 所有输出
}

type Installer struct {
	Offline         bool
	Cmd             string
	FileName        string
	ServerListPath  string
	OfflineFilePath string
	DataDir         string
}

func ParseCommonServerList(yamlPath string) (c CommonServerList) {
	var decodeErr, marshalErr error

	f, err := os.Open(yamlPath)
	if err != nil {
		log.Fatal(err)
	}

	decodeErr = yaml.NewDecoder(f).Decode(&c)
	_, marshalErr = json.Marshal(&c)

	if decodeErr != nil {
		log.Fatal(decodeErr)
	}

	if marshalErr != nil {
		log.Fatal(err)
	}

	var serverList []Server

	for _, v := range c.Server {
		// 192.168.235.[1:3]
		if strings.Contains(v.Host, "[") {
			log.Println("检测到配置文件中含有IP段，开始解析组装...")
			//192.168.235.
			baseAddress := strings.Split(v.Host, "[")[0]
			log.Printf("解析到IP子网网段为：%s...\n", baseAddress)

			// 1:3] -> 1:3
			ipRange := strings.Split(strings.Split(v.Host, "[")[1], "]")[0]
			log.Printf("解析到IP区间为：%s...\n", ipRange)

			// 1:3 -> 1
			begin := strings.Split(ipRange, ":")[0]
			log.Printf("解析到起始IP为：%s...\n", fmt.Sprintf("%s%s", baseAddress, begin))

			// 1:3 -> 3
			end := strings.Split(ipRange, ":")[1]
			log.Printf("解析到末尾IP为：%s...\n", fmt.Sprintf("%s%s", baseAddress, end))

			// string -> int
			beginIndex, _ := strconv.Atoi(begin)
			endIndex, _ := strconv.Atoi(end)

			for i := beginIndex; i <= endIndex; i++ {
				server := Server{
					Host:           fmt.Sprintf("%s%d", baseAddress, i),
					Port:           v.Port,
					Username:       v.Username,
					Password:       v.Password,
					PrivateKeyPath: v.PrivateKeyPath,
				}

				if !util.SliceContain(c.Exclude, server.Host) {
					//log.Printf("add host: %s\n", server.Host)
					serverList = append(serverList, server)
				}

			}
		} else {
			serverList = append(serverList, v)
		}
	}

	c.Server = serverList

	return c
}

func ParseServerList(yamlPath string, v interface{}) (list ServerList) {

	var decodeErr, marshalErr error

	f, err := os.Open(yamlPath)
	if err != nil {
		log.Fatal(err)
	}

	// todo:优化反射方式
	switch reflect.ValueOf(v).Type().String() {
	case "runner.DockerServerList":
		decodeErr = yaml.NewDecoder(f).Decode(&list.Docker)
		_, marshalErr = json.Marshal(&list.Docker)
	case "runner.HarborServerList":
		decodeErr = yaml.NewDecoder(f).Decode(&list.Harbor)
		_, marshalErr = json.Marshal(&list.Harbor)
	case "runner.KeepaliveServerList":
		decodeErr = yaml.NewDecoder(f).Decode(&list.Keepalive)
		_, marshalErr = json.Marshal(&list.Keepalive)
	default:
		decodeErr = yaml.NewDecoder(f).Decode(&list.Common)
		_, marshalErr = json.Marshal(&list.Common)
	}

	if decodeErr != nil {
		log.Fatal(decodeErr)
	}

	if marshalErr != nil {
		log.Fatal(err)
	}

	return list
}

// 远程写文件
func ScpFile(srcPath string, dstPath string, server Server, mode os.FileMode) {
	// init sftp
	sftp, err := sftpConnect(server.Username, server.Password, server.Host, server.Port)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("-> transfer %s to %s:%s", srcPath, server.Host, dstPath)
	dstFile, err := sftp.Create(dstPath)
	_ = sftp.Chmod(dstPath, mode)

	if err != nil {
		log.Fatal(err.Error())
	}

	f, _ := os.Open(srcPath)
	ff, _ := os.Stat(srcPath)
	total := ff.Size()
	reader := io.LimitReader(f, total)

	// 初始化进度条
	p := mpb.New(
		mpb.WithWidth(60),                  // 进度条长度
		mpb.WithRefreshRate(1*time.Second), // 刷新速度
	)

	//
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

	log.Println("传输成功...")
	defer f.Close()
	defer proxyReader.Close()
	defer dstFile.Close()

}

// Scp 远程写文件
func (server Server) Scp(srcPath string, dstPath string, mode os.FileMode) error {
	// init sftp
	sftp, connetErr := sftpConnect(server.Username, server.Password, server.Host, server.Port)
	if connetErr != nil {
		return errors.New(fmt.Sprintf("连接远程主机：%s失败 ->%s",
			server.Host, connetErr.Error()))
	}

	defer sftp.Close()

	log.Printf("-> transfer %s to %s:%s\n", srcPath, server.Host, dstPath)
	dstFile, createErr := sftp.Create(dstPath)
	if createErr != nil {
		return errors.New(fmt.Sprintf("创建远程主机：%s文件句柄: %s失败 ->%s",
			server.Host, dstPath, createErr.Error()))
	}
	defer dstFile.Close()

	modErr := sftp.Chmod(dstPath, mode)
	if modErr != nil {
		return errors.New(fmt.Sprintf("修改%s:%s文件句柄权限失败 ->%s",
			server.Host, dstPath, createErr.Error()))
	}

	// 获取文件大小信息
	f, _ := os.Open(srcPath)
	defer f.Close()
	ff, _ := os.Stat(srcPath)
	reader := io.LimitReader(f, ff.Size())

	// 初始化进度条
	p := mpb.New(
		mpb.WithWidth(60),                  // 进度条长度
		mpb.WithRefreshRate(1*time.Second), // 刷新速度
	)

	//
	bar := p.Add(ff.Size(),
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
	_, ioErr := io.Copy(dstFile, proxyReader)
	if ioErr != nil {
		return errors.New(fmt.Sprintf("传输%s:%s失败 ->%s",
			server.Host, dstPath, createErr.Error()))
	}

	p.Wait()

	defer proxyReader.Close()

	return nil
}

func (server Server) CommandExists(cmd string) bool {

	session, conErr := server.sshConnect()
	if conErr != nil {
		log.Fatal(conErr)
	}

	defer session.Close()

	_, runErr := session.CombinedOutput(fmt.Sprintf("command -v %s > /dev/null 2>&1", cmd))

	if runErr != nil {
		return false
	}

	return true
}

func LocalCommandExists(cmd string) bool {
	re := Shell(fmt.Sprintf("command -v %s > /dev/null 2>&1", cmd))
	if re.ExitCode != 0 {
		return false
	}
	return true
}

func RemoteWriteFile(b []byte, dstPath string, instance Server, mode os.FileMode) {
	// init sftp
	sftp, err := sftpConnect(instance.Username, instance.Password, instance.Host, instance.Port)
	if err != nil {
		fmt.Println(err.Error())
	}

	dstFile, err := sftp.Create(dstPath)
	sftp.Chmod(dstPath, mode)

	if err != nil {
		fmt.Println(err.Error())
	}

	dstFile.Write(b)

	defer dstFile.Close()

}

func (server Server) WriteRemoteFile(b []byte, dstPath string, mode os.FileMode) {
	// init sftp
	sftp, err := sftpConnect(server.Username, server.Password, server.Host, server.Port)
	if err != nil {
		fmt.Println(err.Error())
	}

	dstFile, err := sftp.Create(dstPath)
	sftp.Chmod(dstPath, mode)

	if err != nil {
		fmt.Println(err.Error())
	}

	_, _ = dstFile.Write(b)

	defer dstFile.Close()

}

// 移动目录下文件至新目录
func (server Server) MoveDirFiles(srcDir string, dstDir string) {
	files, _ := ioutil.ReadDir(srcDir)
	for _, f := range files {
		if !f.IsDir() {
			oldpath := srcDir + f.Name()
			newPath := dstDir + "/" + f.Name()
			log.Printf("%s => %s", oldpath, newPath)
			err := os.Rename(oldpath, newPath)
			if err != nil {
				log.Fatal(err.Error())
			}
		}
	}
}

// 清理目录
func (server Server) DelDirectory(dir string) error {

	log.Printf("删除目录: %s:%s", server.Host, dir)
	sftp, err := sftpConnect(server.Username, server.Password, server.Host, server.Port)
	defer sftp.Close()
	if err != nil {
		log.Println(err.Error())
	}

	f, err := sftp.Stat(dir)
	if err != nil {
		return err
	}

	if !f.IsDir() {
		log.Printf("%s不是目录...", dir)
	}
	server.RemoteShell(fmt.Sprintf("rm -rf %s", dir))
	return nil

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

func (server Server) InstallSoft(installScript string) error {
	session, conErr := server.sshConnect()
	if conErr != nil {
		log.Fatal(conErr)
	}

	defer session.Close()

	_, runErr := session.CombinedOutput(installScript)

	if runErr != nil {
		return runErr
	}

	return nil
}

func (server Server) RemoteShell(cmd string) ShellResult {

	var result ShellResult
	if len(cmd) < 60 {
		result.Cmd = cmd
	} else {
		result.Cmd = "built-in shell"
	}
	result.Host = server.Host

	log.Printf("-> [%s] => exec shell...", server.Host)
	session, conErr := server.sshConnect()
	if conErr != nil {
		log.Fatal(conErr)
	}

	defer session.Close()

	combo, runErr := session.CombinedOutput(cmd)

	if runErr != nil {
		e, _ := runErr.(*ssh.ExitError)
		result.Code = e.ExitStatus()
		result.Stderr = runErr
		result.StdErrMsg = runErr.Error()
		log.Print(runErr.Error())
	}

	if result.Code == 0 {
		log.Printf("<- call back [%s] exec shell successful...", server.Host)
		result.Status = "success"
		result.StdOut = string(combo)
	} else {
		log.Printf("<- call back [%s]\n %s", server.Host, string(combo))
		result.Status = "failed"
	}

	return result
}

func (server Server) RemoteShellOutput(cmd string) {

	var result ShellResult
	if len(cmd) < 60 {
		result.Cmd = cmd
	} else {
		result.Cmd = "built-in shell"
	}
	result.Host = server.Host

	//log.Printf("-> [%s] shell => %s", server.Server, cmd)
	log.Printf("-> [%s] => exec shell...", server.Host)
	session, conErr := server.sshConnect()
	if conErr != nil {
		log.Fatal(conErr)
	}

	defer session.Close()

	// combo, runErr := session.CombinedOutput(cmd)
	session.Stderr = os.Stderr
	session.Stdout = os.Stdout
	session.Run(cmd)
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

// 立即返回，执行状态及结果
func ShortShell(command string) (re ExecResult) {

	//函数返回一个*Cmd，用于使用给出的参数执行name指定的程序
	cmd := exec.Command("/bin/bash", "-c", command)

	var stdOut, stdErr bytes.Buffer
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	err := cmd.Run()
	if err != nil {
		log.Printf("执行：%s失败 -> %s...\n", cmd, err.Error())
	}

	re.StdOut = stdOut.String()
	re.StdErr = stdErr.String()

	re.ExitCode = cmd.ProcessState.ExitCode()

	return re
}

func Shell(command string) (re ExecResult) {
	//函数返回一个*Cmd，用于使用给出的参数执行name指定的程序
	cmd := exec.Command("/bin/bash", "-c", command)

	stderr, _ := cmd.StderrPipe()
	stdout, _ := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		re.StdErr = err.Error()
		return re
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
		errBuf.WriteString(s)
		errBuf.WriteString("\n")
		re.StdErr = s
	}

	// 等待命令执行完
	cmd.Wait()
	re.ExitCode = cmd.ProcessState.ExitCode()
	return re
}

// /root or /home/username
func HomeDir(server Server) string {
	switch server.Username {
	case "root":
		return "/root"
	default:
		return fmt.Sprintf("/home/%s", server.Username)
	}
}

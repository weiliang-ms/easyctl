package util

import (
	"bufio"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type SSHInstance struct {
	Host     string
	Port     string
	Username string
	Password string
}

func (instance SSHInstance) ExecuteOriginCmd(cmd string) (msg string, exitCode int) {
	session, conErr := instance.sshConnect()
	fmt.Printf("[shell] 远程执行: %s\n", cmd)
	if conErr != nil {
		log.Fatal(conErr)
	}

	defer session.Close()

	combo, runErr := session.CombinedOutput(cmd)

	if runErr != nil {
		e, _ := runErr.(*ssh.ExitError)
		exitCode = e.ExitStatus()
	}
	return string(combo), exitCode
}
func (instance SSHInstance) ExecuteOriginCmdParallel(cmd string, wg *sync.WaitGroup) (msg string, exitCode int) {
	defer wg.Done()
	session, conErr := instance.sshConnect()
	fmt.Printf("[origin] 远程执行: %s\n", cmd)
	if conErr != nil {
		log.Fatal(conErr)
	}

	defer session.Close()

	combo, runErr := session.CombinedOutput(cmd)

	if runErr != nil {
		e, _ := runErr.(*ssh.ExitError)
		exitCode = e.ExitStatus()
	}

	return string(combo), exitCode
}
func ReadSSHInfoFromFile(hostsFile string) (instances []SSHInstance) {
	f, err := os.OpenFile(hostsFile, os.O_RDONLY, 0644)
	defer f.Close()

	if err != nil {
		log.Fatal(err.Error())
	}

	rd := bufio.NewReader(f)

	for {
		line, err := rd.ReadString('\n') //以'\n'为结束符读入一行
		if strings.Contains(line, "host") {
			instances = append(instances, setSSHObjectValue(line))
		}
		if err != nil || io.EOF == err {
			break
		}
	}
	log.Printf("%+v", instances)
	return instances
}

func setSSHObjectValue(hosts string) (instance SSHInstance) {
	//fmt.Println("-----line----" + hosts)
	cutCharacters := []string{"\n"}
	for _, v := range strings.Split(hosts, " ") {
		//fmt.Printf("---%v", v)
		if strings.Contains(v, "host") {
			instance.Host = CutCharacter(strings.TrimPrefix(v, "host="), cutCharacters)
			//fmt.Printf("host赋值%s",object.Host)
			//fmt.Printf("%+v",object)
		}

		if strings.Contains(v, "port") {
			//fmt.Println("port赋值")
			instance.Port = CutCharacter(strings.TrimPrefix(v, "port="), cutCharacters)
		}

		if strings.Contains(v, "user") {
			instance.Username = CutCharacter(strings.TrimPrefix(v, "user="), cutCharacters)
			//fmt.Printf("username赋值：%s",object.Username)
		}
		if strings.Contains(v, "password") {
			instance.Password = CutCharacter(strings.Trim(v, "password="), cutCharacters)
			//fmt.Printf("password赋值%s",instance.Password)
		}
	}
	return instance
}

func (instance SSHInstance) sshConnect() (*ssh.Session, error) {
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
	auth = append(auth, ssh.Password(instance.Password))

	hostKeyCallbk := func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		return nil
	}

	clientConfig = &ssh.ClientConfig{
		User: instance.Username,
		Auth: auth,
		// Timeout:             30 * time.Second,
		HostKeyCallback: hostKeyCallbk,
	}

	// connet to ssh
	addr = fmt.Sprintf("%s:%s", instance.Host, instance.Port)

	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	// create session
	if session, err = client.NewSession(); err != nil {
		return nil, err
	}

	return session, nil
}

func RemoteHostYumDetection(instance SSHInstance) bool {
	_, code := instance.ExecuteOriginCmd("yum install -y gcc")
	if code == 0 {
		return true
	}

	return false
}

func ScpFile(originPath string, destinationPath string, instance SSHInstance) {
	// 跨服务器拷贝

}

// 远程写文件
func OriginWriteFile(filePath string, b []byte, instance SSHInstance) {
	// init sftp
	sftp, err := SftpConnect(instance.Username, instance.Password, instance.Host, instance.Port)
	if err != nil {
		log.Fatal(err)
	}
	dstFile, err := sftp.Create(filePath)
	defer dstFile.Close()
	dstFile.Write(b)
}

func SftpConnect(user, password, host string, port string) (sftpClient *sftp.Client, err error) { //参数: 远程服务器用户名, 密码, ip, 端口
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

package runner

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"
)

//func ParseServerList(yamlPath string, v interface{}) (list ServerList) {
//
//	var decodeErr, marshalErr error
//
//	f, err := os.Open(yamlPath)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// todo:优化反射方式
//	switch reflect.ValueOf(v).Type().String() {
//	case "runner.DockerServerList":
//		decodeErr = yaml.NewDecoder(f).Decode(&list.Docker)
//		_, marshalErr = json.Marshal(&list.Docker)
//	case "runner.HarborServerList":
//		decodeErr = yaml.NewDecoder(f).Decode(&list.Harbor)
//		_, marshalErr = json.Marshal(&list.Harbor)
//	case "runner.KeepaliveServerList":
//		decodeErr = yaml.NewDecoder(f).Decode(&list.Keepalive)
//		_, marshalErr = json.Marshal(&list.Keepalive)
//	default:
//		decodeErr = yaml.NewDecoder(f).Decode(&list.Common)
//		_, marshalErr = json.Marshal(&list.Common)
//	}
//
//	if decodeErr != nil {
//		log.Fatal(decodeErr)
//	}
//
//	if marshalErr != nil {
//		log.Fatal(err)
//	}
//
//	return list
//}

// CommandExists 判断命令知否存在
func (server ServerInternal) CommandExists(cmd string) bool {

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

// WriteRemoteFile 远程写文件
func (server ServerInternal) WriteRemoteFile(b []byte, dstPath string, mode os.FileMode) {
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

// MoveDirFiles 移动目录下文件至新目录
func (server ServerInternal) MoveDirFiles(srcDir string, dstDir string) {
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

// RemoteShellReturnStd 单机远程执行返回结果
func (server ServerInternal) RemoteShellReturnStd(cmd string) string {
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

// InstallSoft 安装软件
func (server ServerInternal) InstallSoft(installScript string) error {
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

// HomeDir /root or /home/username
func (server ServerInternal) HomeDir() string {
	switch server.Username {
	case "root":
		return "/root"
	default:
		return fmt.Sprintf("/home/%s", server.Username)
	}
}

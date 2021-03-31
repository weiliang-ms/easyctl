package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v2"
	"log"
	"net"
	"os"
	"time"
)

type ServerList struct {
	Server []Server `yaml:"server,flow"`
}
type Server struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func parseServerList(yamlPath string) ServerList {

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

// 远程写文件
func remoteWriteFile(filePath string, b []byte, instance Server) {
	// init sftp
	sftp, err := sftpConnect(instance.Username, instance.Password, instance.Host, instance.Port)
	if err != nil {
		fmt.Println(err.Error())
	}
	dstFile, err := sftp.Create(filePath)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer dstFile.Close()
	dstFile.Write(b)
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

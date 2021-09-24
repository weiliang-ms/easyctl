package runner

import (
	"fmt"
	"github.com/pkg/sftp"
	"github.com/vbauerster/mpb/v6"
	"github.com/vbauerster/mpb/v6/decor"
	"golang.org/x/crypto/ssh"
	"io"
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

// 远程写文件
func ScpFile(srcPath string, dstPath string, server ServerInternal, mode os.FileMode) {
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

func RemoteWriteFile(b []byte, dstPath string, instance ServerInternal, mode os.FileMode) {
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

//// 清理目录
//func (server Server) DelDirectory(dir string) error {
//
//	log.Printf("删除目录: %s:%s", server.Host, dir)
//	sftp, err := sftpConnect(server.Username, server.Password, server.Host, server.Port)
//	defer sftp.Close()
//	if err != nil {
//		log.Println(err.Error())
//	}
//
//	f, err := sftp.Stat(dir)
//	if err != nil {
//		return err
//	}
//
//	if !f.IsDir() {
//		log.Printf("%s不是目录...", dir)
//	}
//	server.RemoteShell(fmt.Sprintf("rm -rf %s", dir))
//	return nil
//
//}

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

//func (server Server) RemoteShell(cmd string) ShellResult {
//
//	var result ShellResult
//	if len(cmd) < 60 {
//		result.Cmd = cmd
//	} else {
//		result.Cmd = "built-in shell"
//	}
//	result.Host = server.Host
//
//	log.Printf("-> [%s] => exec shell...", server.Host)
//	session, conErr := server.sshConnect()
//	if conErr != nil {
//		log.Fatal(conErr)
//	}
//
//	defer session.Close()
//
//	combo, runErr := session.CombinedOutput(cmd)
//
//	if runErr != nil {
//		e, _ := runErr.(*ssh.ExitError)
//		result.Code = e.ExitStatus()
//		result.Stderr = runErr
//		result.StdErrMsg = runErr.Error()
//		log.Print(runErr.Error())
//	}
//
//	if result.Code == 0 {
//		log.Printf("<- call back [%s] exec shell successful...", server.Host)
//		result.Status = "success"
//		result.StdOut = string(combo)
//	} else {
//		log.Printf("<- call back [%s]\n %s", server.Host, string(combo))
//		result.Status = "failed"
//	}
//
//	return result
//}

//func (server Server) RemoteShellOutput(cmd string) {
//
//	var result ShellResult
//	if len(cmd) < 60 {
//		result.Cmd = cmd
//	} else {
//		result.Cmd = "built-in shell"
//	}
//	result.Host = server.Host
//
//	//log.Printf("-> [%s] shell => %s", server.Server, cmd)
//	log.Printf("-> [%s] => exec shell...", server.Host)
//	session, conErr := server.sshConnect()
//	if conErr != nil {
//		log.Fatal(conErr)
//	}
//
//	defer session.Close()
//
//	// combo, runErr := session.CombinedOutput(cmd)
//	session.Stderr = os.Stderr
//	session.Stdout = os.Stdout
//	session.Run(cmd)
//}

// HomeDir /root or /home/username
func (server ServerInternal) HomeDir() string {
	switch server.Username {
	case "root":
		return "/root"
	default:
		return fmt.Sprintf("/home/%s", server.Username)
	}
}

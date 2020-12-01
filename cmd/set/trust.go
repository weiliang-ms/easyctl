package set

import (
	"easyctl/constant"
	"easyctl/util"
	"fmt"
	"log"
	"os"
)

const (
	createSSHDirectoryMsg = "创建~/.ssh目录"
	handoutAKMsg          = "分发authorized_keys"
	handoutRsaMsg         = "分发id_rsa"
	handoutPubMsg         = "分发id_rsa.pub"
)

type trust struct {
	serverList    []util.Server
	rsaContent    string
	rsaPubContent string
	server        util.Server
	rasPath       string // id_rsa等文件路径
}

// 设置主机互信
// todo 待优化代码
func (trust *trust) setTrust() {

	// 1.解析主机互信主机组
	trust.parseServerList()

	// 2.查询是否存在，不存在创建
	trust.searchPubFile()

	// 3.读取文件id_rsa内容
	trust.read()

	// 3.创建目录
	trust.mkdir()

	// 4.分发
	trust.handout()

}

func (trust *trust) parseServerList() {
	if trustServerListFile != "" {
		if _, err := os.Stat(trustServerListFile); err == nil {
			util.PrintDirectBanner([]string{"parse"}, parseYumServerListFile)
			trust.serverList = util.ParseServerList(trustServerListFile).PubServerList
			fmt.Println("[pub]")
			for _, v := range trust.serverList {
				fmt.Printf("Host=%s Port=%s Username=%s Password=%s\n",
					v.Host, v.Port, v.Username, v.Password)
			}
		} else {
			log.Fatal(err.Error())
		}
	}
}

func (trust *trust) servers() {

}

// 生成pub文件
func (trust *trust) searchPubFile() {

	// 判断是否存在id_rsa.*文件
	util.PrintDirectBanner([]string{trust.server.Host}, "判断是否存在id_rsa.*文件")
	if !trust.server.FileDetection(trust.rasPath) {
		trust.generatePubFile()
	}

}

func (trust *trust) setValue() {
	trust.server = trust.serverList[0]
	trust.rasPath = fmt.Sprintf("%s/.ssh/id_rsa.*", trust.server.HomeDir())
}

func (trust *trust) generatePubFile() {
	util.PrintDirectBanner([]string{trust.server.Host}, "生成id_rsa.*文件")
	trust.server.RemoteShellReturnStd(constant.CreateRsaAndRsaPubCmd)
}

// 分发公私钥
func (trust *trust) handout() {
	for _, v := range trust.serverList {
		trust.handoutAK(v)
		trust.handoutRsa(v)
		trust.handoutPub(v)
	}
}

func (trust *trust) handoutRsaPubFile(server util.Server, name string, content string, msg string) {
	util.PrintDirectBanner([]string{server.Host}, msg)
	server.RemoteWriteFile(fmt.Sprintf("%s/.ssh/%s", server.HomeDir(), name), []byte(content))
}

func (trust *trust) handoutAK(server util.Server) {
	trust.handoutRsaPubFile(server, ak, trust.rsaPubContent, handoutAKMsg)
}

func (trust *trust) handoutRsa(server util.Server) {
	trust.handoutRsaPubFile(server, idRsa, trust.rsaContent, handoutRsaMsg)
}

func (trust *trust) handoutPub(server util.Server) {
	trust.handoutRsaPubFile(server, idRsaPub, trust.rsaPubContent, handoutPubMsg)
}

// 创建~/.ssh目录
func (trust *trust) mkdir() {
	for _, v := range trust.serverList {
		util.PrintDirectBanner([]string{v.Host}, createSSHDirectoryMsg)
		v.RemoteShellPrint(fmt.Sprintf("mkdir -p %s/.ssh && chmod 600 %s/.ssh", v.HomeDir(), v.HomeDir()))
	}
}

func (trust *trust) read() {
	trust.readRsa().readPub()
}

func (trust *trust) readRsa() *trust {
	// 获取文件内容
	util.PrintDirectBanner([]string{trust.server.Host}, "获取id_rsa内容")
	trust.rsaContent = trust.server.RemoteShellReturnStd("cat $HOME/.ssh/id_rsa")
	return trust
}

func (trust *trust) readPub() *trust {
	// 获取文件内容
	util.PrintDirectBanner([]string{trust.server.Host}, "获取id_rsa.pub内容")
	trust.rsaPubContent = trust.server.RemoteShellReturnStd("cat $HOME/.ssh/id_rsa.pub")
	return trust
}

// shell
func (trust *trust) shell(cmd string, msg string) {
	util.Shell{
		Cmd:        cmd,
		ServerList: trust.serverList,
		Banner:     trust.returnBanner(msg),
	}.Shell()
}

func (trust *trust) returnBanner(msg string) util.Banner {
	return util.Banner{
		Symbols: nil,
		Msg:     msg,
	}
}

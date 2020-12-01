package set

import (
	"easyctl/constant"
	"easyctl/util"
	"fmt"
	"log"
	"os"
)

const (
	parseYumServerListFile = "解析配置Yum的主机列表"
	backupYumRepoConfig    = "备份原有yum仓储文件"
	mountImageMsg          = "挂载系统iso镜像"
	cleanYumCacheMsg       = "清除yum缓存"
)

type yum struct {
	serverList []util.Server
}

func (yum *yum) setYum() {
	yum.getServerList()
	yum.parseFlags()
}

func (yum *yum) parseFlags() {
	if yumRepo != "" {
		yum.setRepo()
	}
}

// docker server list赋值
func (yum *yum) getServerList() {
	if yumServerListFile != "" {
		if _, err := os.Stat(yumServerListFile); err == nil {
			util.PrintDirectBanner([]string{"parse"}, parseYumServerListFile)
			yum.serverList = util.ParseServerList(yumServerListFile).YumServerList
			for _, v := range yum.serverList {
				fmt.Printf("Host=%s Port=%s Username=%s Password=%s\n",
					v.Host, v.Port, v.Username, v.Password)
			}
		} else {
			log.Fatal(err.Error())
		}
	}
}

func (yum *yum) setRepo() {
	switch yumRepo {
	case ali:
		yum.setAliRepo()
	case local:
		yum.setLocalRepo()
	default:
		log.Fatal("暂不支持")
	}

	// 清除缓存
	yum.cleanCache()
}

// 备份原有
func (yum *yum) backupRepo() {
	yum.shell(constant.BackupYumRepoCmd, backupYumRepoConfig)
}

// 配置阿里云仓库
func (yum *yum) setAliRepo() {
	yum.backupRepo()
	// base
	util.WriteFile(constant.AliBaseRepoPath, []byte(constant.CentOSAliBaseYumContent), yum.serverList)
	// epel
	util.WriteFile(constant.AliEpelRepoPath, []byte(constant.CentOSAliEpelYumContent), yum.serverList)
}

func (yum *yum) setLocalRepo() {
	yum.backupRepo()
	yum.mountImage()
	// base
	util.WriteFile(constant.YumLocalRepoPath, []byte(constant.CentOSLocalYumContent), yum.serverList)
}

// 挂载镜像
func (yum *yum) mountImage() {
	cmd := fmt.Sprintf("mount -o loop %s /media", imageFilePath)
	yum.shell(cmd, mountImageMsg)
}

// 清除缓存
func (yum *yum) cleanCache() {
	yum.shell(constant.CleanYumCacheCmd, cleanYumCacheMsg)
}

// shell
func (yum *yum) shell(cmd string, msg string) {
	util.Shell{
		Cmd:        cmd,
		ServerList: yum.serverList,
		Banner:     yum.returnBanner(msg),
	}.Shell()
}

func (yum *yum) returnBanner(msg string) util.Banner {
	return util.Banner{
		Symbols: nil,
		Msg:     msg,
	}
}

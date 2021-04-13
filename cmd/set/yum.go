package set

import (
	"easyctl/asset"
	"easyctl/pkg/runner"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
)

var (
	repoDir="/etc/yum.repos.d/"
	aliBaseByte []byte
	aliEpelByte []byte
)


const (
	aliBaseRepo = "/etc/yum.repos.d/ali-base.repo"
	aliEpelRepo = "/etc/yum.repos.d/ali-epel.repo"
)


// 配置yum repo
var setYumRepoCmd = &cobra.Command{
	Use:     "yum-repo",
	Short:   "easyctl set yum-repo [flags]",
	Example: "\neasyctl set yum-repo",
	Run: func(cmd *cobra.Command, args []string) {
		setYumRepo()
	},
}

func init() {
	aliBaseByte, _ = asset.Asset("static/conf/yum-ali-base.repo")
	aliEpelByte, _ = asset.Asset("static/conf/yum-ali-epel.repo")
	setYumRepoCmd.Flags().BoolVarP(&aliRepo, "ali-repo", "", false, "阿里云镜像源yum仓库")
	//setYumRepoCmd.Flags().BoolVarP(&tsinghuaRepo, "tsinghua-repo", "", false, "清华镜像源yum仓库")
	setYumRepoCmd.Flags().BoolVarP(&multiNode, "multi-node", "", false, "是否配置多节点")
	setYumRepoCmd.Flags().StringVarP(&serverListFile, "server-list", "", "server.yaml", "服务器列表")
}

// 配置yum repo
// todo: 优化判断语句
func setYumRepo() {
	if !multiNode {
		setLocalYumRepo()
	}else {
		setRemoteYumRepo()
	}
}

// 配置yum repo
func setLocalYumRepo() {

	// backup repo file
	os.MkdirAll(repoDir+"bak", 0644)
	log.Println("开始备份，yum仓库配置文件...")
	files, _ := ioutil.ReadDir(repoDir)
	for _, f := range files {
		if !f.IsDir() {
			oldpath := repoDir + f.Name()
			newPath := repoDir + "bak/" + f.Name()
			log.Printf("%s => %s", oldpath, newPath)
			err := os.Rename(oldpath, newPath)
			if err != nil {
				log.Fatal(err.Error())
			}
		}
	}


	if aliRepo {
		log.Printf("write repo config file -> %s",aliBaseRepo)
		ioutil.WriteFile(aliBaseRepo,aliBaseByte,0644)
		log.Printf("write repo config file -> %s",aliEpelRepo)
		ioutil.WriteFile(aliEpelRepo,aliEpelByte,0644)
	}

	log.Println("配置yum repo成功...")
}

// 配置yum repo
func setRemoteYumRepo() {

	list := runner.ParseServerList(serverListFile)

	for _, v := range list.Server{

		log.Printf("[%s] 备份repo文件",v.Host)
		v.MoveDirFiles(repoDir,repoDir+"/bak")

		log.Printf("[%s] write repo config file -> %s",v.Host,aliBaseRepo)
		v.WriteRemoteFile(aliBaseByte,aliBaseRepo,0644)
		log.Printf("[%s] write repo config file -> %s",v.Host,aliEpelRepo)
		v.WriteRemoteFile(aliEpelByte,aliEpelRepo,0644)
	}

	log.Println("配置yum repo成功...")
}

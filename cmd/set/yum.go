package set

import (
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
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
	setYumRepoCmd.Flags().BoolVarP(&aliRepo, "ali-repo", "", false, "阿里云镜像源yum仓库")
	setYumRepoCmd.Flags().BoolVarP(&tsinghuaRepo, "tsinghua-repo", "", false, "清华镜像源yum仓库")
	setYumRepoCmd.Flags().BoolVarP(&multiNode, "multi-node", "", false, "是否配置多节点")
	setYumRepoCmd.Flags().StringVarP(&serverListFile, "server-list", "", "server.yaml", "服务器列表")
}

// 配置yum repo
func setYumRepo() {
	if !multiNode {
		setLocalYumRepo()
	}
}

// 配置yum repo
func setLocalYumRepo() {
	//var re run.ExecResult

	// backup repo file
	repoDir := "/etc/yum.repos.d/"
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
}

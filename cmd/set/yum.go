package set

import (
	"github.com/spf13/cobra"
)

var repoUrl string
var isoPath string

func init() {
	yumRepoCmd.Flags().StringVarP(&repoUrl, "repo-url", "", "http://mirrors.163.com", "yum仓库地址，默认为")
	yumRepoCmd.Flags().StringVarP(&isoPath, "iso-path", "", "", "iso文件url地址或离线文件地址")
}

// 配置yum repo
// todo:
var yumRepoCmd = &cobra.Command{
	Use:   "yum-repo",
	Short: "配置yum仓库",
	Example: "\neasyctl set yum-repo" +
		"\neasyctl set yum-repo --repo-url=http://mirrors.163.com" +
		"\neasyctl set yum-repo --iso-path=CentOS-7-x86_64-DVD-2009.iso" +
		"\neasyctl set yum-repo --iso-path=https://mirrors.ustc.edu.cn/centos/7.9.2009/isos/x86_64/CentOS-7-x86_64-DVD-2009.iso",
	Run: func(cmd *cobra.Command, args []string) {
		//setYumRepo()
	},
}

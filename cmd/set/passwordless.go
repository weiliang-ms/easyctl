package set

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/asset"
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"io/ioutil"
	"log"
	"os/user"
)

func init() {
	setPasswordLessCmd.Flags().StringVarP(&serverListFile, "server-list", "", "server.yaml", "服务器列表")
}

// 主机互信
var setPasswordLessCmd = &cobra.Command{
	Use:     "password-less [flags]",
	Short:   "easyctl set password-less --server-list=xxx",
	Example: "\neasyctl set password-less --server-list=server.yaml",
	//Aliases: []string{"pka"},
	Args: cobra.ExactValidArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		passwordLess()
	},
}

func passwordLess() {

	local("生成互信文件", passwordScript())

	// 解析主机列表
	list := runner.ParseServerList(serverListFile)

	// 拷贝文件
	u, _ := user.Current()

	rsa, _ := ioutil.ReadFile(fmt.Sprintf("%s/.ssh/id_rsa", u.HomeDir))
	rsaPub, _ := ioutil.ReadFile(fmt.Sprintf("%s/.ssh/id_rsa.pub", u.HomeDir))
	authorizedKeys, _ := ioutil.ReadFile(fmt.Sprintf("%s/.ssh/authorized_keys", u.HomeDir))

	for _, v := range list.Server {

		v.RemoteShell(fmt.Sprintf("mkdir -p %s/.ssh", runner.HomeDir(v)))

		log.Printf("传输数据文件%s至%s...", rsa, v.Host)
		runner.RemoteWriteFile(rsa, fmt.Sprintf("%s/.ssh/id_rsa", runner.HomeDir(v)), v, 0600)
		log.Println("-> done 传输完毕...")

		log.Printf("传输数据文件%s至%s...", rsa, v.Host)
		runner.RemoteWriteFile(rsaPub, fmt.Sprintf("%s/.ssh/id_rsa.pub", runner.HomeDir(v)), v, 0600)
		log.Println("-> done 传输完毕...")

		log.Printf("传输数据文件%s至%s...", rsa, v.Host)
		runner.RemoteWriteFile(authorizedKeys, fmt.Sprintf("%s/.ssh/authorized_keys", runner.HomeDir(v)), v, 0600)
		log.Println("-> done 传输完毕...")
	}

	log.Println("主机免密配置完毕，请验证...")

}

func passwordScript() string {
	script, _ := asset.Asset("static/script/set/passwordless.sh")
	return string(script)
}

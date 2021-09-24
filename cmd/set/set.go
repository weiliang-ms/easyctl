package set

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/util"
	"k8s.io/klog"
	"os"
)

var (
	value          string //默认参数
	serverListFile string // 服务器列表配置文件
	configFile     string
	Debug          bool
)

// RootCmd set命令
var RootCmd = &cobra.Command{
	Use:   "set [OPTIONS] [flags]",
	Short: "设置指令集",
	Args:  cobra.ExactValidArgs(1),
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "配置文件")
	RootCmd.AddCommand(passwordLessCmd)
	RootCmd.AddCommand(ulimitCmd)
	RootCmd.AddCommand(dnsCmd)
	RootCmd.AddCommand(hostResolveCmd)
	RootCmd.AddCommand(timeZoneCmd)
	RootCmd.AddCommand(yumRepoCmd)
}

type Entity struct {
	Cmd *cobra.Command
	Fnc func(b []byte, debug bool) error
}

func Set(entity Entity) error {
	if configFile == "" {
		klog.Infof("检测到配置文件为空，生成配置文件样例 -> %s", util.ConfigFile)
		_ = os.WriteFile(util.ConfigFile, config, 0666)
	}

	flagset := entity.Cmd.Parent().Parent().PersistentFlags()
	debug, err := flagset.GetBool("debug")
	if err != nil {
		fmt.Println(err)
	}

	b, readErr := os.ReadFile(configFile)
	if readErr != nil {
		klog.Fatalf("读取配置文件失败")
	}

	return entity.Fnc(b, debug)
}

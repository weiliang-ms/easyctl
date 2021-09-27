package deny

import (
	_ "embed"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/weiliang-ms/easyctl/pkg/util"
	"k8s.io/klog"
	"os"
)

var (
	configFile string
)

//go:embed asset/config.yaml
var config []byte

type Entity struct {
	Cmd *cobra.Command
	Fnc func(b []byte, logger *logrus.Logger) error
}

// RootCmd 禁用命令
var RootCmd = &cobra.Command{
	Use:   "deny [OPTIONS] [flags]",
	Short: "禁用指令集",
	Args:  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "配置文件")
	RootCmd.AddCommand(denyFirewallCmd)
	RootCmd.AddCommand(denySelinuxCmd)
	RootCmd.AddCommand(denyPingCmd)
}

func Deny(entity Entity) error {
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

	logger := &logrus.Logger{}
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	return entity.Fnc(b, logger)
}

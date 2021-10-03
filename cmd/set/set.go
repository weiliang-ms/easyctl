package set

import (
	// embed
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

// RootCmd set命令
var RootCmd = &cobra.Command{
	Use:   "set [OPTIONS] [flags]",
	Short: "设置指令集",
	Args:  cobra.ExactValidArgs(1),
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "配置文件")
	subCommands := []*cobra.Command{
		dnsCmd,
		hostResolveCmd,
		passwordLessCmd,
		newPasswordCmd,
		timeZoneCmd,
		ulimitCmd,
		yumRepoCmd,
	}

	for _, v := range subCommands {
		RootCmd.AddCommand(v)
	}
}

// Entity Set实体
type Entity struct {
	Cmd           *cobra.Command
	Fnc           func(b []byte, logger *logrus.Logger) error
	DefaultConfig []byte
}

// Set 设置指令通用入口
func Set(entity Entity) error {

	if entity.DefaultConfig == nil {
		entity.DefaultConfig = config
	}

	if configFile == "" {
		klog.Infof("检测到配置文件参数为空，生成配置文件样例 -> %s", util.ConfigFile)
		_ = os.WriteFile(util.ConfigFile, entity.DefaultConfig, 0666)
		os.Exit(0)
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

	logger := logrus.New()
	if debug {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}

	return entity.Fnc(b, logger)
}

package add

import (
	// embed
	_ "embed"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	configFile string
)

// Entity 添加实体
type Entity struct {
	Cmd           *cobra.Command
	Fnc           func(b []byte, logger *logrus.Logger) error
	DefaultConfig []byte
}

//go:embed asset/config.yaml
var config []byte

func init() {
	RootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "配置文件")
	RootCmd.AddCommand(addUserCmd)
}

// RootCmd add命令
var RootCmd = &cobra.Command{
	Use:   "add [flags]",
	Short: "添加指令集",
	Args:  cobra.ExactValidArgs(1),
}

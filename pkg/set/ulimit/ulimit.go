package ulimit

import (
	_ "embed"
	"github.com/pkg/errors"
	"github.com/weiliang-ms/easyctl/pkg/exec"
	"github.com/weiliang-ms/easyctl/pkg/util"
	"gopkg.in/yaml.v2"
	"k8s.io/klog"
	"os"
)

const ulimitShell = `
sed -i ':a;$!{N;ba};s@# easyctl ulimit BEGIN.*# easyctl ulimit END@@' /etc/security/limits.conf
sed -i '/^$/N;/\n$/N;//D' /etc/security/limits.conf

cat >> /etc/security/limits.conf <<EOF
# easyctl ulimit BEGIN
* soft nofile 65536
* hard nofile 65536
* soft nproc 65536
* hard nproc 65536
# easyctl ulimit END
EOF
`

//go:embed config.yaml
var config []byte

func Config(configFile string, level util.LogLevel) error {
	if configFile == "" {
		klog.Infof("检测到配置文件为空，生成配置文件样例 -> %s", util.ConfigFile)
		os.WriteFile(util.ConfigFile, config, 0666)
	}

	b, err := os.ReadFile(configFile)
	if err != nil {
		return errors.New("读取配置文件失败")
	}
	exec, err := parseConfig(b)
	if err != nil {
		return err
	}

	if err := exec.Run(true, level); err != nil {
		return err
	}

	return nil
}

// ParseConfig 解析yaml配置
func parseConfig(content []byte) (exec.ExecutorItem, error) {

	executor := exec.ExecutorItem{}
	sl := exec.ServerList{}
	err := yaml.Unmarshal(content, &sl)
	if err != nil {
		return executor, err
	}

	executor.Server = sl.ParseServerList()
	executor.Script = ulimitShell

	return executor, nil
}

package install

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/weiliang-ms/easyctl/pkg/install/tmpl"
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"github.com/weiliang-ms/easyctl/pkg/util"
	"github.com/weiliang-ms/easyctl/pkg/util/errors"
	"gopkg.in/yaml.v2"
	"os"
)

// RedisClusterConfig redis安装配置反序列化对象
type RedisClusterConfig struct {
	RedisCluster struct {
		Paasword    string           `yaml:"paasword"`
		ClusterType RedisClusterType `yaml:"cluster-type"`
		Package     string           `yaml:"package"`
	} `yaml:"redis-cluster"`
}

// 内部对象
type redisClusterConfig struct {
	Servers    []runner.ServerInternal
	Password   string
	CluterType RedisClusterType
	Package    string
	Logger     *logrus.Logger
}

// RedisClusterType redis cluster部署模式
type RedisClusterType int8

const (
	local RedisClusterType = iota
	threeNodesThreeShards
	sixNodesThreeShards
)

const pruneRedisShell = `
rm -f /etc/redis-* || true
rm -rf /tmp/redis* || true
rm -rf /var/lib/redis || true
rm -rf /var/log/redis || true
rm -f /var/run/redis* || true
`

// RedisCluster 部署redis集群
func RedisCluster(b []byte, logger *logrus.Logger) error {
	var config RedisClusterConfig

	logger.Info("解析redis cluster安装配置")
	if err := yaml.Unmarshal(b, &config); err != nil {
		return err
	}

	// 深拷贝属性
	redisCluster := config.deepCopy()
	redisCluster.Logger = logger
	servers, err := runner.ParseServerList(b, logger)

	if err != nil {
		return fmt.Errorf("[redis-cluster] 反序列化主机列表失败 -> %v", err)
	}
	redisCluster.Servers = servers

	return install(redisCluster)
}

// Detect 调用运行时检测依赖
func (config *redisClusterConfig) Detect() error {
	const check = "gcc -v"

	if _, err := os.Stat(config.Package); err != nil {
		return errors.FileNotFoundErr(config.Package)
	}

	if config.CluterType == threeNodesThreeShards && len(config.Servers) != 3 {
		return errors.NumNotEqualErr("节点", 3, len(config.Servers))
	}

	if config.CluterType == sixNodesThreeShards && len(config.Servers) != 6 {
		return errors.NumNotEqualErr("节点", 6, len(config.Servers))
	}

	config.Logger.Infoln("检测依赖环境...")
	exec := runner.ExecutorInternal{
		Servers: config.Servers,
		Script:  check,
		Logger:  config.Logger,
	}

	if config.CluterType == local {
		return runner.LocalRun(check, config.Logger)
	}

	for v := range exec.ParallelRun() {
		if v.Err != nil {
			return fmt.Errorf("%s 依赖检测失败 -> %s", v.Host, v.Err)
		}
	}

	return nil
}

// Prune 清理历史文件
func (config *redisClusterConfig) Prune() error {

	config.Logger.Infoln("清理redis历史文件...")
	if config.CluterType == local {
		return runner.LocalRun(pruneRedisShell, config.Logger)
	}

	exec := runner.ExecutorInternal{
		Servers: config.Servers,
		Script:  pruneRedisShell,
		Logger:  config.Logger,
	}

	ch := exec.ParallelRun()
	for v := range ch {
		if v.Err != nil {
			return fmt.Errorf("[%s] 执行清理指令失败 %s", v.Host, v.Err)
		}
	}

	return nil
}

// HandPackage 分发安装包
func (config *redisClusterConfig) HandPackage() error {

	if config.CluterType == local {
		return os.Rename(config.Package, fmt.Sprintf("/tmp/%s",
			util.SubFileName(config.Package)))
	}

	config.Logger.Infoln("分发package...")
	ch := runner.ParallelScp(runner.ScpItem{
		Servers: config.Servers,
		SrcPath: config.Package,
		DstPath: fmt.Sprintf("/tmp/%s",
			util.SubFileName(config.Package)),
		Mode:   0755,
		Logger: config.Logger,
	})

	for v := range ch {
		if v != nil {
			return v
		}
	}

	config.Logger.Infoln("分发redis安装包完毕...")
	return nil
}

// Compile 编译
func (config *redisClusterConfig) Compile() error {

	compileCmd, err := util.Render(tmpl.RedisCompileTmpl, util.TmplRenderData{
		"PackageName": util.SubFileName(config.Package),
	})

	if err != nil {
		return fmt.Errorf("生成编译指令模板失败, %s", err)
	}

	exec := runner.ExecutorInternal{
		Servers: config.Servers,
		Script:  compileCmd,
		Logger:  config.Logger,
	}

	ch := exec.ParallelRun()
	for v := range ch {
		if v.Err != nil {
			return v.Err
		}
	}

	return nil
}

func (config RedisClusterConfig) deepCopy() *redisClusterConfig {
	return &redisClusterConfig{
		Password:   config.RedisCluster.Paasword,
		CluterType: config.RedisCluster.ClusterType,
		Package:    config.RedisCluster.Package,
	}
}

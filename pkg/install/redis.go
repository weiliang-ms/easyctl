package install

import (
	"fmt"
	"github.com/lithammer/dedent"
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"github.com/weiliang-ms/easyctl/pkg/util"
	"gopkg.in/yaml.v2"
	"k8s.io/klog/v2"
	"text/template"
)

var compileTmpl = template.Must(template.New("compileTmpl").Parse(dedent.Dedent(`
cd /tmp
for i in $(ls |grep redis-);do
	rm -rf $i
done
tar zxvf {{ .PackageName }}
packageName=$(sed "s#{{ .PackageName }}#.tar.gz#g")
echo $packageName
cd %packageName
sed -i "s#\$(PREFIX)/bin#%s#g" src/Makefile
make -j $(nproc)
make install
`)))

// RedisMeta redis元数据
type RedisMeta struct {
	Password string
	Package  string
}

// RedisItem 序列化对象
type RedisItem struct {
	Redis struct {
		Password string `yaml:"password"`
		Package  string `yaml:"package"`
	} `yaml:"redis"`
}

// RedisCluster 部署redis集群
func RedisCluster(b []byte, debug bool) error {
	meta, err := parseItem(b, debug)
	if err != nil {
		return err
	}
	if err := Install(&meta, b, debug); err != nil {
		return err
	}

	return nil
}

func parseItem(b []byte, debug bool) (RedisMeta, error) {
	klog.Infoln("解析配置...")
	item := RedisItem{}

	if err := yaml.Unmarshal(b, &item); err != nil {
		return RedisMeta{}, err
	}

	if debug {
		fmt.Printf("%v", item)
	}

	return RedisMeta{Password: item.Redis.Password, Package: item.Redis.Package}, nil
}

// Combine 组装执行器
func (meta *RedisMeta) Combine(servers []runner.ServerInternal) Executor {
	return Executor{Servers: servers, Meta: meta}
}

// Detect 检测依赖
func (meta *RedisMeta) Detect(executor Executor, debug bool) error {
	klog.Infoln("检测依赖环境...")
	check := "gcc -v"
	exec := runner.ExecutorInternal{
		Servers: executor.Servers,
		Script:  check,
	}
	// todo:  fix logger
	for v := range exec.ParallelRun(nil) {
		if v.Err != nil {
			return fmt.Errorf("依赖检测失败 -> %s", v.Err)
		}
	}

	return nil
}

// HandPackage 分发安装包
func (meta *RedisMeta) HandPackage(executor Executor, debug bool) error {

	klog.Infoln("分发package...")
	ch := runner.ParallelScp(executor.Servers, meta.Package, fmt.Sprintf("/tmp/%s", util.SubFileName(meta.Package)), 0755)

	count := 1
	for {
		if count > len(executor.Servers) {
			fmt.Println("receive close chan")
			break
		}

		if err, _ := <-ch; err != nil {
			return err
		}
		count++
	}

	klog.Infoln("分发package完毕...")
	return nil
}

// Compile 编译
func (meta *RedisMeta) Compile(executor Executor, debug bool) error {
	compileCmd, err := util.Render(compileTmpl, util.Data{
		"PackageName": util.SubFileName(meta.Package),
	})
	if err != nil {
		return err
	}

	exec := runner.ExecutorInternal{
		Servers: executor.Servers,
		Script:  compileCmd,
	}

	// todo fix logger
	ch := exec.ParallelRun(nil)
	for v := range ch {
		if v.Err != nil {
			return v.Err
		}
	}

	return nil
}

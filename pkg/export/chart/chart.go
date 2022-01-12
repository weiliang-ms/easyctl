package chart

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"github.com/weiliang-ms/easyctl/pkg/util/log"
	"github.com/weiliang-ms/easyctl/pkg/util/osutil"
	"github.com/weiliang-ms/easyctl/pkg/util/request"
	"gopkg.in/yaml.v3"
	"net/http"
	"sort"
	"sync"
	"time"
)

// Manager 与helm仓库交互的执行器
type Manager struct {
	Config
	Logger             *logrus.Logger
	Handler            HandlerInterface
	HttpRequestTimeout time.Duration
}

type Config struct {
	Schema      string
	Endpoint    string
	Domain      string
	Username    string
	Password    string
	PreserveDir string
	Package     bool
	RepoName    string
}

type itemList []Item

type ConfigManifest struct {
	HelmRepo struct {
		Endpoint    string `yaml:"endpoint"`
		Domain      string `yaml:"domain"`
		Username    string `yaml:"username"`
		Password    string `yaml:"password"`
		PreserveDir string `yaml:"preserveDir"`
		Package     bool   `yaml:"package"`
		RepoName    string `yaml:"repo-name"`
	} `yaml:"helm-repo"`
}

// HelmRepo helm仓库配置
type HelmRepo struct {
	Endpoint    string `yaml:"endpoint"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	PreserveDir string `yaml:"preserveDir"`
	Package     bool   `yaml:"package"`
	RepoName    string `yaml:"Repository-name"`
}

// Item chart对象，用于反序列化
type Item struct {
	Name          string    `json:"name"`
	TotalVersions int       `json:"total_versions"`
	LatestVersion string    `json:"latest_version"`
	Created       time.Time `json:"created"`
	Updated       time.Time `json:"updated"`
	Icon          string    `json:"icon"`
	Home          string    `json:"home"`
	Deprecated    bool      `json:"deprecated"`
}

// Run 批量下载chart
func Run(item command.OperationItem) command.RunErr {

	// set default
	item.Logger = log.SetDefault(item.Logger)

	item.Logger.Info("解析chart仓库配置...")
	config, err := parseConfig(item.B, item.Logger)
	if err != nil {
		return command.RunErr{Err: err}
	}

	m := Manager{
		Config: config,
		Logger: item.Logger,
	}

	m.Handler = getHandlerInterface(item.Interface)

	return command.RunErr{Err: m.Save()}
}

// ParseHelmRepoConfig 解析helm仓库配置
func parseConfig(b []byte, logger *logrus.Logger) (Config, error) {
	c := ConfigManifest{}
	if err := yaml.Unmarshal(b, &c); err != nil {
		return Config{}, err
	}
	logger.Debugf("chart导出器结构体: %v", c)
	return deepCopy(c), nil
}

func deepCopy(c ConfigManifest) Config {
	return Config{
		Schema:      "http",
		Endpoint:    c.HelmRepo.Endpoint,
		Domain:      c.HelmRepo.Domain,
		Username:    c.HelmRepo.Username,
		Password:    c.HelmRepo.Password,
		PreserveDir: c.HelmRepo.PreserveDir,
		Package:     c.HelmRepo.Package,
		RepoName:    c.HelmRepo.RepoName,
	}
}

// List 获取chart列表
func (m *Manager) List() (itemList, error) {

	// todo: ssl
	var items []Item
	b, err := m.Handler.GetChartList(request.HTTPRequestItem{
		Url:      ApiListChartsUrl(m.Config.Schema, m.Config.Endpoint),
		Method:   http.MethodGet,
		Body:     nil,
		Headers:  nil,
		User:     m.Config.Username,
		Timeout:  m.HttpRequestTimeout,
		Password: m.Config.Password,
		Mock:     false,
	})

	if err != nil {
		return nil, err
	}

	unmarshalErr := json.Unmarshal(b, &items)
	return items, unmarshalErr
}

// Save 保存chart列表
func (m *Manager) Save() error {

	list, err := m.List()

	if err != nil {
		return command.RunErr{Err: err}
	}

	m.Logger.Infof("待导出chart数量为: %d", len(list))

	m.Logger.Infof("导出chart -> 创建目录: %s", m.Config.PreserveDir)
	osutil.MkDirPanicErr(m.Config.PreserveDir, 0755)

	m.Logger.Info("逐一导出chart中...")

	sort.Sort(list)
	byteSlice, err := m.GetChartsByte(list)
	if err != nil {
		return err
	}

	for k, v := range byteSlice {
		m.saveChart(k, v)
	}

	m.Logger.Infof("导出完毕，chart总数为:%d", len(list))
	// todo: tar打包功能
	return nil
}

type chartByteResponse struct {
	B    []byte
	Err  error
	Path string
}

func (m *Manager) GetChartsByte(list itemList) (map[string][]byte, error) {
	result := make(map[string][]byte)

	ch := make(chan chartByteResponse, len(list))
	wg := sync.WaitGroup{}
	wg.Add(len(list))

	f := func(name string) ([]byte, error) {
		return m.Handler.GetChartByte(request.HTTPRequestItem{
			Url:      ApiChartBytesUrl(m.Config.Schema, m.Config.Endpoint, name),
			Method:   http.MethodGet,
			Body:     nil,
			Headers:  nil,
			User:     m.Config.Username,
			Timeout:  m.HttpRequestTimeout,
			Password: m.Config.Password,
			Mock:     false,
		})
	}

	for _, v := range list {
		name := fmt.Sprintf("%s-%s.tgz", v.Name, v.LatestVersion)
		go func(rName string) {
			b, err := f(rName)
			ch <- chartByteResponse{
				B:    b,
				Err:  err,
				Path: fmt.Sprintf("%s/%s", m.Config.PreserveDir, rName),
			}
			wg.Done()
		}(name)

	}

	wg.Wait()
	close(ch)

	for v := range ch {
		if v.Err != nil {
			return result, v.Err
		}
		result[v.Path] = v.B
	}

	return result, nil
}

// SaveChart 持久化chart文件
func (m *Manager) saveChart(filepath string, b []byte) {
	count := osutil.WriteFilePanicErr(filepath, b)
	m.Logger.Debugf("生成文件:%s内容，文件大小为: %d", filepath, count)
}

func (list itemList) Len() int { return len(list) }

func (list itemList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list itemList) Less(i, j int) bool {
	return list[i].Name > list[j].Name
}

func getHandlerInterface(i interface{}) HandlerInterface {
	handlerInterface, _ := i.(HandlerInterface)
	if handlerInterface == nil {
		return new(Handler)
	}
	return handlerInterface
}

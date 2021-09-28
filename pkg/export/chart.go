package export

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io"
	"net/http"
	"os"
	"time"
)

// ChartRepoConfig chart仓库实体
type ChartRepoConfig struct {
	HelmRepo HelmRepo `yaml:"helm-Repository"`
}

type ChartExecutor struct {
	Config ChartRepoConfig
	Logger *logrus.Logger
}

type HelmRepo struct {
	Endpoint    string `yaml:"endpoint"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	PreserveDir string `yaml:"preserveDir"`
	Package     bool   `yaml:"package"`
	RepoName    string `yaml:"Repository-name"`
}

type ChartItem struct {
	Name          string    `json:"name"`
	TotalVersions int       `json:"total_versions"`
	LatestVersion string    `json:"latest_version"`
	Created       time.Time `json:"created"`
	Updated       time.Time `json:"updated"`
	Icon          string    `json:"icon"`
	Home          string    `json:"home"`
	Deprecated    bool      `json:"deprecated"`
}

func Chart(b []byte, logger *logrus.Logger) error {
	// todo: 生成样例文件
	logger.Info("解析chart仓库配置...")
	config, err := ParseHelmRepoConfig(b, logger)
	if err != nil {
		return err
	}

	executor := ChartExecutor{
		Config: config,
		Logger: logger,
	}

	list, err := executor.ChartList()
	if err != nil {
		return err
	}

	logger.Infof("待导出chart数量为: %d", len(list))

	executor.Save(list)

	return nil
}

func ParseHelmRepoConfig(b []byte, logger *logrus.Logger) (ChartRepoConfig, error) {
	config := ChartRepoConfig{}
	if err := yaml.Unmarshal(b, &config); err != nil {
		return ChartRepoConfig{}, err
	} else {
		logger.Debugf("chart导出器结构体: %v", config)
		return config, nil
	}
}

// ChartList 获取chart列表
func (executor ChartExecutor) ChartList() ([]ChartItem, error) {
	var items []ChartItem
	// todo ssl
	url := fmt.Sprintf("http://%s/api/chartrepo/charts/charts", executor.Config.HelmRepo.Endpoint)
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(executor.Config.HelmRepo.Username, executor.Config.HelmRepo.Password)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	b, _ := io.ReadAll(resp.Body)
	unmarshalErr := json.Unmarshal(b, &items)
	if unmarshalErr != nil {
		return nil, unmarshalErr
	}

	return items, nil
}

// Save
// 保存chart列表
func (executor ChartExecutor) Save(list []ChartItem) {

	executor.Logger.Info("导出chart...")

	executor.Logger.Infof("创建目录: %s\n", executor.Config.HelmRepo.PreserveDir)
	err := os.Mkdir(executor.Config.HelmRepo.PreserveDir, 0755)
	if err != nil {
		if err.(*os.PathError) == nil {
			panic(err)
		}
	}

	executor.Logger.Info("逐一导出chart中...")

	for _, v := range list {
		name := fmt.Sprintf("%s-%s.tgz", v.Name, v.LatestVersion)
		url := fmt.Sprintf("http://%s/chartrepo/charts/charts/%s", executor.Config.HelmRepo.Endpoint, name)
		resp, err := executor.doGET(url)
		if err != nil {
			executor.Logger.Fatal(err)
		}
		path := fmt.Sprintf("%s/%s", executor.Config.HelmRepo.PreserveDir, name)
		executor.SaveChart(path, resp.Body)
	}

	executor.Logger.Infof("导出完毕，chart总数为:%d", len(list))
	// todo: tar打包功能
}

func (executor ChartExecutor) doGET(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(executor.Config.HelmRepo.Username, executor.Config.HelmRepo.Password)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func (executor ChartExecutor) SaveChart(dst string, src io.Reader) {

	executor.Logger.Debugf("创建文件句柄: %s", dst)
	out, err := os.Create(dst)
	defer out.Close()

	executor.Logger.Debugf("生成文件:%s内容", dst)
	count, err := io.Copy(bufio.NewWriter(out), src)
	if err != nil {
		panic(err)
	}

	executor.Logger.Debugf("文件大小为: %d", count)
}

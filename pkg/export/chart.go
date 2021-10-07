package export

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"github.com/weiliang-ms/easyctl/pkg/util/log"
	"gopkg.in/yaml.v3"
	"io"
	"net/http"
	"os"
	"time"
)

// ChartRepoConfig chart仓库实体
type ChartRepoConfig struct {
	HelmRepo HelmRepo `yaml:"helm-Repository"`
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

// ChartItem chart对象，用于反序列化
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

const GetChartListFunc = "getChartListFunc"
const GetChartsByteFunc = "getChartsByteFunc"

type getChartListFunc func(endpoint, user, password string) ([]byte, error)
type getChartsByteFunc func(list []ChartItem) (map[string][]byte, error)

// Chart 批量下载chart
func Chart(item command.OperationItem) error {

	// set default
	item.Logger = log.SetDefault(item.Logger)

	item.Logger.Info("解析chart仓库配置...")
	config, err := parseHelmRepoConfig(item.B, item.Logger)
	if err != nil {
		return err
	}

	listFunc, ok := item.OptionFunc[GetChartListFunc].(func(endpoint, user, password string) ([]byte, error))
	if !ok {
		return fmt.Errorf("%s 入参非法", GetChartListFunc)
	}

	chartBytes, ok := item.OptionFunc[GetChartsByteFunc].(func(list []ChartItem) (map[string][]byte, error))
	if !ok {
		return fmt.Errorf("%s 入参非法", GetChartsByteFunc)
	}

	executor := ChartExecutor{
		Config:         config,
		Logger:         item.Logger,
		ChartsByteFunc: chartBytes,
		ChartListFunc:  listFunc,
	}

	list, err := executor.List()

	if err != nil {
		return err
	}

	item.Logger.Infof("待导出chart数量为: %d", len(list))

	return executor.Save(list)
}

// ParseHelmRepoConfig 解析helm仓库配置
func parseHelmRepoConfig(b []byte, logger *logrus.Logger) (ChartRepoConfig, error) {
	config := ChartRepoConfig{}
	if err := yaml.Unmarshal(b, &config); err != nil {
		return ChartRepoConfig{}, err
	}
	logger.Debugf("chart导出器结构体: %v", config)
	return config, nil
}

// List 获取chart列表
func (executor *ChartExecutor) List() ([]ChartItem, error) {
	// todo ssl

	var items []ChartItem
	b, err := executor.ChartListFunc(executor.Config.HelmRepo.Endpoint,
		executor.Config.HelmRepo.Username,
		executor.Config.HelmRepo.Password)

	if err != nil {
		return nil, err
	}

	unmarshalErr := json.Unmarshal(b, &items)
	if unmarshalErr != nil {
		return nil, unmarshalErr
	}

	return items, nil
}

func GetChartList(endpoint, user, password string) ([]byte, error) {

	url := fmt.Sprintf("http://%s/api/chartrepo/charts/charts", endpoint)
	resp, err := doGET(url, user, password)

	if err != nil {
		return nil, err
	}

	return io.ReadAll(resp.Body)
}

// Save 保存chart列表
func (executor ChartExecutor) Save(list []ChartItem) error {

	executor.Logger.Info("导出chart -> 创建目录: %s", executor.Config.HelmRepo.PreserveDir)

	if err := os.Mkdir(executor.Config.HelmRepo.PreserveDir, 0755); err != nil {
		return err
	}

	executor.Logger.Info("逐一导出chart中...")

	byteSlice, err := executor.ChartsByteFunc(list)
	if err != nil {
		return err
	}

	for k, v := range byteSlice {
		err := saveChart(executor.Logger, k, v)
		if err != nil {
			return err
		}
	}

	executor.Logger.Infof("导出完毕，chart总数为:%d", len(list))
	// todo: tar打包功能

	return nil
}

func GetChartsByte(list []ChartItem, endpoint, user, password, preserveDir string) (map[string][]byte, error) {
	result := make(map[string][]byte)
	for _, v := range list {
		name := fmt.Sprintf("%s-%s.tgz", v.Name, v.LatestVersion)
		url := fmt.Sprintf("http://%s/chartrepo/charts/charts/%s", endpoint, name)
		resp, err := doGET(url, user, password)
		if err != nil {
			return result, err
		}

		path := fmt.Sprintf("%s/%s", preserveDir, name)

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return result, err
		}

		result[path] = b

	}

	return result, nil
}

func doGET(url string, user, password string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(user, password)

	return http.DefaultClient.Do(req)
}

// SaveChart 持久化chart文件
func saveChart(logger *logrus.Logger, dst string, src []byte) error {

	logger.Debugf("创建文件句柄: %s", dst)
	out, err := os.Create(dst)
	defer out.Close()

	logger.Debugf("生成文件:%s内容", dst)
	count, err := out.Write(src)
	if err != nil {
		return err
	}

	logger.Debugf("文件大小为: %d", count)
	return nil
}

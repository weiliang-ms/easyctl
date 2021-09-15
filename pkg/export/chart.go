package export

import (
	"bufio"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type ChartExporter struct {
	HelmRepo HelmRepo `yaml:"helm-repo"`
}

type HelmRepo struct {
	Endpoint    string `yaml:"endpoint"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	PreserveDir string `yaml:"preserveDir"`
	Package     bool   `yaml:"package"`
	RepoName    string `yaml:"repo-name"`
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

func Chart(configPath string) (err error) {
	// todo: 生成样例文件
	exporter, err := parseHelmExporter(configPath)
	if err != nil {
		return err
	}

	err = exporter.chart()

	return nil
}

func parseHelmExporter(configPath string) (*ChartExporter, error) {

	exporter := &ChartExporter{}

	f, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(exporter)

	log.Printf("chart导出器结构体: %v", exporter)
	return exporter, err
}

func (export ChartExporter) chart() error {
	log.Println("导出逻辑...")
	list, err := export.chartList()
	if err != nil {
		return err
	}
	log.Printf("待导出chart数量为: %d", len(list))

	export.save(list)

	return nil
}

// 获取chart列表
func (export ChartExporter) chartList() ([]ChartItem, error) {
	var items []ChartItem
	url := fmt.Sprintf("http://%s/api/chartrepo/charts/charts", export.HelmRepo.Endpoint)
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(export.HelmRepo.Username, export.HelmRepo.Password)
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

// 保存chart列表
func (export ChartExporter) save(list []ChartItem) {

	log.Printf("创建目录: %s\n", export.HelmRepo.PreserveDir)
	err := os.Mkdir(export.HelmRepo.PreserveDir, 0755)
	if err != nil {
		panic(err)
	}

	for _, v := range list {
		name := fmt.Sprintf("%s-%s.tgz", v.Name, v.LatestVersion)
		url := fmt.Sprintf("http://%s/chartrepo/charts/charts/%s", export.HelmRepo.Endpoint, name)
		resp, err := export.doGET(url)
		if err != nil {
			log.Fatal(err)
		}

		path := fmt.Sprintf("%s/%s", export.HelmRepo.PreserveDir, name)
		export.saveChart(path, resp.Body)

	}

	// todo: tar打包功能
}

func (export ChartExporter) doGET(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(export.HelmRepo.Username, export.HelmRepo.Password)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func (export ChartExporter) saveChart(dst string, src io.Reader) {

	log.Printf("创建文件句柄: %s", dst)
	out, err := os.Create(dst)
	defer out.Close()

	log.Printf("生成文件:%s内容", dst)
	count, err := io.Copy(bufio.NewWriter(out), src)
	if err != nil {
		panic(err)
	}

	log.Printf("文件大小为: %d", count)
}

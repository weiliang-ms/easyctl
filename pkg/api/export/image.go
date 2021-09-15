package export

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/weiliang-ms/easyctl/pkg/util"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

// 包含关系
// harbor -> 包含多个project
// project -> 包含多个repo
// repo -> 包含多个tag

type harbor struct {
	Harbor harborMeta `yaml:"harbor"`
}

type harborMeta struct {
	User        string   `yaml:"user"`
	Password    string   `yaml:"password"`
	Address     string   `yaml:"address"`
	ExportAll   bool     `yaml:"export-all"`
	Projects    []string `yaml:"projects"`
	Domain      string   `yaml:"domain"`
	PreserveDir string   `yaml:"preserve-dir"`
}

type Project struct {
	CreationTime       time.Time `json:"creation_time"`
	CurrentUserRoleId  int       `json:"current_user_role_id"`
	CurrentUserRoleIds []int     `json:"current_user_role_ids"`
	CveAllowlist       struct {
		CreationTime time.Time     `json:"creation_time"`
		Id           int           `json:"id"`
		Items        []interface{} `json:"items"`
		ProjectId    int           `json:"project_id"`
		UpdateTime   time.Time     `json:"update_time"`
	} `json:"cve_allowlist"`
	Metadata struct {
		Public               string `json:"public"`
		RetentionId          string `json:"retention_id,omitempty"`
		AutoScan             string `json:"auto_scan,omitempty"`
		EnableContentTrust   string `json:"enable_content_trust,omitempty"`
		PreventVul           string `json:"prevent_vul,omitempty"`
		ReuseSysCveAllowlist string `json:"reuse_sys_cve_allowlist,omitempty"`
		Severity             string `json:"severity,omitempty"`
	} `json:"metadata"`
	Name       string    `json:"name"`
	OwnerId    int       `json:"owner_id"`
	OwnerName  string    `json:"owner_name"`
	ProjectId  int       `json:"project_id"`
	RepoCount  int       `json:"repo_count,omitempty"`
	UpdateTime time.Time `json:"update_time"`
	ChartCount int       `json:"chart_count,omitempty"`
}

type HarborExecutor struct {
	HarborMeta          harborMeta
	Projects            []Project
	ImageTagsWithinRepo map[string][]string
}

type repo struct {
	Name string `json:"name"`
}

type tags struct {
	Tags []tag `json:"tags"`
}

type tag struct {
	Name string `json:"name"`
}

//go:embed config.yaml
var defaultConfig []byte

// ImageList HarborImage 单机本地离线
func ImageList(path string) {

	if path == "" {
		log.Println("未检测到配置文件，已生成默认配置文件, 请调整内容后重新执行...")
		_ = os.WriteFile("config.yaml", defaultConfig, 0755)
		return
	}

	//var p []project
	hm := parseHarbor(path).Harbor
	excutor := &HarborExecutor{}
	excutor.HarborMeta = hm

	if hm.ExportAll {
		excutor.projects()
	} else {
		excutor.projectsByName()
	}

	excutor.imagesWithTag().writeImageList()
}

// 获取所有projects
func (executor *HarborExecutor) projects() *HarborExecutor {
	//http://192.168.174.95/api/v2.0/projects?page=1&page_size=15
	page := 0

	var projects []Project
	for {
		page++
		url := fmt.Sprintf("http://%s/api/v2.0/projects?page=%d&page_size=15", executor.HarborMeta.Address, page)

		var p []Project
		log.Printf("请求url为：%s\n", url)
		resp, err := executor.doRequest(url, http.MethodGet, nil, nil)
		if err != nil {
			panic(err)
		}

		b, err := io.ReadAll(resp.Body)
		//log.Printf("response结果为: %s\n", string(b))

		if err != nil {
			panic(err)
		}
		if err := json.Unmarshal(b, &p); err != nil {
			panic(err)
		}

		//log.Printf("反序列化对象结果: %+v\n", p)

		for _, v := range p {
			projects = append(projects, v)
		}

		if len(p) < 1 {
			executor.Projects = projects
			return executor
		}
	}
}

func (executor HarborExecutor) doRequest(url, method string, body io.Reader, headers map[string]string) (*http.Response, error) {
	request, err := http.NewRequest(method, url, body)
	request.SetBasicAuth(executor.HarborMeta.User, executor.HarborMeta.Password)
	if err != nil {
		return nil, err
	}
	client := http.DefaultClient
	return client.Do(request)
}

// 解析自定义导出项目组下repo列表
func (executor *HarborExecutor) imagesWithTag() *HarborExecutor {

	repoProjectMap := make(map[string][]string)
	// 遍历项目名称
	for _, v := range executor.Projects {
		//fmt.Printf("project: %s", v)
		log.Printf("导出项目：%s下镜像列表\n", v.Name)
		// 遍历v项目下repo列表
		tags := []string{}
		for _, r := range executor.reposWithinProject(v.Name) {
			// 输出tag
			for _, t := range executor.tagsWithinRepo(v.Name, r) {
				tags = append(tags, t)
			}
		}
		log.Printf("项目：%s下的repo tag集合为: %s", v.Name, tags)
		repoProjectMap[v.Name] = tags
	}

	executor.ImageTagsWithinRepo = repoProjectMap
	return executor
}

// 生成文件
func (executor *HarborExecutor) writeImageList() {

	path := executor.HarborMeta.PreserveDir
	err := os.Mkdir(path, 0755)
	if err != nil && !os.IsExist(err) {
		log.Println(err.Error())
	}

	allProjectImagesListPath := fmt.Sprintf("%s/images-list.txt", executor.HarborMeta.PreserveDir)
	if _, err := os.Stat(allProjectImagesListPath); err == nil {
		os.Remove(allProjectImagesListPath)
	}

	// 初始化image-list.txt
	allProjectImagesList, _ := os.Create(allProjectImagesListPath)

	for k, v := range executor.ImageTagsWithinRepo {
		path := fmt.Sprintf("%s%s%s", executor.HarborMeta.PreserveDir, slash(), k)
		err := os.Mkdir(path, 0755)
		if err != nil && !os.IsExist(err) {
			panic(err)
		}
		imageListPath := fmt.Sprintf("%s/%s/image-list.txt", executor.HarborMeta.PreserveDir, k)
		if _, err := os.Stat(imageListPath); err == nil {
			os.Remove(imageListPath)
		}

		// 初始化image-list.txt
		imageList, _ := os.Create(imageListPath)
		for _, t := range v {
			fmt.Println(t)
			_, _ = allProjectImagesList.WriteString(t)
			_, _ = allProjectImagesList.WriteString("\n")
			_, _ = imageList.WriteString(t)
			_, _ = imageList.WriteString("\n")
		}
	}
}

// 获取项目下repo列表
func (executor *HarborExecutor) reposWithinProject(projectName string) []string {
	var repos []string
	i := 1
	for {
		re := listRepo(i, executor.HarborMeta.User, executor.HarborMeta.Password, executor.HarborMeta.Address, projectName)
		if len(re) == 0 {
			break
		}
		for _, v := range re {
			repos = append(repos, strings.TrimPrefix(v.Name, fmt.Sprintf("%s/", projectName)))
		}
		i++
	}

	return repos
}

// 解析自定义导出项目组下repo列表
func (executor *HarborExecutor) projectsByName() *HarborExecutor {
	var projects []Project
	for _, v := range executor.HarborMeta.Projects {
		p := Project{Name: v}
		projects = append(projects, p)
	}
	executor.Projects = projects
	return executor
}

// 按`project name + repo name`查询所有镜像tag
func (executor *HarborExecutor) tagsWithinRepo(projectName string, repoName string) []string {
	var tags []string
	i := 1
	for {
		re := executor.tagsWithinRepoByPage(i, executor.HarborMeta.Address, projectName, repoName)
		if len(re) == 0 {
			break
		}
		for _, v := range re {
			for _, t := range v.Tags {
				tag := fmt.Sprintf("%s/%s/%s:%s", executor.HarborMeta.Domain, projectName, repoName, t.Name)
				tags = append(tags, tag)
			}
		}
		i++
	}
	return tags
}

// 获取repo的全部tag
func (executor *HarborExecutor) tagsWithinRepoByPage(page int, address string, projectName string, repoName string) []tags {

	var tags []tags

	args := fmt.Sprintf("page=%d"+
		"&page_size=10&with_tag=true"+
		"&with_label=false&with_scan_overview=false"+
		"&with_signature=false"+
		"&with_immutable_status=false", page)

	url := fmt.Sprintf("http://%s/api/v2.0/projects/%s/repositories/%s/artifacts?%s",
		address, projectName, strings.ReplaceAll(repoName, "/", "%252F"), args)

	//log.Printf("请求url -> %s", url)

	//var r io.Reader
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}
	//req.SetBasicAuth(username, password)
	req.Header.Add("accept", "application/json")
	req.SetBasicAuth(executor.HarborMeta.User, executor.HarborMeta.Password)
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		panic(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &tags); err != nil {
		fmt.Println(err.Error())
	}
	return tags
}

// 调用harbor api获取project列表，根据project name获取repo列表
func listRepo(page int, username string, password string, address string, projectName string) []repo {
	//log.Printf("入参为：\npage: %d \nusername: %s \npassword: %s \naddress: %s \nprojectName: %s\n",
	//page, username, password, address, projectName)

	var repo []repo
	url := fmt.Sprintf("http://%s/api/v2.0/projects/%s/repositories?page=%d&page_size=10", address, projectName, page)
	//log.Printf("request url: %s", url)

	body := util.Get(url, username, password)
	if err := json.Unmarshal(body, &repo); err != nil {
		fmt.Printf("获取项目：%s下的所有repo失败，信息：%s\n", projectName, err.Error())
	}

	return repo
}

// 解析harbor镜像导出配置文件
func parseHarbor(yamlPath string) harbor {
	var harbor harbor
	if f, err := os.Open(yamlPath); err != nil {
		log.Println("open yaml...")
		log.Fatal(err)
	} else {
		decodeErr := yaml.NewDecoder(f).Decode(&harbor)
		if decodeErr != nil {
			log.Println("decode failed...")
			log.Fatal(decodeErr)
		}
	}

	_, err := json.Marshal(harbor)

	if err != nil {
		log.Println("marshal failed...")
		log.Fatal(err)
	}

	return harbor
}

func slash() string {
	switch runtime.GOOS {
	case "linux":
		return "/"
	case "windows":
		return "\\"
	default:
		return "/"
	}
}

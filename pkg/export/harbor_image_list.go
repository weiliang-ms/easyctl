package export

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/weiliang-ms/easyctl/pkg/util"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

// 包含关系
// harbor -> 包含多个project
// project -> 包含多个repo
// Repository -> 包含多个tag

// HarborExecutor 与harbor交互的执行器
type HarborExecutor struct {
	HarborConfig   HarborConfig
	Logger         *logrus.Logger
	ProjectSlice   []ProjectInternel
	ReposInProject map[string][]string // 项目下镜像repo集合
	TagsInProject  map[string][]string // 项目下镜像tag集合
}

// HarborConfigExternel 用于反序列化harbor-repo对象配置
type HarborConfigExternel struct {
	HarborRepo struct {
		Schema      string   `yaml:"schema"`
		Address     string   `yaml:"address"`
		Domain      string   `yaml:"domain"`
		User        string   `yaml:"user"`
		Password    string   `yaml:"password"`
		PreserveDir string   `yaml:"preserve-dir"`
		Projects    []string `yaml:"projects"`
		Excludes    []string `yaml:"excludes"`
	} `yaml:"harbor-repo"`
}

// HarborConfig harbor-repo对象配置，用于内部使用
type HarborConfig struct {
	Schema      string
	Address     string
	Domain      string
	User        string
	Password    string
	PreserveDir string
	Projects    []string
	Excludes    []string
}

// ProjectInternel harbor project内部对象，包含部分必需属性
type ProjectInternel struct {
	Name      string
	ProjectId int
}

// ProjectExternel harbor project外部对象，用于反序列化
type ProjectExternel struct {
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

// Repository harbor repo外部对象，用于反序列化
type Repository struct {
	Name string `json:"name"`
}

// SearchResult 用于反序列化 查询harbor结果的对象结果集
type SearchResult struct {
	Project    []ProjectExternel `json:"project"`
	Repository []struct {
	} `json:"repository"`
	Chart []struct {
		Name  string `json:"Name"`
		Score int    `json:"Score"`
		Chart struct {
			Name        string    `json:"name"`
			Version     string    `json:"version"`
			Description string    `json:"description"`
			ApiVersion  string    `json:"apiVersion"`
			AppVersion  string    `json:"appVersion"`
			Type        string    `json:"type"`
			Urls        []string  `json:"urls"`
			Created     time.Time `json:"created"`
			Digest      string    `json:"digest"`
		} `json:"Chart"`
	} `json:"chart"`
}

// Artifact 用于反序列化制品属性
type Artifact struct {
	AdditionLinks struct {
		BuildHistory struct {
			Absolute bool   `json:"absolute"`
			Href     string `json:"href"`
		} `json:"build_history"`
		Vulnerabilities struct {
			Absolute bool   `json:"absolute"`
			Href     string `json:"href"`
		} `json:"vulnerabilities"`
	} `json:"addition_links"`
	Digest     string `json:"digest"`
	ExtraAttrs struct {
		Architecture string      `json:"architecture"`
		Author       interface{} `json:"author"`
		Created      time.Time   `json:"created"`
		Os           string      `json:"os"`
	} `json:"extra_attrs"`
	Icon              string        `json:"icon"`
	Id                int           `json:"id"`
	Labels            interface{}   `json:"labels"`
	ManifestMediaType string        `json:"manifest_media_type"`
	MediaType         string        `json:"media_type"`
	ProjectId         int           `json:"project_id"`
	PullTime          time.Time     `json:"pull_time"`
	PushTime          time.Time     `json:"push_time"`
	References        interface{}   `json:"references"`
	RepositoryId      int           `json:"repository_id"`
	Size              int           `json:"size"`
	Tags              []TagExternel `json:"tags"`
	Type              string        `json:"type"`
}

// TagExternel 用于反序列化repo tag等
type TagExternel struct {
	ArtifactId   int       `json:"artifact_id"`
	Id           int       `json:"id"`
	Immutable    bool      `json:"immutable"`
	Name         string    `json:"name"`
	PullTime     time.Time `json:"pull_time"`
	PushTime     time.Time `json:"push_time"`
	RepositoryId int       `json:"repository_id"`
	Signed       bool      `json:"signed"`
}

// HarborImageList harbor镜像列表导出
func HarborImageList(b []byte, logger *logrus.Logger) error {

	executor, err := ParseHarborConfig(b, logger)
	if err != nil {
		return err
	}

	executor.ReposInProject = make(map[string][]string)
	executor.TagsInProject = make(map[string][]string)

	task := []func() error{
		executor.ProjectList,
		executor.FilterProjects,
		executor.ReposWithinProjects,
		executor.TagsWithinProjects,
		executor.GenerateImageList,
	}

	for _, v := range task {
		if err := exec(v); err != nil {
			return err
		}
	}

	return nil
}

func exec(fnc func() error) error {
	return fnc()
}

// ParseHarborConfig 解析harbor配置
func ParseHarborConfig(b []byte, logger *logrus.Logger) (*HarborExecutor, error) {

	executor := &HarborExecutor{}
	executor.Logger = logger

	logger.Info("解析harbor配置信息...")
	config := HarborConfigExternel{}
	err := yaml.Unmarshal(b, &config)
	if err != nil {
		return nil, err
	}
	executor.HarborConfig = config.configDeepCopy()
	logger.Debugf("harbor配置信息: %v", executor.HarborConfig)
	return executor, nil
}

func (p ProjectExternel) projectDeepCopy() ProjectInternel {
	return ProjectInternel{Name: p.Name, ProjectId: p.ProjectId}
}

func (c HarborConfigExternel) configDeepCopy() HarborConfig {
	return HarborConfig{
		Schema:      c.HarborRepo.Schema,
		Address:     c.HarborRepo.Address,
		Domain:      c.HarborRepo.Domain,
		User:        c.HarborRepo.User,
		Password:    c.HarborRepo.Password,
		PreserveDir: c.HarborRepo.PreserveDir,
		Projects:    c.HarborRepo.Projects,
		Excludes:    c.HarborRepo.Excludes,
	}
}

// ProjectList 获取所有projects
func (executor *HarborExecutor) ProjectList() error {
	executor.Logger.Info("获取待导出镜像tag列表的project集合...")
	if len(executor.HarborConfig.Projects) == 0 {
		return executor.AllProjects()
	}

	return executor.ProjectsByNames(executor.HarborConfig.Projects)
}

// AllProjects 获取harbor所有projects
func (executor *HarborExecutor) AllProjects() error {
	executor.Logger.Info("获取harbor的project集合...")
	page := 0
	for {

		page++
		list, err := executor.ProjectsByPage(page)
		if err != nil {
			return err
		}

		if len(list) < 1 {
			break
		}

		for _, v := range list {
			executor.ProjectSlice = append(executor.ProjectSlice, v)
		}
	}

	return nil
}

// ProjectsByPage 按页查询harbor内project
func (executor *HarborExecutor) ProjectsByPage(page int) ([]ProjectInternel, error) {
	var projectExternel []ProjectExternel
	var projectInternel []ProjectInternel

	url := fmt.Sprintf("%s://%s/api/v2.0/projects?page=%d&page_size=15",
		executor.HarborConfig.Schema, executor.HarborConfig.Address, page)

	executor.Logger.Debugf("请求url为：%s", url)
	resp, err := executor.doRequest(url, http.MethodGet, nil, nil)
	if err != nil {
		return []ProjectInternel{}, err
	}

	b, err := io.ReadAll(resp.Body)
	executor.Logger.Debugf("response结果为: %s\n", string(b))
	if err != nil {
		return []ProjectInternel{}, err
	}
	if err := json.Unmarshal(b, &projectExternel); err != nil {
		return []ProjectInternel{}, err
	}

	executor.Logger.Debugf("反序列化对象结果: %+v\n", projectExternel)

	for _, v := range projectExternel {
		projectInternel = append(projectInternel, v.projectDeepCopy())
	}

	return projectInternel, nil
}

// ProjectsByName 按名称获取harbor内的项目
func (executor *HarborExecutor) ProjectsByName(projectName string) (ProjectInternel, error) {
	var result SearchResult
	url := fmt.Sprintf("%s://%s/api/v2.0/search?q=%s",
		executor.HarborConfig.Schema, executor.HarborConfig.Address, projectName)

	executor.Logger.Debugf("请求url为：%s", url)
	resp, err := executor.doRequest(url, http.MethodGet, nil, nil)
	if err != nil {
		return ProjectInternel{}, err
	}

	b, err := io.ReadAll(resp.Body)
	executor.Logger.Debugf("response结果为: %s\n", string(b))
	if err != nil {
		return ProjectInternel{}, err
	}
	if err := json.Unmarshal(b, &result); err != nil {
		return ProjectInternel{}, err
	}

	if len(result.Project) == 1 {
		return result.Project[0].projectDeepCopy(), nil
	}

	return ProjectInternel{}, errors.New(fmt.Sprintf("harbor为未查询到%s project", projectName))
}

// ProjectsByNames 按名称集合获取harbor内的项目
func (executor *HarborExecutor) ProjectsByNames(projectsName []string) error {
	var projects []ProjectInternel
	for _, v := range projectsName {
		p, err := executor.ProjectsByName(v)
		if err != nil {
			return err
		}
		projects = append(projects, p)
	}
	executor.ProjectSlice = projects
	return nil
}

// FilterProjects 过滤排除的project
func (executor *HarborExecutor) FilterProjects() error {
	executor.Logger.Info("过滤excludes中的project...")
	var projects []ProjectInternel
	for _, v := range executor.ProjectSlice {
		if !util.SliceContain(executor.HarborConfig.Excludes, v.Name) {
			projects = append(projects, v)
		}
	}

	executor.Logger.Infof("过滤后project集合数量为: %d -> %v...", len(projects), projects)
	executor.ProjectSlice = projects

	return nil
}

func (executor *HarborExecutor) doRequest(url, method string, body io.Reader, headers map[string]string) (*http.Response, error) {
	request, err := http.NewRequest(method, url, body)
	request.SetBasicAuth(executor.HarborConfig.User, executor.HarborConfig.Password)
	if err != nil {
		return nil, err
	}
	client := http.DefaultClient
	return client.Do(request)
}

// GenerateImageList 生成文件
func (executor *HarborExecutor) GenerateImageList() error {

	path := executor.HarborConfig.PreserveDir
	err := os.Mkdir(path, 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}

	allProjectImagesListPath := fmt.Sprintf("%s/images-list.txt", executor.HarborConfig.PreserveDir)
	if _, err := os.Stat(allProjectImagesListPath); err == nil {
		_ = os.Remove(allProjectImagesListPath)
	}

	// 初始化image-list.txt
	executor.Logger.Debugf("创建文件: %s", allProjectImagesListPath)
	allProjectImagesList, err := os.Create(allProjectImagesListPath)
	if err != nil {
		return nil
	}

	for k, v := range executor.TagsInProject {
		path := fmt.Sprintf("%s%s%s", executor.HarborConfig.PreserveDir, slash(), k)
		err := os.Mkdir(path, 0755)
		if err != nil && !os.IsExist(err) {
			panic(err)
		}
		imageListPath := fmt.Sprintf("%s/%s/image-list.txt", executor.HarborConfig.PreserveDir, k)
		if _, err := os.Stat(imageListPath); err == nil {
			_ = os.Remove(imageListPath)
		}

		// 初始化image-list.txt
		imageList, _ := os.Create(imageListPath)
		for _, t := range v {
			if executor.Logger.Level == logrus.DebugLevel {
				fmt.Println(t)
			}
			_, _ = allProjectImagesList.WriteString(t)
			_, _ = allProjectImagesList.WriteString("\n")
			_, _ = imageList.WriteString(t)
			_, _ = imageList.WriteString("\n")
		}
	}

	return nil
}

// ReposWithinProjects 获取多个项目内repo集合
func (executor *HarborExecutor) ReposWithinProjects() error {
	for _, v := range executor.ProjectSlice {
		err := executor.ReposWithinProject(v.Name)
		if err != nil {
			return err
		}
	}

	executor.Logger.Debugf("repo集合：%v", executor.ReposInProject)
	return nil
}

// ReposWithinProject 获取项目下repo列表
func (executor *HarborExecutor) ReposWithinProject(projectName string) error {
	var repos []string
	i := 1
	for {
		re, err := executor.ListRepoByPage(i, projectName)
		if err != nil {
			return fmt.Errorf("获取项目：%s第%d页 repo列表失败 ->%s", projectName, i, err)
		}
		if len(re) == 0 {
			executor.ReposInProject[projectName] = repos
			break
		}

		for _, v := range re {
			repos = append(repos, strings.TrimPrefix(v.Name, fmt.Sprintf("%s/", projectName)))
		}

		i++
	}

	return nil
}

type result struct {
	projectName string
	tags        []string
	err         error
}

// TagsWithinProjects 获取项目下镜像tag列表
func (executor *HarborExecutor) TagsWithinProjects() error {

	executor.Logger.Info("检索tag列表...")

	var wg sync.WaitGroup
	ch := make(chan result, len(executor.ProjectSlice))
	tags := make(map[string][]string)

	for k := range executor.ReposInProject {
		wg.Add(1)
		go func(projectName string) {
			re, err := executor.TagsWithinProject(projectName)
			ch <- result{
				projectName: projectName,
				tags:        re,
				err:         err,
			}
			wg.Done()
		}(k)
	}

	wg.Wait()
	close(ch)

	for v := range ch {
		if v.err != nil {
			return v.err
		}
		tags[v.projectName] = v.tags
	}

	executor.TagsInProject = tags
	executor.Logger.Debugf("%v", executor.TagsInProject)
	return nil
}

// TagsWithinProject 按`project name + Repository name`查询所有镜像tag
func (executor *HarborExecutor) TagsWithinProject(projectName string) ([]string, error) {
	var tags []string

	var wg sync.WaitGroup
	executor.Logger.Infof("解析项目%s下镜像tag集合，共计%d个tag",
		projectName, len(executor.ReposInProject[projectName]))
	ch := make(chan result, len(executor.ReposInProject[projectName]))
	for _, r := range executor.ReposInProject[projectName] {
		wg.Add(1)
		go func(repoName string) {
			re, err := executor.TagsWithinRepo(projectName, repoName)
			ch <- result{
				projectName: projectName,
				tags:        re,
				err:         err,
			}
			wg.Done()
		}(r)
	}

	wg.Wait()
	close(ch)

	for v := range ch {
		if v.err != nil {
			return tags, v.err
		}
		for _, t := range v.tags {
			tags = append(tags, t)
		}
	}

	return tags, nil
}

// TagsWithinRepo 按`project name + Repository name`查询所有镜像tag
func (executor *HarborExecutor) TagsWithinRepo(projectName string, repoName string) ([]string, error) {
	var tags []string
	i := 1

	for {
		t, err := executor.TagsWithinRepoByPage(i, projectName, repoName)
		if err != nil {
			return nil, err
		}
		if len(t) == 0 {
			break
		}
		for _, v := range t {
			tag := fmt.Sprintf("%s/%s/%s:%s", executor.HarborConfig.Domain, projectName, repoName, v)
			tags = append(tags, tag)
		}
		i++
	}
	return tags, nil
}

// TagsWithinRepoByPage 获取repo的全部tag
func (executor *HarborExecutor) TagsWithinRepoByPage(page int, projectName string, repoName string) ([]string, error) {

	var tags []Artifact

	args := fmt.Sprintf("page=%d"+
		"&page_size=10&with_tag=true"+
		"&with_label=false&with_scan_overview=false"+
		"&with_signature=false"+
		"&with_immutable_status=false", page)

	url := fmt.Sprintf("%s://%s/api/v2.0/projects/%s/repositories/%s/artifacts?%s",
		executor.HarborConfig.Schema,
		executor.HarborConfig.Address,
		projectName,
		strings.ReplaceAll(repoName, "/", "%252F"),
		args)

	executor.Logger.Debugf("请求url为: %s", url)

	resp, err := executor.doRequest(url, http.MethodGet, nil, map[string]string{"accept": "application/json"})
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &tags); err != nil {
		return nil, err
	}

	var t []string
	for _, v := range tags {
		for _, tag := range v.Tags {
			t = append(t, tag.Name)
		}
	}

	executor.Logger.Debugf("%s/%s tag集合为 -> %v",
		projectName, repoName, t)

	return t, nil
}

// ListRepoByPage 调用harbor api获取project列表，根据project name获取repo列表
func (executor *HarborExecutor) ListRepoByPage(page int, projectName string) ([]Repository, error) {

	executor.Logger.Debugf("获取项目:%s 第%d页 repo列表", projectName, page)

	var repo []Repository
	url := fmt.Sprintf("%s://%s/api/v2.0/projects/%s/repositories?page=%d&page_size=10",
		executor.HarborConfig.Schema,
		executor.HarborConfig.Address,
		projectName,
		page)

	executor.Logger.Debugf("请求url为：%s", url)

	resp, err := executor.doRequest(url, http.MethodGet, nil, nil)
	if err != nil {
		return []Repository{}, err
	}

	b, err := io.ReadAll(resp.Body)
	executor.Logger.Debugf("response结果为: %s\n", string(b))
	if err != nil {
		return []Repository{}, err
	}

	if err := json.Unmarshal(b, &repo); err != nil {
		fmt.Printf("获取项目：%s下的所有repo失败，信息：%s\n", projectName, err.Error())
	}

	executor.Logger.Debugf("page: %d 项目：%s -> 仓库集合为：%v", page, projectName, repo)

	return repo, nil
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

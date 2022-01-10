package harbor

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"github.com/weiliang-ms/easyctl/pkg/util/log"
	"github.com/weiliang-ms/easyctl/pkg/util/platform"
	"github.com/weiliang-ms/easyctl/pkg/util/request"
	"github.com/weiliang-ms/easyctl/pkg/util/slice"
	"gopkg.in/yaml.v2"
	"io/fs"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

const DefaultRequestTimeout = time.Second * 5

type ProjectNotFoundErr struct {
	ProjectName string
}

func (p ProjectNotFoundErr) Error() string {
	return fmt.Sprintf("[not found] 未查询到project -> %s", p.ProjectName)
}

/*
	1.get projects list
      - ProjectsToSearch != nil -> projects list == ProjectsToSearch
      - ProjectsToSearch == nil -> all projects in `harbor`
    2.filter projects(exclude someone) -> projects list
    3.map <project & repository> -> map[projectName][]string{repoName...}
    4.map <repository & tag>     -> map[repositoryName][]string{tagName...}
	5.map <project & project/repo:tag>     -> map[projectName][]string{domain/project/repo:ta...}
    6.write to files
*/

// ImageList harbor镜像列表导出
func ImageList(item command.OperationItem) command.RunErr {

	executor, err := ParseHarborConfig(item.B, item.Logger)
	if err != nil {
		return command.RunErr{Err: err}
	}

	executor.HandlerInterface = getHandlerInterface(item.Interface)
	executor.TagsInProject = make(map[string][]string)
	projects, err := executor.ProjectList()

	if err != nil {
		return command.RunErr{Err: err}
	}

	projects = executor.FilterProjects(projects)

	projectAndReposMap, err := executor.ReposWithinProjects(projects)
	if err != nil {
		return command.RunErr{Err: err}
	}

	projectAndTagsMap, err := executor.TagsWithinProjects(projectAndReposMap)
	if err != nil {
		return command.RunErr{Err: err}
	}

	return command.RunErr{Err: executor.GenerateImageList(projectAndTagsMap)}
}

// ParseHarborConfig 解析harbor配置
func ParseHarborConfig(b []byte, logger *logrus.Logger) (*Executor, error) {

	executor := &Executor{}
	executor.Logger = logger

	logger.Info("解析harbor配置信息...")
	config := ConfigExternal{}
	err := yaml.Unmarshal(b, &config)
	if err != nil {
		return nil, err
	}
	executor.Config = config.configDeepCopy()
	logger.Debugf("harbor配置信息: %v", executor.Config)
	return executor, nil
}

func (p ProjectExternal) projectDeepCopy() ProjectInternal {
	return ProjectInternal{Name: p.Name, ProjectId: p.ProjectId}
}

func (c ConfigExternal) configDeepCopy() Config {
	return Config{
		HarborSchema:             c.HarborRepo.Schema,
		HarborAddress:            c.HarborRepo.Address,
		HarborDomain:             c.HarborRepo.Domain,
		HarborUser:               c.HarborRepo.User,
		HarborPassword:           c.HarborRepo.Password,
		TagWithDomain:            c.HarborRepo.TagWithDomain,
		PreserveDir:              c.HarborRepo.PreserveDir,
		ProjectsToSearch:         c.HarborRepo.Projects,
		ProjectsToSearchExcludes: c.HarborRepo.Excludes,
	}
}

// ProjectList 获取所有projects
func (executor *Executor) ProjectList() ([]ProjectInternal, error) {
	executor.Logger.Info("获取待导出镜像tag列表的project集合...")
	if len(executor.ProjectsToSearch) == 0 {
		b, err := executor.HandlerInterface.ProjectCount(executor.HarborUser, executor.HarborPassword, executor.HarborSchema, executor.HarborAddress, DefaultRequestTimeout)
		if err != nil {
			return []ProjectInternal{}, err
		}

		var statistics Statistics
		_ = json.Unmarshal(b, &statistics)

		return executor.AllProjects(statistics.TotalProjectCount)
	}

	return executor.ProjectsByNames(executor.ProjectsToSearch)
}

// AllProjects 获取harbor所有projects
func (executor *Executor) AllProjects(projectCount int) (projects []ProjectInternal, err error) {

	logger := log.SetDefault(executor.Logger)
	logger.Info("获取harbor的project集合...")

	count := projectCount / 10
	for i := 1; i < count+2; i++ {
		b, err := executor.HandlerInterface.ProjectsByPage(i, 10, logger, executor.HarborSchema, executor.HarborAddress, executor.HarborUser, executor.HarborPassword, DefaultRequestTimeout)
		if err != nil {
			return projects, err
		}

		var projectExternal []ProjectExternal
		_ = json.Unmarshal(b, &projectExternal)
		logger.Debugf("反序列化对象结果: %+v\n", projectExternal)

		for _, v := range projectExternal {
			projects = append(projects, v.projectDeepCopy())
		}
	}

	return projects, nil
}

// ProjectsByName 按名称获取harbor内的项目
func (executor *Executor) ProjectsByName(projectName string) (ProjectInternal, error) {

	var result SearchResult
	b, err := executor.HandlerInterface.DoRequest(request.HTTPRequestItem{
		Url:      ApiGetProjectsByName(executor.HarborSchema, executor.HarborAddress, projectName),
		Method:   http.MethodGet,
		Body:     nil,
		Headers:  nil,
		User:     executor.HarborUser,
		Password: executor.HarborPassword,
	})

	if err != nil {
		return ProjectInternal{}, err
	}

	_ = json.Unmarshal(b, &result)

	// 接口返回数据为模糊查询数据，需要二次校验
	for _, v := range result.Project {
		if v.Name == projectName {
			return v.projectDeepCopy(), nil
		}
	}

	return ProjectInternal{}, ProjectNotFoundErr{ProjectName: projectName}
}

// ProjectsByNames 按名称集合获取harbor内的项目
func (executor *Executor) ProjectsByNames(projectsName []string) ([]ProjectInternal, error) {
	var projects []ProjectInternal
	for _, v := range projectsName {
		p, err := executor.ProjectsByName(v)
		if err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, nil
}

// FilterProjects 过滤排除的project
func (executor *Executor) FilterProjects(projects []ProjectInternal) []ProjectInternal {
	executor.Logger.Infof("过滤excludes中的project -> %s...", executor.ProjectsToSearchExcludes)
	var filterProjects []ProjectInternal
	for _, v := range projects {
		if !slice.StringSliceContain(executor.ProjectsToSearchExcludes, v.Name) {
			filterProjects = append(filterProjects, v)
		}
	}
	executor.Logger.Infof("过滤后project集合数量为: %d -> %v...", len(filterProjects), filterProjects)
	return filterProjects
}

// GenerateImageList 生成文件
func (executor *Executor) GenerateImageList(tagsInProject map[string][]string) error {

	logger := log.SetDefault(executor.Logger)
	// todo: 合法性检测
	if err := mkDirIfNotExist(executor.PreserveDir, 0755); err != nil {
		return err
	}

	// 初始化image-list.txt
	allProjectImagesListPath := fmt.Sprintf("%s/images-list.txt", executor.Config.PreserveDir)
	logger.Debugf("创建文件: %s", allProjectImagesListPath)
	allProjectImagesList, _ := os.Create(allProjectImagesListPath)
	defer allProjectImagesList.Close()

	for k, v := range tagsInProject {
		path := fmt.Sprintf("%s%s%s", executor.Config.PreserveDir, platform.Slash(runtime.GOOS), k)
		_ = mkDirIfNotExist(path, 0755)
		imageListPath := fmt.Sprintf("%s/%s/image-list.txt", executor.Config.PreserveDir, k)
		// 初始化image-list.txt
		imageList, _ := os.Create(imageListPath)
		for _, t := range v {

			tag := t
			if executor.TagWithDomain {
				tag = fmt.Sprintf("%s/%s", executor.HarborDomain, t)
			}

			executor.Logger.Info(tag)
			_, _ = allProjectImagesList.WriteString(tag)
			_, _ = allProjectImagesList.WriteString("\n")
			_, _ = imageList.WriteString(tag)
			_, _ = imageList.WriteString("\n")
		}
		imageList.Close()
	}

	return nil
}

// ReposWithinProjects 获取多个项目内repo集合
func (executor *Executor) ReposWithinProjects(projects []ProjectInternal) (map[string][]string, error) {

	repoMap := make(map[string][]string)

	for _, v := range projects {
		tags, err := executor.ReposWithinProject(v.Name, v.ProjectId)
		if err != nil {
			return repoMap, err
		}

		repoMap[v.Name] = tags
	}

	executor.Logger.Debugf("repo集合：%v", executor.ReposInProject)
	return repoMap, nil
}

// ReposWithinProject 获取项目下repo列表
func (executor *Executor) ReposWithinProject(projectName string, projectId int) ([]string, error) {

	var repos []string
	logger := log.SetDefault(executor.Logger)

	b, err := executor.HandlerInterface.RepoCount(
		executor.HarborUser,
		executor.HarborPassword,
		executor.HarborSchema,
		executor.HarborAddress,
		projectId, DefaultRequestTimeout)

	if err != nil {
		return repos, err
	}

	var meta ProjectMeta
	_ = json.Unmarshal(b, &meta)
	logger.Debugf("反序列化对象结果: %+v\n", meta)

	count := meta.RepoCount / 10
	for i := 1; i < count+2; i++ {

		data, err := executor.HandlerInterface.ListRepoByPage(
			i, 10,
			executor.HarborUser, executor.HarborPassword,
			executor.HarborSchema, executor.HarborAddress,
			projectName, DefaultRequestTimeout)
		if err != nil {
			return repos, fmt.Errorf("获取项目：%s第%d页 repo列表失败 ->%s", projectName, i, err)
		}

		var repoSlice []Repository
		_ = json.Unmarshal(data, &repoSlice)

		for _, v := range repoSlice {
			repos = append(repos, strings.TrimPrefix(v.Name, fmt.Sprintf("%s/", projectName)))
		}
	}

	return repos, nil
}

type result struct {
	err            error
	tags           []string
	projectMapTags projectMapTags
}

// TagsWithinProjects 获取项目下镜像tag列表
func (executor *Executor) TagsWithinProjects(projectAndReposMap map[string][]string) (map[string][]string, error) {

	logger := log.SetDefault(executor.Logger)
	logger.Info("检索tag列表...")

	var wg sync.WaitGroup
	ch := make(chan result, len(projectAndReposMap))
	tags := make(map[string][]string)

	for k, v := range projectAndReposMap {
		wg.Add(1)
		go func(projectName string, repos []string) {
			re, err := executor.TagsWithinProject(projectName, repos)
			ch <- result{err: err, projectMapTags: re}
			wg.Done()
		}(k, v)
	}

	wg.Wait()
	close(ch)

	for v := range ch {
		if v.err != nil {
			return tags, v.err
		}
		tags[v.projectMapTags.ProjectName] = v.projectMapTags.Tags
	}

	logger.Debugf("%v", tags)
	return tags, nil
}

// TagsWithinProject 按`project name + Repository name`查询所有镜像tag
func (executor *Executor) TagsWithinProject(projectName string, repos []string) (projectMapTags, error) {

	wg := sync.WaitGroup{}
	logger := log.SetDefault(executor.Logger)

	logger.Infof("解析项目%s下镜像tag集合，共计%d个tag",
		projectName, len(repos))
	ch := make(chan result, len(repos))

	for _, r := range repos {
		wg.Add(1)
		go func(repoName string) {
			tags, err := executor.generateRepoTagsSlice(projectName, repoName)
			ch <- result{tags: tags, err: err}
			wg.Done()
		}(r)
	}

	wg.Wait()
	close(ch)

	var tags []string
	for v := range ch {
		if v.err != nil {
			return projectMapTags{}, v.err
		}
		tags = slice.StringSliceAppend(tags, v.tags)
	}

	return projectMapTags{Tags: tags, ProjectName: projectName}, nil
}

// TagsWithinRepo 按`project name + Repository name`查询所有镜像tag
func (executor *Executor) TagsWithinRepo(projectName string, repoName string, tagsCount int) ([]string, error) {

	var tags []string

	logger := log.SetDefault(executor.Logger)
	logger.Infof("获取%s/%s下所有tag集合...", projectName, repoName)

	if tagsCount == 0 {
		return tags, fmt.Errorf("%s/%s 未查询到合法tag", projectName, repoName)
	}

	count := tagsCount / 10

	for i := 1; i < count+2; i++ {
		b, err := executor.HandlerInterface.TagsWithinRepoByPage(
			i, 10,
			executor.HarborSchema,
			executor.HarborAddress,
			projectName,
			repoName,
			executor.HarborUser,
			executor.HarborPassword,
			DefaultRequestTimeout,
		)
		if err != nil {
			return tags, err
		}

		var artifact []Artifact
		_ = json.Unmarshal(b, &artifact)

		var t []string
		for _, v := range artifact {
			for _, tag := range v.Tags {
				t = append(t, tag.Name)
			}
		}

		for _, v := range t {
			tag := fmt.Sprintf("%s/%s:%s", projectName, repoName, v)
			tags = append(tags, tag)
		}
	}

	return tags, nil
}

/*
	入参： projectName, repoName
    出参：
	[]string{
		projectName/repo1:tag1,
		projectName/repo1:tag2,
	}
*/
func (executor *Executor) generateRepoTagsSlice(projectName, repoName string) ([]string, error) {

	var tags []string
	// 获取repo下tags数量（用于后续分页查询）
	b, err := executor.HandlerInterface.TagsNumWithRepo(
		executor.HarborUser,
		executor.HarborPassword,
		executor.HarborSchema,
		executor.HarborAddress,
		projectName,
		repoName, DefaultRequestTimeout)

	if err != nil {
		return tags, err
	}

	var repoArtifactInfo RepoArtifactInfo
	_ = json.Unmarshal(b, &repoArtifactInfo)

	data, err := executor.TagsWithinRepo(projectName, repoName, repoArtifactInfo.ArtifactCount)
	if err != nil {
		return tags, err
	}

	for _, v := range data {
		tags = append(tags, v)
	}

	return tags, nil
}

func getHandlerInterface(i interface{}) HandlerInterface {
	handlerInterface, _ := i.(HandlerInterface)
	if handlerInterface == nil {
		return new(Requester)
	}
	return handlerInterface
}

func mkDirIfNotExist(dirName string, mode fs.FileMode) error {
	err := os.Mkdir(dirName, mode)
	if err != nil && strings.Contains(err.Error(), "already exists") {
		return nil
	}
	return err
}

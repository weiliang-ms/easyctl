package harbor

import (
	//
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/weiliang-ms/easyctl/pkg/export/harbor/mocks"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"github.com/weiliang-ms/easyctl/pkg/util/request"
	"net/http"
	"os"
	"runtime"
	"sort"
	"strings"
	"testing"
)

//go:embed mocks/busybox.json
var busyboxProjectInfoBytes []byte

/* elastic mock data */
//go:embed mocks/elastic.json
var elasticProjectInfoBytes []byte

//go:embed mocks/elastic-meta-info.json
var elasticMetaInfoBytes []byte

//go:embed mocks/elasticReposByPage1.json
var elasticReposByPage1Bytes []byte

//go:embed mocks/filebeat-tags-num.json
var elasticFilebeatTagsNumBytes []byte

//go:embed mocks/elasticsearch-tags-num.json
var elasticElasticsearchTagsNumBytes []byte

//go:embed mocks/filebeat-tags-by-page.json
var elasticFilebeatTagsByPageBytes []byte

//go:embed mocks/elasticsearch-tags-by-page.json
var elasticElasticsearchTagsByPageBytes []byte

//go:embed mocks/search-project.json
var searchProjectBytes []byte

//go:embed mocks/repoTagsByPage1.json
var repoTagsByPageOneBytes []byte

//go:embed mocks/repoTagsByPage2.json
var repoTagsByPageTwoBytes []byte

//go:embed mocks/projects-by-page-1.json
var projectsByPageFirstBytes []byte

//go:embed mocks/projects-by-page-2.json
var projectsByPageSecondBytes []byte

//go:embed mocks/projects-by-page-3.json
var projectsByPageThirdBytes []byte

//go:embed mocks/projects-by-page-4.json
var projectsByPageFourthBytes []byte

//go:embed mocks/projects-by-page-5.json
var projectsByPageFifthBytes []byte

//go:embed mocks/statistics.json
var statisticsBytes []byte

//go:embed mocks/kubernetes-meta-info.json
var kubernetesMetaInfoBytes []byte

//go:embed mocks/kubernetes.json
var kubernetesProjectInfoBytes []byte

//go:embed mocks/kubernetesReposByPage1.json
var kubernetesReposByPage1Bytes []byte

//go:embed mocks/kubernetesReposByPage2.json
var kubernetesReposByPage2Bytes []byte

//go:embed mocks/repo-cephcsi-tags-num.json
var cephcsiRepoTagsNumBytes []byte

//go:embed mocks/repo-cephcsi-tags.json
var cephcsiRepoTagsBytes []byte

//go:embed mocks/repo-reloader-tags-num.json
var reloaderRepoTagsNumBytes []byte

//go:embed mocks/repo-reloader-tags.json
var reloaderRepoTagsBytes []byte

//go:embed mocks/apache.json
var apacheProjectInfoBytes []byte

//go:embed mocks/apache-meta-info.json
var apacheProjectMetaInfoBytes []byte

//go:embed mocks/apacheReposByPage.json
var apacheReposByPageBytes []byte

//go:embed mocks/repo-skywalking-ui-tags-num.json
var skywalkingUITagsNumBytes []byte

//go:embed mocks/repo-skywalking-ui-tags.json
var skywalkingUITagsBytes []byte

//go:embed mocks/repo-skywalking-oap-server-tags-num.json
var skywalkingServerTagsNumBytes []byte

//go:embed mocks/repo-skywalking-oap-server-tags.json
var skywalkingServerTagsBytes []byte

//go:embed mocks/repo-skywalking-java-agent-tags-num.json
var skywalkingAgentTagsNumBytes []byte

//go:embed mocks/repo-skywalking-java-agent-tags.json
var skywalkingAgentTagsBytes []byte

/*
	moc func test
*/

var (
	MockDoRequestFunc            = "DoRequest"
	MockProjectCountFunc         = "ProjectCount"
	MockTagsWithinRepoByPageFunc = "TagsWithinRepoByPage"
	MockListRepoByPageFunc       = "ListRepoByPage"
	MockRepoCountFunc            = "RepoCount"

	mockApiGetStatisticsUrl     = fmt.Sprintf("%s://%s/api/v2.0/statistics", mockSchema, mockAddress)
	mockApiGetProjectsByPageUrl = fmt.Sprintf("%s://%s/api/v2.0/projects?page=%d&page_size=%d",
		mockSchema, mockAddress, mockPage, mockPageSize)
	mockApiGetProjectsByNameUrl = fmt.Sprintf("%s://%s/api/v2.0/search?q=%s",
		mockSchema, mockAddress, mockProjectName)
	mockApiGetListRepoByPageUrl = fmt.Sprintf("%s://%s/api/v2.0/projects/%s/repositories?page=%d&page_size=%d",
		mockSchema,
		mockAddress,
		mockProjectName,
		mockPage,
		mockPageSize,
	)
	mockApiGetProjectMetaInfoUrl = fmt.Sprintf("%s://%s/api/v2.0/projects/%d", mockSchema, mockAddress, mockProjectId)

	args = fmt.Sprintf("page=%d"+
		"&page_size=%d&with_tag=true"+
		"&with_label=false&with_scan_overview=false"+
		"&with_signature=false"+
		"&with_immutable_status=false", mockPage, mockPageSize)
	mockApiGetTagsWithinRepoByPageUrl = fmt.Sprintf("%s://%s/api/v2.0/projects/%s/repositories/%s/artifacts?%s",
		mockSchema,
		mockAddress,
		mockProjectName,
		strings.ReplaceAll(mockRepoName, "/", "%252F"),
		args)

	mockApiGetTagsNumWithRepoUrl = fmt.Sprintf("%s://%s/api/v2.0/projects/%s/repositories/%s",
		mockSchema,
		mockAddress,
		mockProjectName,
		strings.ReplaceAll(mockRepoName, "/", "%252F"),
	)

	mockSchema           = "http"
	mockAddress          = "192.168.1.1:80"
	mockUser             = "admin"
	mockPassword         = "123456"
	mockDomain           = "docker.wl.io"
	mockDataDir          = "/tmp"
	mockPage             = 1
	mockPageSize         = 10
	mockProjectName      = "kubernetes"
	mockProjectId        = 5
	mockProjectRepoCount = 15

	mockAllProjectsCount = 44
	mockProjectsToSearch = []string{
		"aaa",
		"bbb",
	}
	mockProjectsExcludedToSearch = []string{
		"ddd",
	}
	mockProjects = []ProjectInternal{
		{
			Name:      "kubernetes",
			ProjectId: 0,
		},
		{
			Name:      "library",
			ProjectId: 1,
		},
		{
			Name:      "aaa",
			ProjectId: 0,
		},
	}
	mockRepoName = "cephcsi"
	mockLogger   = logrus.New()
	mockNetErr   = fmt.Errorf("网络异常")
	mockExecutor = Executor{
		Config: Config{
			HarborSchema:             mockSchema,
			HarborAddress:            mockAddress,
			HarborDomain:             mockDomain,
			HarborUser:               mockUser,
			HarborPassword:           mockPassword,
			PreserveDir:              mockDataDir,
			ProjectsToSearch:         nil,
			ProjectsToSearchExcludes: nil,
		},
		Logger:         mockLogger,
		ProjectSlice:   nil,
		ReposInProject: nil,
		TagsInProject:  nil,
	}
)

func TestParseHarborConfig(t *testing.T) {

	content := `
harbor-repo:
  schema: http
  address: 192.168.1.1:80           # harbor连接地址
  domain: docker.wl.io              # harbor域
  user: admin                       # harbor用户
  password: 123456                  # harbor用户密码
  preserve-dir: /tmp                # 持久化tag
  projects:                         # 导出哪些项目下的镜像tag（如果为空表示全库导出）
    - aaa                        # project名称
    - bbb
  excludes:                         # 配置'projects'空值使用，过滤某些project
    - ddd
`
	logger := logrus.New()
	re, err := ParseHarborConfig([]byte(content), logger)
	require.Equal(t, nil, err)
	require.Equal(t, &Executor{
		Config: Config{
			HarborSchema:             mockSchema,
			HarborAddress:            mockAddress,
			HarborDomain:             mockDomain,
			HarborUser:               mockUser,
			HarborPassword:           mockPassword,
			PreserveDir:              mockDataDir,
			ProjectsToSearch:         mockProjectsToSearch,
			ProjectsToSearchExcludes: mockProjectsExcludedToSearch,
		},
		Logger:           logger,
		ProjectSlice:     nil,
		ReposInProject:   nil,
		TagsInProject:    nil,
		HandlerInterface: nil,
	}, re)
}

func TestParseHarborConfig_Err(t *testing.T) {

	content := `
harbor-repo:
  schema: 
    - http
  address: 192.168.1.1:80           # harbor连接地址
  domain: docker.wl.io              # harbor域
  user: admin                       # harbor用户
  password: 123456                  # harbor用户密码
  preserve-dir: /tmp                # 持久化tag
  projects:                         # 导出哪些项目下的镜像tag（如果为空表示全库导出）
    - aaa                        # project名称
    - bbb
  excludes:                         # 配置'projects'空值使用，过滤某些project
    - ddd
`
	logger := logrus.New()
	_, err := ParseHarborConfig([]byte(content), logger)
	require.NotNil(t, err)
}

func TestProjectNotFoundErr_Error(t *testing.T) {
	err := ProjectNotFoundErr{ProjectName: mockProjectName}
	require.Equal(t, fmt.Sprintf("[not found] 未查询到project -> %s", mockProjectName), err.Error())
}

func TestExecutor_FilterProjects(t *testing.T) {

	mockExecutor.ProjectsToSearchExcludes = []string{"aaa", "bbb"}
	projects := mockExecutor.FilterProjects(mockProjects)
	require.Equal(t, []ProjectInternal{
		{
			Name:      "kubernetes",
			ProjectId: 0,
		},
		{
			Name:      "library",
			ProjectId: 1,
		},
	}, projects)
}

/*
	mock func
*/

func Test_ProjectsByName_Mock(t *testing.T) {

	r := request.HTTPRequestItem{
		Url:      mockApiGetProjectsByNameUrl,
		Method:   http.MethodGet,
		Body:     nil,
		Headers:  nil,
		User:     mockUser,
		Password: mockPassword,
	}

	mockInterface := &mocks.HandlerInterface{}
	mockExecutor.HandlerInterface = mockInterface

	mockInterface.On(MockDoRequestFunc, r).Return(searchProjectBytes, nil).Once()

	p, err := mockExecutor.ProjectsByName(mockProjectName)
	fmt.Println(p)

	var searchResult SearchResult

	json.Unmarshal(searchProjectBytes, &searchResult)

	require.Equal(t, nil, err)
	require.Equal(t, searchResult.Project[0].projectDeepCopy(), p)
}

func Test_ProjectsByNameUnmarshalErr_Mock(t *testing.T) {

	r := request.HTTPRequestItem{
		Url:      mockApiGetProjectsByNameUrl,
		Method:   http.MethodGet,
		Body:     nil,
		Headers:  nil,
		User:     mockUser,
		Password: mockPassword,
	}

	mockInterface := &mocks.HandlerInterface{}
	mockExecutor.HandlerInterface = mockInterface

	mockInterface.On(MockDoRequestFunc, r).Return(nil, nil).Once()

	p, err := mockExecutor.ProjectsByName(mockProjectName)

	var project ProjectInternal
	_, ok := err.(*json.SyntaxError)
	require.Equal(t, true, ok)
	require.Equal(t, project, p)
}

func Test_ProjectsByName_ProjectNotFound_Mock(t *testing.T) {

	r := request.HTTPRequestItem{
		Url:      mockApiGetProjectsByNameUrl,
		Method:   http.MethodGet,
		Body:     nil,
		Headers:  nil,
		User:     mockUser,
		Password: mockPassword,
	}

	mockInterface := &mocks.HandlerInterface{}
	mockExecutor.HandlerInterface = mockInterface

	mockInterface.On(MockDoRequestFunc, r).Return(nil, nil).Once()

	p, err := mockExecutor.ProjectsByName(mockProjectName)

	var project ProjectInternal
	_, ok := err.(*json.SyntaxError)
	require.Equal(t, true, ok)
	require.Equal(t, project, p)
}

func Test_ProjectsByName_NetErr_Mock(t *testing.T) {

	mockInterface := &mocks.HandlerInterface{}
	mockExecutor.HandlerInterface = mockInterface
	r := request.HTTPRequestItem{
		Url:      mockApiGetProjectsByNameUrl,
		Method:   http.MethodGet,
		Body:     nil,
		Headers:  nil,
		User:     mockUser,
		Password: mockPassword,
	}

	mockInterface.On(MockDoRequestFunc, r).Return(nil, mockNetErr).Once()

	p, err := mockExecutor.ProjectsByName(mockProjectName)
	var searchResult SearchResult
	json.Unmarshal(nil, &searchResult)
	require.Equal(t, mockNetErr, err)
	require.Equal(t, ProjectInternal{}, p)
}

func Test_ProjectsByNames_Mock(t *testing.T) {

	r := request.HTTPRequestItem{
		Url:      ApiGetProjectsByName(mockSchema, mockAddress, "busybox"),
		Method:   http.MethodGet,
		Body:     nil,
		Headers:  nil,
		User:     mockUser,
		Password: mockPassword,
	}

	mockInterface := &mocks.HandlerInterface{}
	mockExecutor.HandlerInterface = mockInterface

	mockInterface.On(MockDoRequestFunc, r).Return(busyboxProjectInfoBytes, nil).Once()

	r = request.HTTPRequestItem{
		Url:      ApiGetProjectsByName(mockSchema, mockAddress, "elastic"),
		Method:   http.MethodGet,
		Body:     nil,
		Headers:  nil,
		User:     mockUser,
		Password: mockPassword,
	}

	mockInterface.On(MockDoRequestFunc, r).Return(elasticProjectInfoBytes, nil).Once()

	result, err := mockExecutor.ProjectsByNames([]string{"busybox", "elastic"})
	require.Nil(t, err)
	require.Equal(t, 2, len(result))

}

func Test_ProjectsByNames_Err_Mock(t *testing.T) {

	r := request.HTTPRequestItem{
		Url:      ApiGetProjectsByName(mockSchema, mockAddress, "busybox"),
		Method:   http.MethodGet,
		Body:     nil,
		Headers:  nil,
		User:     mockUser,
		Password: mockPassword,
	}

	mockInterface := &mocks.HandlerInterface{}
	mockExecutor.HandlerInterface = mockInterface
	mockInterface.On(MockDoRequestFunc, r).Return(nil, mockNetErr).Once()

	result, err := mockExecutor.ProjectsByNames([]string{"busybox", "elastic"})
	require.Equal(t, err, mockNetErr)
	require.Nil(t, result)

}

func Test_TagsWithinRepoByPage_Mock(t *testing.T) {

	mockInterface := &mocks.HandlerInterface{}
	mockExecutor.HandlerInterface = mockInterface

	mockInterface.On(MockTagsWithinRepoByPageFunc, mockPage, mockPageSize,
		mockSchema, mockAddress,
		mockProjectName, mockRepoName,
		mockUser, mockPassword, DefaultRequestTimeout).Return(repoTagsByPageOneBytes, nil).Once()

	b, err := mockExecutor.TagsWithinRepoByPage(mockPage, mockPageSize,
		mockSchema, mockAddress,
		mockProjectName, mockRepoName,
		mockUser, mockPassword, DefaultRequestTimeout)

	require.Equal(t, nil, err)
	require.Equal(t, repoTagsByPageOneBytes, b)
}

func Test_TagsWithinRepoByPage_NetErr_Mock(t *testing.T) {

	mockInterface := &mocks.HandlerInterface{}
	mockExecutor.HandlerInterface = mockInterface

	mockInterface.On(MockTagsWithinRepoByPageFunc, mockPage, mockPageSize,
		mockSchema, mockAddress,
		mockProjectName, mockRepoName,
		mockUser, mockPassword, DefaultRequestTimeout).Return(nil, mockNetErr).Once()

	p, err := mockExecutor.TagsWithinRepoByPage(
		mockPage, mockPageSize,
		mockSchema, mockAddress,
		mockProjectName, mockRepoName,
		mockUser, mockPassword, DefaultRequestTimeout,
	)

	require.Equal(t, mockNetErr, err)
	require.Equal(t, 0, len(p))
}

func Test_ProjectsByPage_mock(t *testing.T) {

	mockInterface := &mocks.HandlerInterface{}
	mockExecutor.HandlerInterface = mockInterface

	mockInterface.On("ProjectsByPage",
		mockPage,
		mockPageSize,
		mockLogger,
		mockSchema,
		mockAddress,
		mockUser,
		mockPassword,
		DefaultRequestTimeout,
	).Return(projectsByPageFirstBytes, nil).Once()

	p, err := mockExecutor.ProjectsByPage(
		mockPage,
		mockPageSize,
		mockLogger,
		mockSchema,
		mockAddress,
		mockUser,
		mockPassword,
		DefaultRequestTimeout,
	)

	require.Equal(t, nil, err)
	require.Equal(t, projectsByPageFirstBytes, p)
}

func Test_ProjectsByPage_NetErr_mock(t *testing.T) {

	mockInterface := &mocks.HandlerInterface{}
	mockExecutor.HandlerInterface = mockInterface

	mockInterface.On("ProjectsByPage",
		mockPage,
		mockPageSize,
		mockLogger,
		mockSchema,
		mockAddress,
		mockUser,
		mockPassword,
		DefaultRequestTimeout,
	).Return(nil, mockNetErr).Once()

	p, err := mockExecutor.ProjectsByPage(
		mockPage,
		mockPageSize,
		mockLogger,
		mockSchema,
		mockAddress,
		mockUser,
		mockPassword,
		DefaultRequestTimeout,
	)

	require.Equal(t, mockNetErr, err)
	require.Equal(t, "", string(p))
}

func Test_ProjectList_Mock(t *testing.T) {

	mockInterface := &mocks.HandlerInterface{}
	mockExecutor.HandlerInterface = mockInterface
	mockExecutor.ProjectsToSearch = []string{"elastic"}

	r := request.HTTPRequestItem{
		Url:      ApiGetProjectsByName(mockSchema, mockAddress, "elastic"),
		Method:   http.MethodGet,
		Body:     nil,
		Headers:  nil,
		User:     mockUser,
		Password: mockPassword,
	}
	mockInterface.On(MockDoRequestFunc, r).Return(elasticProjectInfoBytes, nil).Once()
	results, err := mockExecutor.ProjectList()
	require.Nil(t, err)
	require.Equal(t, 1, len(results))
}

func Test_ProjectList_AllCase_Mock(t *testing.T) {

	mockInterface := &mocks.HandlerInterface{}
	mockExecutor.HandlerInterface = mockInterface

	mockExecutor.ProjectsToSearch = []string{}
	// mock -> executor.HandlerInterface.ProjectCount
	mockInterface.On(MockProjectCountFunc, mockUser, mockPassword, mockSchema, mockAddress, DefaultRequestTimeout).Return(statisticsBytes, nil).Once()

	// mock -> AllProjects
	mockInterface.On("ProjectsByPage",
		1,
		mockPageSize,
		mockLogger,
		mockSchema,
		mockAddress,
		mockUser,
		mockPassword,
		DefaultRequestTimeout,
	).Return(projectsByPageFirstBytes, nil).Once()

	mockInterface.On("ProjectsByPage",
		2,
		mockPageSize,
		mockLogger,
		mockSchema,
		mockAddress,
		mockUser,
		mockPassword,
		DefaultRequestTimeout,
	).Return(projectsByPageSecondBytes, nil).Once()

	mockInterface.On("ProjectsByPage",
		3,
		mockPageSize,
		mockLogger,
		mockSchema,
		mockAddress,
		mockUser,
		mockPassword,
		DefaultRequestTimeout,
	).Return(projectsByPageThirdBytes, nil).Once()

	mockInterface.On("ProjectsByPage",
		4,
		mockPageSize,
		mockLogger,
		mockSchema,
		mockAddress,
		mockUser,
		mockPassword,
		DefaultRequestTimeout,
	).Return(projectsByPageFourthBytes, nil).Once()

	mockInterface.On("ProjectsByPage",
		5,
		mockPageSize,
		mockLogger,
		mockSchema,
		mockAddress,
		mockUser,
		mockPassword,
		DefaultRequestTimeout,
	).Return(projectsByPageFifthBytes, nil).Once()

	results, err := mockExecutor.ProjectList()
	require.Nil(t, err)
	require.Equal(t, mockAllProjectsCount, len(results))
}

func Test_ProjectList_AllErrCase_Mock(t *testing.T) {

	mockInterface := &mocks.HandlerInterface{}
	mockExecutor.HandlerInterface = mockInterface

	mockExecutor.ProjectsToSearch = []string{}
	// mock -> executor.HandlerInterface.ProjectCount
	mockInterface.On(MockProjectCountFunc, mockUser, mockPassword, mockSchema, mockAddress, DefaultRequestTimeout).Return(nil, mockNetErr).Once()

	results, err := mockExecutor.ProjectList()
	require.Equal(t, mockNetErr, err)
	require.Equal(t, 0, len(results))
}

func Test_ProjectCount_Mock(t *testing.T) {

	mockInterface := &mocks.HandlerInterface{}
	mockExecutor.HandlerInterface = mockInterface

	mockInterface.On(MockProjectCountFunc, mockUser, mockPassword, mockSchema, mockAddress, DefaultRequestTimeout).Return(statisticsBytes, nil).Once()
	p, err := mockExecutor.ProjectCount(mockUser, mockPassword, mockSchema, mockAddress, DefaultRequestTimeout)

	require.Equal(t, nil, err)
	require.Equal(t, statisticsBytes, p)
}

func Test_ProjectCount_NetErr_Mock(t *testing.T) {

	mockInterface := &mocks.HandlerInterface{}
	mockExecutor.HandlerInterface = mockInterface

	mockInterface.On(MockProjectCountFunc, mockUser, mockPassword, mockSchema, mockAddress, DefaultRequestTimeout).Return(nil, mockNetErr).Once()
	p, err := mockExecutor.ProjectCount(mockUser, mockPassword, mockSchema, mockAddress, DefaultRequestTimeout)

	require.Equal(t, mockNetErr, err)
	require.Equal(t, "", string(p))
}

func Test_AllProjects_Mock(t *testing.T) {

	mockInterface := &mocks.HandlerInterface{}
	mockExecutor.HandlerInterface = mockInterface

	mockInterface.On("ProjectsByPage",
		1,
		mockPageSize,
		mockLogger,
		mockSchema,
		mockAddress,
		mockUser,
		mockPassword,
		DefaultRequestTimeout,
	).Return(projectsByPageFirstBytes, nil).Once()

	mockInterface.On("ProjectsByPage",
		2,
		mockPageSize,
		mockLogger,
		mockSchema,
		mockAddress,
		mockUser,
		mockPassword,
		DefaultRequestTimeout,
	).Return(projectsByPageSecondBytes, nil).Once()

	mockInterface.On("ProjectsByPage",
		3,
		mockPageSize,
		mockLogger,
		mockSchema,
		mockAddress,
		mockUser,
		mockPassword,
		DefaultRequestTimeout,
	).Return(projectsByPageThirdBytes, nil).Once()

	mockInterface.On("ProjectsByPage",
		4,
		mockPageSize,
		mockLogger,
		mockSchema,
		mockAddress,
		mockUser,
		mockPassword,
		DefaultRequestTimeout,
	).Return(projectsByPageFourthBytes, nil).Once()

	mockInterface.On("ProjectsByPage",
		5,
		mockPageSize,
		mockLogger,
		mockSchema,
		mockAddress,
		mockUser,
		mockPassword,
		DefaultRequestTimeout,
	).Return(projectsByPageFifthBytes, nil).Once()

	//mockExecutor.AllProjects(10)
	projects, err := mockExecutor.AllProjects(44)

	require.Equal(t, nil, err)
	require.Equal(t, mockAllProjectsCount, len(projects))
}

func Test_AllProjects_ErrCase_Mock(t *testing.T) {

	mockInterface := &mocks.HandlerInterface{}
	mockExecutor.HandlerInterface = mockInterface

	mockInterface.On("ProjectsByPage",
		1,
		mockPageSize,
		mockLogger,
		mockSchema,
		mockAddress,
		mockUser,
		mockPassword,
		DefaultRequestTimeout,
	).Return(nil, mockNetErr).Once()

	projects, err := mockExecutor.AllProjects(44)

	require.Equal(t, mockNetErr, err)
	require.Equal(t, 0, len(projects))
}

func Test_ReposWithinProject_Mock(t *testing.T) {

	mockInterface := &mocks.HandlerInterface{}
	mockExecutor.HandlerInterface = mockInterface

	mockInterface.On(MockRepoCountFunc,
		mockUser,
		mockPassword,
		mockSchema,
		mockAddress,
		mockProjectId,
		DefaultRequestTimeout,
	).Return(kubernetesMetaInfoBytes, nil).Once()

	mockInterface.On(MockListRepoByPageFunc,
		1, 10,
		mockUser, mockPassword,
		mockSchema, mockAddress,
		mockProjectName,
		DefaultRequestTimeout,
	).Return(kubernetesReposByPage1Bytes, nil).Once()

	mockInterface.On(MockListRepoByPageFunc,
		2, 10,
		mockUser, mockPassword,
		mockSchema, mockAddress,
		mockProjectName,
		DefaultRequestTimeout,
	).Return(kubernetesReposByPage2Bytes, nil).Once()

	results, err := mockExecutor.ReposWithinProject(mockProjectName, mockProjectId)
	require.Nil(t, err)
	require.Equal(t, mockProjectRepoCount, len(results))
}

func Test_ReposWithinProject_Err_Mock(t *testing.T) {

	mockInterface := &mocks.HandlerInterface{}
	mockExecutor.HandlerInterface = mockInterface

	// MockListRepoByPageFunc -> mock err
	mockInterface.On(MockRepoCountFunc,
		mockUser,
		mockPassword,
		mockSchema,
		mockAddress,
		mockProjectId,
		DefaultRequestTimeout,
	).Return(kubernetesMetaInfoBytes, mockNetErr).Once()

	mockInterface.On(MockListRepoByPageFunc,
		1, 10,
		mockUser, mockPassword,
		mockSchema, mockAddress,
		mockProjectName,
		DefaultRequestTimeout,
	).Return(kubernetesReposByPage1Bytes, nil).Once()

	mockInterface.On(MockListRepoByPageFunc,
		2, 10,
		mockUser, mockPassword,
		mockSchema, mockAddress,
		mockProjectName,
		DefaultRequestTimeout,
	).Return(kubernetesReposByPage2Bytes, nil).Once()

	results, err := mockExecutor.ReposWithinProject(mockProjectName, mockProjectId)
	require.Equal(t, mockNetErr, err)
	require.Equal(t, 0, len(results))
}

func Test_ReposWithinProject_Err2_Mock(t *testing.T) {

	mockInterface := &mocks.HandlerInterface{}
	mockExecutor.HandlerInterface = mockInterface

	mockInterface.On(MockRepoCountFunc,
		mockUser,
		mockPassword,
		mockSchema,
		mockAddress,
		mockProjectId,
		DefaultRequestTimeout,
	).Return(kubernetesMetaInfoBytes, nil).Once()

	mockInterface.On(MockListRepoByPageFunc,
		1, 10,
		mockUser, mockPassword,
		mockSchema, mockAddress,
		mockProjectName,
		DefaultRequestTimeout,
	).Return(nil, mockNetErr).Once()

	mockInterface.On(MockListRepoByPageFunc,
		2, 10,
		mockUser, mockPassword,
		mockSchema, mockAddress,
		mockProjectName,
		DefaultRequestTimeout,
	).Return(nil, mockNetErr).Once()

	results, err := mockExecutor.ReposWithinProject(mockProjectName, mockProjectId)
	require.NotNil(t, err)
	require.Equal(t, 0, len(results))
}

func Test_TagsWithinRepo_Mock(t *testing.T) {
	mockInterface := &mocks.HandlerInterface{}
	mockExecutor.HandlerInterface = mockInterface

	mockInterface.On(MockTagsWithinRepoByPageFunc, 1, mockPageSize,
		mockSchema, mockAddress,
		mockProjectName, mockRepoName,
		mockUser, mockPassword, DefaultRequestTimeout).Return(repoTagsByPageOneBytes, nil).Once()

	mockInterface.On(MockTagsWithinRepoByPageFunc, 2, mockPageSize,
		mockSchema, mockAddress,
		mockProjectName, mockRepoName,
		mockUser, mockPassword, DefaultRequestTimeout).Return(repoTagsByPageTwoBytes, nil).Once()

	// kubernetes/cephcsi -> mocks/{repoTagsByPage1.json,repoTagsByPage2.json}
	results, err := mockExecutor.TagsWithinRepo(mockProjectName, mockRepoName, 12)
	require.Nil(t, err)
	require.Equal(t, 12, len(results))
}

func Test_TagsWithinRepo_Err_Mock(t *testing.T) {
	mockInterface := &mocks.HandlerInterface{}
	mockExecutor.HandlerInterface = mockInterface

	mockInterface.On(MockTagsWithinRepoByPageFunc, 1, mockPageSize,
		mockSchema, mockAddress,
		mockProjectName, mockRepoName,
		mockUser, mockPassword, DefaultRequestTimeout).Return(repoTagsByPageOneBytes, mockNetErr).Once()

	mockInterface.On(MockTagsWithinRepoByPageFunc, 2, mockPageSize,
		mockSchema, mockAddress,
		mockProjectName, mockRepoName,
		mockUser, mockPassword, DefaultRequestTimeout).Return(repoTagsByPageTwoBytes, mockNetErr).Once()

	// kubernetes/cephcsi -> mocks/{repoTagsByPage1.json,repoTagsByPage2.json}
	results, err := mockExecutor.TagsWithinRepo(mockProjectName, mockRepoName, 12)
	require.NotNil(t, err)
	require.Equal(t, 0, len(results))
}

func Test_ReposWithinProjects_Mock(t *testing.T) {
	mockInterface := &mocks.HandlerInterface{}
	mockExecutor.HandlerInterface = mockInterface

	projects := []ProjectInternal{
		{
			Name:      "kubernetes",
			ProjectId: 5,
		},
		{
			Name:      "elastic",
			ProjectId: 38,
		},
	}

	// mock kubernetes
	mockInterface.On(MockRepoCountFunc,
		mockUser,
		mockPassword,
		mockSchema,
		mockAddress,
		5,
		DefaultRequestTimeout,
	).Return(kubernetesMetaInfoBytes, nil).Once()

	mockInterface.On(MockListRepoByPageFunc,
		1, 10,
		mockUser, mockPassword,
		mockSchema, mockAddress,
		"kubernetes",
		DefaultRequestTimeout,
	).Return(kubernetesReposByPage1Bytes, nil).Once()
	//
	mockInterface.On(MockListRepoByPageFunc,
		2, 10,
		mockUser, mockPassword,
		mockSchema, mockAddress,
		"kubernetes",
		DefaultRequestTimeout,
	).Return(kubernetesReposByPage2Bytes, nil).Once()

	// mock elastic
	mockInterface.On(MockRepoCountFunc,
		mockUser,
		mockPassword,
		mockSchema,
		mockAddress,
		38,
		DefaultRequestTimeout,
	).Return(elasticMetaInfoBytes, nil).Once()

	mockInterface.On(MockListRepoByPageFunc,
		1, 10,
		mockUser, mockPassword,
		mockSchema, mockAddress,
		"elastic",
		DefaultRequestTimeout,
	).Return(elasticReposByPage1Bytes, nil).Once()

	// elasticMetaInfoBytes

	results, err := mockExecutor.ReposWithinProjects(projects)
	require.Nil(t, err)
	require.Equal(t, 2, len(results["elastic"]))
	require.Equal(t, 15, len(results["kubernetes"]))

	fmt.Println(results)
}

func Test_ReposWithinProjects_Err_Mock(t *testing.T) {
	mockInterface := &mocks.HandlerInterface{}
	mockExecutor.HandlerInterface = mockInterface

	projects := []ProjectInternal{
		{
			Name:      "kubernetes",
			ProjectId: 5,
		},
		{
			Name:      "elastic",
			ProjectId: 38,
		},
	}

	// mock kubernetes
	mockInterface.On(MockRepoCountFunc,
		mockUser,
		mockPassword,
		mockSchema,
		mockAddress,
		5,
		DefaultRequestTimeout,
	).Return(nil, mockNetErr).Once()

	results, err := mockExecutor.ReposWithinProjects(projects)
	require.NotNil(t, err)
	require.Equal(t, 0, len(results["elastic"]))
	require.Equal(t, 0, len(results["kubernetes"]))
}

func Test_generateRepoTagsSlice_Mock(t *testing.T) {
	mockInterface := &mocks.HandlerInterface{}
	mockExecutor.HandlerInterface = mockInterface

	mockInterface.On("TagsNumWithRepo",
		mockUser, mockPassword, mockSchema, mockAddress, mockProjectName, mockRepoName, DefaultRequestTimeout,
	).Return(cephcsiRepoTagsNumBytes, nil).Once()

	mockInterface.On("TagsWithinRepoByPage",
		mockPage, mockPageSize,
		mockSchema, mockAddress,
		mockProjectName, mockRepoName,
		mockUser, mockPassword,
		DefaultRequestTimeout,
	).Return(cephcsiRepoTagsBytes, nil)

	re, err := mockExecutor.generateRepoTagsSlice(mockProjectName, mockRepoName)
	require.Nil(t, err)
	require.Equal(t, []string{
		"kubernetes/cephcsi:v3.2.0",
		"kubernetes/cephcsi:v3.2.1",
	}, re)
}

func Test_generateRepoTagsSlice_Err_Mock(t *testing.T) {
	mockInterface := &mocks.HandlerInterface{}
	mockExecutor.HandlerInterface = mockInterface

	mockInterface.On("TagsNumWithRepo",
		mockUser, mockPassword, mockSchema, mockAddress, mockProjectName, mockRepoName, DefaultRequestTimeout,
	).Return(cephcsiRepoTagsNumBytes, mockNetErr).Once()

	re, err := mockExecutor.generateRepoTagsSlice(mockProjectName, mockRepoName)
	require.Equal(t, mockNetErr, err)
	var s []string
	require.Equal(t, s, re)

	mockInterface2 := &mocks.HandlerInterface{}
	mockExecutor.HandlerInterface = mockInterface2

	mockInterface2.On("TagsNumWithRepo",
		mockUser, mockPassword, mockSchema, mockAddress, mockProjectName, mockRepoName, DefaultRequestTimeout,
	).Return(cephcsiRepoTagsNumBytes, nil).Once()

	mockInterface2.On("TagsWithinRepoByPage",
		mockPage, mockPageSize,
		mockSchema, mockAddress,
		mockProjectName, mockRepoName,
		mockUser, mockPassword, DefaultRequestTimeout,
	).Return(nil, mockNetErr)

	re2, err2 := mockExecutor.generateRepoTagsSlice(mockProjectName, mockRepoName)
	require.Equal(t, mockNetErr, err2)
	require.Equal(t, s, re2)
}

func Test_TagsWithinProject_Mock(t *testing.T) {

	kubernetesRepos := []string{"cephcsi", "stakater/reloader"}
	mockInterface := &mocks.HandlerInterface{}
	mockExecutor.HandlerInterface = mockInterface

	mockInterface.On("TagsNumWithRepo",
		mockUser, mockPassword, mockSchema, mockAddress, mockProjectName, "cephcsi", DefaultRequestTimeout,
	).Return(cephcsiRepoTagsNumBytes, nil).Once()

	mockInterface.On("TagsNumWithRepo",
		mockUser, mockPassword, mockSchema, mockAddress, mockProjectName, "stakater/reloader", DefaultRequestTimeout,
	).Return(reloaderRepoTagsNumBytes, nil).Once()

	mockInterface.On("TagsWithinRepoByPage",
		mockPage, mockPageSize,
		mockSchema, mockAddress,
		mockProjectName, mockRepoName,
		mockUser, mockPassword, DefaultRequestTimeout,
	).Return(cephcsiRepoTagsBytes, nil)

	mockInterface.On("TagsWithinRepoByPage",
		mockPage, mockPageSize,
		mockSchema, mockAddress,
		mockProjectName, "stakater/reloader",
		mockUser, mockPassword, DefaultRequestTimeout,
	).Return(reloaderRepoTagsBytes, nil)

	re, err := mockExecutor.TagsWithinProject(mockProjectName, kubernetesRepos)

	sort.Strings(re.Tags)

	require.Nil(t, err)
	require.Equal(t, projectMapTags{
		ProjectName: mockProjectName,
		Tags: []string{
			"kubernetes/cephcsi:v3.2.0",
			"kubernetes/cephcsi:v3.2.1",
			"kubernetes/stakater/reloader:v0.0.97",
		},
	}, re)
}

func Test_TagsWithinProject_Error_Mock(t *testing.T) {

	// 1. mock http error
	kubernetesRepos := []string{"cephcsi", "stakater/reloader"}
	mockInterface := &mocks.HandlerInterface{}
	mockExecutor.HandlerInterface = mockInterface

	// 协程情况下，返回顺序不确定，因此两个函数均Mock error
	mockInterface.On("TagsNumWithRepo",
		mockUser, mockPassword, mockSchema, mockAddress, mockProjectName, "cephcsi", DefaultRequestTimeout,
	).Return(nil, mockNetErr).Once()

	mockInterface.On("TagsNumWithRepo",
		mockUser, mockPassword, mockSchema, mockAddress, mockProjectName, "stakater/reloader", DefaultRequestTimeout,
	).Return(nil, mockNetErr).Once()

	re, err := mockExecutor.TagsWithinProject(mockProjectName, kubernetesRepos)

	sort.Strings(re.Tags)
	require.NotNil(t, err)
	require.Equal(t, projectMapTags{}, re)

	// 2. mock tagsCount == 0
	mockInterface2 := &mocks.HandlerInterface{}
	mockExecutor.HandlerInterface = mockInterface2

	// 协程情况下，返回顺序不确定，因此两个函数均Mock error
	mockInterface2.On("TagsNumWithRepo",
		mockUser, mockPassword, mockSchema, mockAddress, mockProjectName, "cephcsi", DefaultRequestTimeout,
	).Return([]byte("ddd"), nil).Once()

	mockInterface2.On("TagsNumWithRepo",
		mockUser, mockPassword, mockSchema, mockAddress, mockProjectName, "stakater/reloader", DefaultRequestTimeout,
	).Return([]byte("ddd"), nil).Once()

	re2, err2 := mockExecutor.TagsWithinProject(mockProjectName, kubernetesRepos)

	sort.Strings(re2.Tags)
	require.NotNil(t, err2)
	require.Equal(t, projectMapTags{}, re)
}

func Test_TagsWithinProjects_Mock(t *testing.T) {

	kubernetesRepos := []string{"cephcsi", "stakater/reloader"}
	apacheRepos := []string{"skywalking-ui", "skywalking-oap-server", "skywalking-java-agent"}

	mockInterface := &mocks.HandlerInterface{}
	mockExecutor.HandlerInterface = mockInterface

	mockInterface.On("TagsNumWithRepo",
		mockUser, mockPassword, mockSchema, mockAddress, mockProjectName, "cephcsi", DefaultRequestTimeout,
	).Return(cephcsiRepoTagsNumBytes, nil).Once()

	mockInterface.On("TagsNumWithRepo",
		mockUser, mockPassword, mockSchema, mockAddress, mockProjectName, "stakater/reloader", DefaultRequestTimeout,
	).Return(reloaderRepoTagsNumBytes, nil).Once()

	mockInterface.On("TagsWithinRepoByPage",
		mockPage, mockPageSize,
		mockSchema, mockAddress,
		mockProjectName, mockRepoName,
		mockUser, mockPassword, DefaultRequestTimeout,
	).Return(cephcsiRepoTagsBytes, nil)

	mockInterface.On("TagsWithinRepoByPage",
		mockPage, mockPageSize,
		mockSchema, mockAddress,
		mockProjectName, "stakater/reloader",
		mockUser, mockPassword, DefaultRequestTimeout,
	).Return(reloaderRepoTagsBytes, nil)

	mockInterface.On("TagsNumWithRepo",
		mockUser, mockPassword, mockSchema, mockAddress, "apache", "skywalking-ui", DefaultRequestTimeout,
	).Return(skywalkingUITagsNumBytes, nil).Once()

	mockInterface.On("TagsNumWithRepo",
		mockUser, mockPassword, mockSchema, mockAddress, "apache", "skywalking-oap-server", DefaultRequestTimeout,
	).Return(skywalkingServerTagsNumBytes, nil).Once()

	mockInterface.On("TagsNumWithRepo",
		mockUser, mockPassword, mockSchema, mockAddress, "apache", "skywalking-java-agent", DefaultRequestTimeout,
	).Return(skywalkingAgentTagsNumBytes, nil).Once()

	mockInterface.On("TagsWithinRepoByPage",
		mockPage, mockPageSize,
		mockSchema, mockAddress,
		"apache", "skywalking-ui",
		mockUser, mockPassword, DefaultRequestTimeout,
	).Return(skywalkingUITagsBytes, nil)

	mockInterface.On("TagsWithinRepoByPage",
		mockPage, mockPageSize,
		mockSchema, mockAddress,
		"apache", "skywalking-oap-server",
		mockUser, mockPassword, DefaultRequestTimeout,
	).Return(skywalkingServerTagsBytes, nil)

	mockInterface.On("TagsWithinRepoByPage",
		mockPage, mockPageSize,
		mockSchema, mockAddress,
		"apache", "skywalking-java-agent",
		mockUser, mockPassword, DefaultRequestTimeout,
	).Return(skywalkingAgentTagsBytes, nil)

	re, err := mockExecutor.TagsWithinProjects(map[string][]string{
		"apache":     apacheRepos,
		"kubernetes": kubernetesRepos,
	})

	require.Nil(t, err)

	expect := make(map[string][]string)
	expect["apache"] = []string{
		"apache/skywalking-java-agent:8.6.0-alpine",
		"apache/skywalking-ui:8.6.0",
		"apache/skywalking-oap-server:8.6.0-es7",
	}
	expect["kubernetes"] = []string{
		"kubernetes/cephcsi:v3.2.0",
		"kubernetes/cephcsi:v3.2.1",
		"kubernetes/stakater/reloader:v0.0.97",
	}

	sort.Strings(expect["apache"])
	sort.Strings(re["apache"])
	require.Equal(t, expect["apache"], re["apache"])

	sort.Strings(expect["kubernetes"])
	sort.Strings(re["kubernetes"])
	require.Equal(t, expect["kubernetes"], re["kubernetes"])
}

func Test_TagsWithinProjects_Err_Mock(t *testing.T) {

	kubernetesRepos := []string{"cephcsi", "stakater/reloader"}
	apacheRepos := []string{"skywalking-ui", "skywalking-oap-server", "skywalking-java-agent"}

	mockInterface := &mocks.HandlerInterface{}
	mockExecutor.HandlerInterface = mockInterface

	mockInterface.On("TagsNumWithRepo",
		mockUser, mockPassword, mockSchema, mockAddress, mockProjectName, "cephcsi", DefaultRequestTimeout,
	).Return(nil, mockNetErr).Once()

	mockInterface.On("TagsNumWithRepo",
		mockUser, mockPassword, mockSchema, mockAddress, mockProjectName, "stakater/reloader", DefaultRequestTimeout,
	).Return(nil, mockNetErr).Once()

	mockInterface.On("TagsNumWithRepo",
		mockUser, mockPassword, mockSchema, mockAddress, "apache", "skywalking-ui", DefaultRequestTimeout,
	).Return(nil, mockNetErr).Once()

	mockInterface.On("TagsNumWithRepo",
		mockUser, mockPassword, mockSchema, mockAddress, "apache", "skywalking-oap-server", DefaultRequestTimeout,
	).Return(nil, mockNetErr).Once()

	mockInterface.On("TagsNumWithRepo",
		mockUser, mockPassword, mockSchema, mockAddress, "apache", "skywalking-java-agent", DefaultRequestTimeout,
	).Return(nil, mockNetErr).Once()

	_, err := mockExecutor.TagsWithinProjects(map[string][]string{
		"apache":     apacheRepos,
		"kubernetes": kubernetesRepos,
	})

	require.NotNil(t, err)
}

func Test_GenerateImageList_Mock(t *testing.T) {

	expect := make(map[string][]string)
	expect["apache"] = []string{
		"apache/skywalking-java-agent:8.6.0-alpine",
		"apache/skywalking-ui:8.6.0",
		"apache/skywalking-oap-server:8.6.0-es7",
	}
	expect["kubernetes"] = []string{
		"kubernetes/cephcsi:v3.2.0",
		"kubernetes/cephcsi:v3.2.1",
		"kubernetes/stakater/reloader:v0.0.97",
	}

	mockExecutor.TagWithDomain = true
	mockExecutor.PreserveDir = "ddd"
	err := mockExecutor.GenerateImageList(expect)

	require.Nil(t, err)

	err = os.RemoveAll("ddd")
	if err != nil {
		panic(err)
	}
}

func Test_GenerateImageList_BadPreserveDirErr_Mock(t *testing.T) {

	expect := make(map[string][]string)
	expect["apache"] = []string{
		"apache/skywalking-java-agent:8.6.0-alpine",
		"apache/skywalking-ui:8.6.0",
		"apache/skywalking-oap-server:8.6.0-es7",
	}
	expect["kubernetes"] = []string{
		"kubernetes/cephcsi:v3.2.0",
		"kubernetes/cephcsi:v3.2.1",
		"kubernetes/stakater/reloader:v0.0.97",
	}

	mockExecutor.TagWithDomain = true
	switch runtime.GOOS {
	case "windows":
		mockExecutor.PreserveDir = "/root/ddd"
	case "linux":
		mockExecutor.PreserveDir = "c://tmp"
	default:
		mockExecutor.PreserveDir = "/root/ddd"
	}
	err := mockExecutor.GenerateImageList(expect)

	require.NotNil(t, err)
}

func Test_mkDirIfNotExist(t *testing.T) {
	os.Mkdir("ooo", 0644)
	err := mkDirIfNotExist("ooo", 0644)
	require.Nil(t, err)
	os.RemoveAll("ooo")

	err2 := mkDirIfNotExist("asd", 0644)
	require.Nil(t, err2)
	os.RemoveAll("asd")
}

func Test_ImageList_Mock(t *testing.T) {

	defer os.RemoveAll("harbor-image-list")

	content := `
harbor-repo:
  schema: http
  address: 192.168.1.1:80           # harbor连接地址
  domain: harbor.wl.io              # harbor域
  user: admin                       # harbor用户
  password: 123456                  # harbor用户密码
  preserve-dir: harbor-image-list   # 持久化tag
  withDomain: true                  # 镜像tag是否包含harbor domain (harbor.wl.io/library/busybox:latest,即含有harbor.wl.io)
  projects:                         # 导出哪些项目下的镜像tag（如果为空表示全库导出）
    - apache                        # project名称
    - elastic
  excludes:                         # 配置'projects'空值使用，过滤某些project
    - ddd
`
	mockInterface := &mocks.HandlerInterface{}

	// 1.mock list projects
	r1 := request.HTTPRequestItem{
		Url:      ApiGetProjectsByName(mockSchema, mockAddress, "apache"),
		Method:   http.MethodGet,
		Body:     nil,
		Headers:  nil,
		User:     mockUser,
		Password: mockPassword,
	}
	mockInterface.On(MockDoRequestFunc, r1).Return(apacheProjectInfoBytes, nil).Once()

	r2 := request.HTTPRequestItem{
		Url:      ApiGetProjectsByName(mockSchema, mockAddress, "elastic"),
		Method:   http.MethodGet,
		Body:     nil,
		Headers:  nil,
		User:     mockUser,
		Password: mockPassword,
	}
	mockInterface.On(MockDoRequestFunc, r2).Return(elasticProjectInfoBytes, nil).Once()

	// 2.1 mock `apache` repo num && repos by page
	mockInterface.On(MockRepoCountFunc,
		mockUser,
		mockPassword,
		mockSchema,
		mockAddress,
		64,
		DefaultRequestTimeout,
	).Return(apacheProjectMetaInfoBytes, nil).Once()

	mockInterface.On(MockListRepoByPageFunc,
		1, 10,
		mockUser, mockPassword,
		mockSchema, mockAddress,
		"apache",
		DefaultRequestTimeout,
	).Return(apacheReposByPageBytes, nil).Once()

	// 2.1 mock `elastic` repo num && repos by page
	mockInterface.On(MockRepoCountFunc,
		mockUser,
		mockPassword,
		mockSchema,
		mockAddress,
		38,
		DefaultRequestTimeout,
	).Return(elasticMetaInfoBytes, nil).Once()

	mockInterface.On(MockListRepoByPageFunc,
		1, 10,
		mockUser, mockPassword,
		mockSchema, mockAddress,
		"elastic",
		DefaultRequestTimeout,
	).Return(elasticReposByPage1Bytes, nil).Once()

	mockInterface.On(MockRepoCountFunc,
		mockUser,
		mockPassword,
		mockSchema,
		mockAddress,
		38,
		DefaultRequestTimeout,
	).Return(elasticMetaInfoBytes, nil).Once()

	// 3.1 mock elastic's repo tags num

	mockInterface.On("TagsNumWithRepo",
		mockUser, mockPassword, mockSchema, mockAddress, "elastic", "filebeat", DefaultRequestTimeout,
	).Return(elasticFilebeatTagsNumBytes, nil).Once()

	mockInterface.On("TagsNumWithRepo",
		mockUser, mockPassword, mockSchema, mockAddress, "elastic", "elasticsearch", DefaultRequestTimeout,
	).Return(elasticElasticsearchTagsNumBytes, nil).Once()

	// 3.2 mock apache's repo tags num

	mockInterface.On("TagsNumWithRepo",
		mockUser, mockPassword, mockSchema, mockAddress, "apache", "skywalking-ui", DefaultRequestTimeout,
	).Return(skywalkingUITagsNumBytes, nil).Once()

	mockInterface.On("TagsNumWithRepo",
		mockUser, mockPassword, mockSchema, mockAddress, "apache", "skywalking-oap-server", DefaultRequestTimeout,
	).Return(skywalkingServerTagsNumBytes, nil).Once()

	mockInterface.On("TagsNumWithRepo",
		mockUser, mockPassword, mockSchema, mockAddress, "apache", "skywalking-java-agent", DefaultRequestTimeout,
	).Return(skywalkingAgentTagsNumBytes, nil).Once()

	// 4.1 mock elastic's repo tags by page
	mockInterface.On("TagsWithinRepoByPage",
		mockPage, mockPageSize,
		mockSchema, mockAddress,
		"apache", "skywalking-ui",
		mockUser, mockPassword, DefaultRequestTimeout,
	).Return(skywalkingUITagsBytes, nil)

	mockInterface.On("TagsWithinRepoByPage",
		mockPage, mockPageSize,
		mockSchema, mockAddress,
		"apache", "skywalking-oap-server",
		mockUser, mockPassword, DefaultRequestTimeout,
	).Return(skywalkingServerTagsBytes, nil)

	mockInterface.On("TagsWithinRepoByPage",
		mockPage, mockPageSize,
		mockSchema, mockAddress,
		"apache", "skywalking-java-agent",
		mockUser, mockPassword, DefaultRequestTimeout,
	).Return(skywalkingAgentTagsBytes, nil)

	// 4.2 mock elastic's repo tags by page
	mockInterface.On("TagsWithinRepoByPage",
		mockPage, mockPageSize,
		mockSchema, mockAddress,
		"elastic", "filebeat",
		mockUser, mockPassword, DefaultRequestTimeout,
	).Return(elasticFilebeatTagsByPageBytes, nil)

	mockInterface.On("TagsWithinRepoByPage",
		mockPage, mockPageSize,
		mockSchema, mockAddress,
		"elastic", "elasticsearch",
		mockUser, mockPassword, DefaultRequestTimeout,
	).Return(elasticElasticsearchTagsByPageBytes, nil)

	ImageList(command.OperationItem{
		B:         []byte(content),
		Logger:    mockLogger,
		Local:     false,
		Interface: mockInterface,
	})
}

func Test_ImageList_ParseConfigErr_Mock(t *testing.T) {
	content := `
harbor-repo:
  - schema: http
  address: 192.168.1.1:80           # harbor连接地址
  domain: harbor.wl.io              # harbor域
  user: admin                       # harbor用户
  password: 123456                  # harbor用户密码
  preserve-dir: harbor-image-list   # 持久化tag
  withDomain: true                  # 镜像tag是否包含harbor domain (harbor.wl.io/library/busybox:latest,即含有harbor.wl.io)
  projects:                         # 导出哪些项目下的镜像tag（如果为空表示全库导出）
    - apache                        # project名称
    - elastic
  excludes:                         # 配置'projects'空值使用，过滤某些project
    - ddd
`
	mockInterface := &mocks.HandlerInterface{}

	err := ImageList(command.OperationItem{
		B:         []byte(content),
		Logger:    mockLogger,
		Local:     false,
		Interface: mockInterface,
	})

	require.NotNil(t, err)
}

func Test_ImageList_GetProjectListErr_Mock(t *testing.T) {
	content := `
harbor-repo:
  schema: http
  address: 192.168.1.1:80           # harbor连接地址
  domain: harbor.wl.io              # harbor域
  user: admin                       # harbor用户
  password: 123456                  # harbor用户密码
  preserve-dir: harbor-image-list   # 持久化tag
  withDomain: true                  # 镜像tag是否包含harbor domain (harbor.wl.io/library/busybox:latest,即含有harbor.wl.io)
  projects:                         # 导出哪些项目下的镜像tag（如果为空表示全库导出）
    - apache                        # project名称
    - elastic
  excludes:                         # 配置'projects'空值使用，过滤某些project
    - ddd
`
	mockInterface := &mocks.HandlerInterface{}

	r1 := request.HTTPRequestItem{
		Url:      ApiGetProjectsByName(mockSchema, mockAddress, "apache"),
		Method:   http.MethodGet,
		Body:     nil,
		Headers:  nil,
		User:     mockUser,
		Password: mockPassword,
	}
	mockInterface.On(MockDoRequestFunc, r1).Return(nil, mockNetErr).Once()

	r2 := request.HTTPRequestItem{
		Url:      ApiGetProjectsByName(mockSchema, mockAddress, "elastic"),
		Method:   http.MethodGet,
		Body:     nil,
		Headers:  nil,
		User:     mockUser,
		Password: mockPassword,
	}
	mockInterface.On(MockDoRequestFunc, r2).Return(nil, mockNetErr).Once()

	err := ImageList(command.OperationItem{
		B:         []byte(content),
		Logger:    mockLogger,
		Local:     false,
		Interface: mockInterface,
	})

	require.NotNil(t, err)
}

func Test_ImageList_GetReposWithinProjectsErr_Mock(t *testing.T) {
	content := `
harbor-repo:
  schema: http
  address: 192.168.1.1:80           # harbor连接地址
  domain: harbor.wl.io              # harbor域
  user: admin                       # harbor用户
  password: 123456                  # harbor用户密码
  preserve-dir: harbor-image-list   # 持久化tag
  withDomain: true                  # 镜像tag是否包含harbor domain (harbor.wl.io/library/busybox:latest,即含有harbor.wl.io)
  projects:                         # 导出哪些项目下的镜像tag（如果为空表示全库导出）
    - apache                        # project名称
    - elastic
  excludes:                         # 配置'projects'空值使用，过滤某些project
    - ddd
`
	mockInterface := &mocks.HandlerInterface{}

	// 1.mock list projects
	r1 := request.HTTPRequestItem{
		Url:      ApiGetProjectsByName(mockSchema, mockAddress, "apache"),
		Method:   http.MethodGet,
		Body:     nil,
		Headers:  nil,
		User:     mockUser,
		Password: mockPassword,
	}
	mockInterface.On(MockDoRequestFunc, r1).Return(apacheProjectInfoBytes, nil).Once()

	r2 := request.HTTPRequestItem{
		Url:      ApiGetProjectsByName(mockSchema, mockAddress, "elastic"),
		Method:   http.MethodGet,
		Body:     nil,
		Headers:  nil,
		User:     mockUser,
		Password: mockPassword,
	}
	mockInterface.On(MockDoRequestFunc, r2).Return(elasticProjectInfoBytes, nil).Once()

	// 2.1 mock `apache` repo num && repos by page
	mockInterface.On(MockRepoCountFunc,
		mockUser,
		mockPassword,
		mockSchema,
		mockAddress,
		64, DefaultRequestTimeout,
	).Return(apacheProjectMetaInfoBytes, nil).Once()

	mockInterface.On(MockListRepoByPageFunc,
		1, 10,
		mockUser, mockPassword,
		mockSchema, mockAddress,
		"apache", DefaultRequestTimeout,
	).Return(apacheReposByPageBytes, nil).Once()

	// 2.1 mock `elastic` repo num && repos by page
	mockInterface.On(MockRepoCountFunc,
		mockUser,
		mockPassword,
		mockSchema,
		mockAddress,
		38, DefaultRequestTimeout,
	).Return(nil, mockNetErr).Once()

	mockInterface.On(MockListRepoByPageFunc,
		1, 10,
		mockUser, mockPassword,
		mockSchema, mockAddress,
		"elastic", DefaultRequestTimeout,
	).Return(elasticReposByPage1Bytes, mockNetErr).Once()

	mockInterface.On(MockRepoCountFunc,
		mockUser,
		mockPassword,
		mockSchema,
		mockAddress,
		38,
		DefaultRequestTimeout,
	).Return(elasticMetaInfoBytes, mockNetErr).Once()

	err := ImageList(command.OperationItem{
		B:         []byte(content),
		Logger:    mockLogger,
		Local:     false,
		Interface: mockInterface,
	})

	require.NotNil(t, err)
}

func Test_ImageList_GetTagsWithinProjectsErr_Mock(t *testing.T) {
	content := `
harbor-repo:
  schema: http
  address: 192.168.1.1:80           # harbor连接地址
  domain: harbor.wl.io              # harbor域
  user: admin                       # harbor用户
  password: 123456                  # harbor用户密码
  preserve-dir: harbor-image-list   # 持久化tag
  withDomain: true                  # 镜像tag是否包含harbor domain (harbor.wl.io/library/busybox:latest,即含有harbor.wl.io)
  projects:                         # 导出哪些项目下的镜像tag（如果为空表示全库导出）
    - apache                        # project名称
    - elastic
  excludes:                         # 配置'projects'空值使用，过滤某些project
    - ddd
`
	mockInterface := &mocks.HandlerInterface{}

	// 1.mock list projects
	r1 := request.HTTPRequestItem{
		Url:      ApiGetProjectsByName(mockSchema, mockAddress, "apache"),
		Method:   http.MethodGet,
		Body:     nil,
		Headers:  nil,
		User:     mockUser,
		Password: mockPassword,
	}
	mockInterface.On(MockDoRequestFunc, r1).Return(apacheProjectInfoBytes, nil).Once()

	r2 := request.HTTPRequestItem{
		Url:      ApiGetProjectsByName(mockSchema, mockAddress, "elastic"),
		Method:   http.MethodGet,
		Body:     nil,
		Headers:  nil,
		User:     mockUser,
		Password: mockPassword,
	}
	mockInterface.On(MockDoRequestFunc, r2).Return(elasticProjectInfoBytes, nil).Once()

	// 2.1 mock `apache` repo num && repos by page
	mockInterface.On(MockRepoCountFunc,
		mockUser,
		mockPassword,
		mockSchema,
		mockAddress,
		64, DefaultRequestTimeout,
	).Return(apacheProjectMetaInfoBytes, nil).Once()

	mockInterface.On(MockListRepoByPageFunc,
		1, 10,
		mockUser, mockPassword,
		mockSchema, mockAddress,
		"apache", DefaultRequestTimeout,
	).Return(apacheReposByPageBytes, nil).Once()

	// 2.1 mock `elastic` repo num && repos by page
	mockInterface.On(MockRepoCountFunc,
		mockUser,
		mockPassword,
		mockSchema,
		mockAddress,
		38, DefaultRequestTimeout,
	).Return(elasticMetaInfoBytes, nil).Once()

	mockInterface.On(MockListRepoByPageFunc,
		1, 10,
		mockUser, mockPassword,
		mockSchema, mockAddress,
		"elastic", DefaultRequestTimeout,
	).Return(elasticReposByPage1Bytes, nil).Once()

	mockInterface.On(MockRepoCountFunc,
		mockUser,
		mockPassword,
		mockSchema,
		mockAddress,
		38, DefaultRequestTimeout,
	).Return(elasticMetaInfoBytes, nil).Once()

	// 3.1 mock elastic's repo tags num

	mockInterface.On("TagsNumWithRepo",
		mockUser, mockPassword, mockSchema, mockAddress, "elastic", "filebeat", DefaultRequestTimeout,
	).Return(elasticFilebeatTagsNumBytes, nil).Once()

	mockInterface.On("TagsNumWithRepo",
		mockUser, mockPassword, mockSchema, mockAddress, "elastic", "elasticsearch", DefaultRequestTimeout,
	).Return(elasticElasticsearchTagsNumBytes, nil).Once()

	// 3.2 mock apache's repo tags num

	mockInterface.On("TagsNumWithRepo",
		mockUser, mockPassword, mockSchema, mockAddress, "apache", "skywalking-ui", DefaultRequestTimeout,
	).Return(skywalkingUITagsNumBytes, nil).Once()

	mockInterface.On("TagsNumWithRepo",
		mockUser, mockPassword, mockSchema, mockAddress, "apache", "skywalking-oap-server", DefaultRequestTimeout,
	).Return(skywalkingServerTagsNumBytes, nil).Once()

	mockInterface.On("TagsNumWithRepo",
		mockUser, mockPassword, mockSchema, mockAddress, "apache", "skywalking-java-agent", DefaultRequestTimeout,
	).Return(skywalkingAgentTagsNumBytes, nil).Once()

	// 4.1 mock elastic's repo tags by page
	mockInterface.On("TagsWithinRepoByPage",
		mockPage, mockPageSize,
		mockSchema, mockAddress,
		"apache", "skywalking-ui",
		mockUser, mockPassword, DefaultRequestTimeout,
	).Return(nil, mockNetErr)

	mockInterface.On("TagsWithinRepoByPage",
		mockPage, mockPageSize,
		mockSchema, mockAddress,
		"apache", "skywalking-oap-server",
		mockUser, mockPassword, DefaultRequestTimeout,
	).Return(nil, mockNetErr)

	mockInterface.On("TagsWithinRepoByPage",
		mockPage, mockPageSize,
		mockSchema, mockAddress,
		"apache", "skywalking-java-agent",
		mockUser, mockPassword, DefaultRequestTimeout,
	).Return(nil, mockNetErr)

	// 4.2 mock elastic's repo tags by page
	mockInterface.On("TagsWithinRepoByPage",
		mockPage, mockPageSize,
		mockSchema, mockAddress,
		"elastic", "filebeat",
		mockUser, mockPassword, DefaultRequestTimeout,
	).Return(nil, mockNetErr)

	mockInterface.On("TagsWithinRepoByPage",
		mockPage, mockPageSize,
		mockSchema, mockAddress,
		"elastic", "elasticsearch",
		mockUser, mockPassword, DefaultRequestTimeout,
	).Return(nil, mockNetErr)

	ImageList(command.OperationItem{
		B:         []byte(content),
		Logger:    mockLogger,
		Local:     false,
		Interface: mockInterface,
	})
}

func Test_getHandlerInterface(t *testing.T) {

	var h HandlerInterface
	r := getHandlerInterface(h)
	require.Equal(t, new(Requester), r)

	h2 := &mocks.HandlerInterface{}
	r2 := getHandlerInterface(h2)
	require.Equal(t, h2, r2)
}

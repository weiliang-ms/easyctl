package harbor

import (
	"github.com/sirupsen/logrus"
	"github.com/weiliang-ms/easyctl/pkg/util/log"
	"github.com/weiliang-ms/easyctl/pkg/util/request"
	"io"
	"net/http"
	"time"
)

//go:generate mockery --name=HandlerInterface
type HandlerInterface interface {
	DoRequest(httpRequestItem request.HTTPRequestItem) ([]byte, error)
	ProjectsByPage(
		page int, pageSize int,
		logger *logrus.Logger,
		schema string,
		address string,
		user string,
		password string, timeout time.Duration,
	) ([]byte, error)
	ProjectCount(user, password, schema, address string, timeout time.Duration) ([]byte, error)
	TagsWithinRepoByPage(page, pageSize int, schema, address, projectName, repoName, user, password string, timeout time.Duration) ([]byte, error)
	TagsNumWithRepo(user, password, schema, address, projectName, repoName string, timeout time.Duration) ([]byte, error)
	RepoCount(user, password, schema, address string, projectId int, timeout time.Duration) ([]byte, error)
	ListRepoByPage(page, pageSize int, user, password, schema, address, projectName string, timeout time.Duration) ([]byte, error)
}

// Requester request harbor
type Requester struct{}

func (r Requester) DoRequest(httpRequestItem request.HTTPRequestItem) ([]byte, error) {
	request, err := http.NewRequest(httpRequestItem.Method, httpRequestItem.Url, httpRequestItem.Body)
	if err != nil {
		return nil, err
	}
	request.SetBasicAuth(httpRequestItem.User, httpRequestItem.Password)

	client := http.Client{Timeout: httpRequestItem.Timeout}
	resp, err := client.Do(request)
	if err != nil && !httpRequestItem.Mock {
		return nil, err
	}

	return io.ReadAll(resp.Body)
}

// ProjectsByPage 按页查询harbor内project
func (r Requester) ProjectsByPage(
	page int, pageSize int,
	logger *logrus.Logger,
	schema string,
	address string,
	user string,
	password string,
	timeout time.Duration,
) ([]byte, error) {
	logger = log.SetDefault(logger)

	return r.DoRequest(request.HTTPRequestItem{
		Url:      ApiGetProjectsByPage(schema, address, page, pageSize),
		Method:   http.MethodGet,
		Body:     nil,
		Headers:  nil,
		Timeout:  timeout,
		User:     user,
		Password: password,
	})
}

// TagsWithinRepoByPage 获取repo的全部tag
func (r Requester) TagsWithinRepoByPage(page, pageSize int, schema, address, projectName, repoName, user, password string, timeout time.Duration) ([]byte, error) {
	return r.DoRequest(request.HTTPRequestItem{
		Url: ApiGetTagsWithinRepoByPage(
			schema,
			address,
			projectName,
			repoName,
			page,
			pageSize),
		Method:   http.MethodGet,
		Body:     nil,
		Headers:  map[string]string{"accept": "application/json"},
		Timeout:  timeout,
		User:     user,
		Password: password,
	})
}

func (r Requester) TagsNumWithRepo(user, password, schema, address, projectName, repoName string, timeout time.Duration) ([]byte, error) {
	return r.DoRequest(request.HTTPRequestItem{
		Url: ApiGetTagsNumWithRepo(
			schema,
			address,
			projectName,
			repoName),
		Method:   http.MethodGet,
		Body:     nil,
		Headers:  map[string]string{"accept": "application/json"},
		Timeout:  timeout,
		User:     user,
		Password: password,
	})
}

func (r Requester) ProjectCount(user, password, schema, address string, timeout time.Duration) ([]byte, error) {
	return r.DoRequest(request.HTTPRequestItem{
		Url:      ApiGetStatistics(schema, address),
		Method:   http.MethodGet,
		Body:     nil,
		Headers:  nil,
		Timeout:  timeout,
		User:     user,
		Password: password,
	})
}

func (r Requester) RepoCount(user, password, schema, address string, projectId int, timeout time.Duration) ([]byte, error) {
	return r.DoRequest(request.HTTPRequestItem{
		Url:      ApiGetProjectMetaInfo(schema, address, projectId),
		Method:   http.MethodGet,
		Body:     nil,
		Headers:  nil,
		Timeout:  timeout,
		User:     user,
		Password: password,
	})
}

// ListRepoByPage 调用harbor api获取project列表，根据project name获取repo列表
func (r Requester) ListRepoByPage(page, pageSize int, user, password, schema, address, projectName string, timeout time.Duration) ([]byte, error) {
	return r.DoRequest(request.HTTPRequestItem{
		Url:      ApiGetListRepoByPage(schema, address, projectName, page, pageSize),
		Method:   http.MethodGet,
		Body:     nil,
		Headers:  nil,
		Timeout:  timeout,
		User:     user,
		Password: password,
	})
}

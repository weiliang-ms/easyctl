package chart

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/weiliang-ms/easyctl/pkg/export/chart/mocks"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	"github.com/weiliang-ms/easyctl/pkg/util/request"
	"net/http"
	"os"
	"sort"
	"testing"
	"time"
)

var (
	adminerBytes      = []byte("adminer-1.0.0.tgz")
	httpdBytes        = []byte("httpd-2.0.0.tgz")
	nginxBytes        = []byte("nginx-3.0.0.tgz")
	mockSchema        = "http"
	mockEndpoint      = "1.1.1.1:80"
	mockUser          = "admin"
	mockPasswd        = "admin"
	mockPreserveDir   = "charts"
	mockPackage       = false
	mockRepoName      = "chart"
	mockRequestTimout = time.Millisecond
	mockLogger        = logrus.New()
	mockErr           = fmt.Errorf("mockErr")
	mockChartList     = itemList{
		{
			Name:          "adminer",
			LatestVersion: "1.0.0",
		},
		{
			Name:          "nginx",
			LatestVersion: "3.0.0",
		},
		{
			Name:          "httpd",
			LatestVersion: "2.0.0",
		},
	}
)

//go:embed mocks/chart-list.json
var chartListBytes []byte

func Test_ListNotFound_Mock(t *testing.T) {
	h := &mocks.HandlerInterface{}

	r := request.HTTPRequestItem{
		Url:      ApiListChartsUrl(mockSchema, mockEndpoint),
		Method:   http.MethodGet,
		Body:     nil,
		Headers:  nil,
		User:     mockUser,
		Timeout:  mockRequestTimout,
		Password: mockPasswd,
		Mock:     false,
	}
	m := &Manager{
		Config: Config{
			Schema:      mockSchema,
			Endpoint:    mockEndpoint,
			Username:    mockUser,
			Password:    mockPasswd,
			PreserveDir: mockPreserveDir,
			Package:     mockPackage,
			RepoName:    mockRepoName,
		},
		Logger:             mockLogger,
		Handler:            h,
		HttpRequestTimeout: mockRequestTimout,
	}

	h.On("GetChartList", r).Return(nil, nil)

	re, err := m.List()
	_, ok := err.(*json.SyntaxError)
	require.Equal(t, 0, len(re))
	require.Equal(t, true, ok)
}

func Test_List_Mock(t *testing.T) {
	h := &mocks.HandlerInterface{}

	r := request.HTTPRequestItem{
		Url:      ApiListChartsUrl(mockSchema, mockEndpoint),
		Method:   http.MethodGet,
		Body:     nil,
		Headers:  nil,
		User:     mockUser,
		Timeout:  mockRequestTimout,
		Password: mockPasswd,
		Mock:     false,
	}

	m := &Manager{
		Config: Config{
			Schema:      mockSchema,
			Endpoint:    mockEndpoint,
			Username:    mockUser,
			Password:    mockPasswd,
			PreserveDir: mockPreserveDir,
			Package:     mockPackage,
			RepoName:    mockRepoName,
		},
		Logger:             mockLogger,
		Handler:            h,
		HttpRequestTimeout: mockRequestTimout,
	}

	h.On("GetChartList", r).Return(chartListBytes, nil)

	re, err := m.List()
	require.Nil(t, err)
	require.Equal(t, 3, len(re))
}

func Test_ListErr_Mock(t *testing.T) {
	h := &mocks.HandlerInterface{}

	r := request.HTTPRequestItem{
		Url:      ApiListChartsUrl(mockSchema, mockEndpoint),
		Method:   http.MethodGet,
		Body:     nil,
		Headers:  nil,
		User:     mockUser,
		Timeout:  mockRequestTimout,
		Password: mockPasswd,
		Mock:     false,
	}
	m := &Manager{
		Config: Config{
			Schema:      mockSchema,
			Endpoint:    mockEndpoint,
			Username:    mockUser,
			Password:    mockPasswd,
			PreserveDir: mockPreserveDir,
			Package:     mockPackage,
			RepoName:    mockRepoName,
		},
		Logger:             mockLogger,
		Handler:            h,
		HttpRequestTimeout: mockRequestTimout,
	}

	h.On("GetChartList", r).Return(nil, mockErr)

	re, err := m.List()
	require.Equal(t, mockErr, err)
	require.Equal(t, 0, len(re))
}

func Test_GetChartsByte_Mock(t *testing.T) {
	h := &mocks.HandlerInterface{}

	url1 := ApiChartBytesUrl(mockSchema, mockEndpoint, "adminer-1.0.0.tgz")
	url2 := ApiChartBytesUrl(mockSchema, mockEndpoint, "httpd-2.0.0.tgz")
	url3 := ApiChartBytesUrl(mockSchema, mockEndpoint, "nginx-3.0.0.tgz")
	r := request.HTTPRequestItem{
		Url:      url1,
		Method:   http.MethodGet,
		Body:     nil,
		Headers:  nil,
		User:     mockUser,
		Timeout:  mockRequestTimout,
		Password: mockPasswd,
		Mock:     false,
	}
	m := &Manager{
		Config: Config{
			Schema:      mockSchema,
			Endpoint:    mockEndpoint,
			Username:    mockUser,
			Password:    mockPasswd,
			PreserveDir: mockPreserveDir,
			Package:     mockPackage,
			RepoName:    mockRepoName,
		},
		Logger:             mockLogger,
		Handler:            h,
		HttpRequestTimeout: mockRequestTimout,
	}

	r.Url = url1
	h.On("GetChartByte", r).Return(adminerBytes, nil)
	r.Url = url2
	h.On("GetChartByte", r).Return(httpdBytes, nil)
	r.Url = url3
	h.On("GetChartByte", r).Return(nginxBytes, nil)

	re, err := m.GetChartsByte(mockChartList)
	require.Nil(t, err)
	require.Equal(t, 3, len(re))
}

func Test_GetChartsByteErr_Mock(t *testing.T) {
	h := &mocks.HandlerInterface{}

	url1 := ApiChartBytesUrl(mockSchema, mockEndpoint, "adminer-1.0.0.tgz")
	url2 := ApiChartBytesUrl(mockSchema, mockEndpoint, "httpd-2.0.0.tgz")
	url3 := ApiChartBytesUrl(mockSchema, mockEndpoint, "nginx-3.0.0.tgz")
	r := request.HTTPRequestItem{
		Url:      url1,
		Method:   http.MethodGet,
		Body:     nil,
		Headers:  nil,
		User:     mockUser,
		Timeout:  mockRequestTimout,
		Password: mockPasswd,
		Mock:     false,
	}

	m := &Manager{
		Config: Config{
			Schema:      mockSchema,
			Endpoint:    mockEndpoint,
			Username:    mockUser,
			Password:    mockPasswd,
			PreserveDir: mockPreserveDir,
			Package:     mockPackage,
			RepoName:    mockRepoName,
		},
		Logger:             mockLogger,
		Handler:            h,
		HttpRequestTimeout: mockRequestTimout,
	}

	r.Url = url1
	h.On("GetChartByte", r).Return(nil, mockErr)
	r.Url = url2
	h.On("GetChartByte", r).Return(nil, mockErr)
	r.Url = url3
	h.On("GetChartByte", r).Return(nil, mockErr)

	re, err := m.GetChartsByte(mockChartList)
	require.Equal(t, mockErr, err)
	require.Equal(t, 0, len(re))
}

func Test_parseConfig(t *testing.T) {
	content := `
helm-repo:
  endpoint: 10.10.1.3:80   # harbor访问地址
  domain: harbor.wl.io      # harbor域
  username: admin           # harbor用户
  password: 123456          # harbor密码
  preserveDir: /root/charts # chart包持久化目录
  package: true             # 是否打成tar包
  repo-name: charts         # chart repo harbor内的名称
`
	c, err := parseConfig([]byte(content), mockLogger)
	require.Nil(t, err)
	require.Equal(t, "10.10.1.3:80", c.Endpoint)
	require.Equal(t, "harbor.wl.io", c.Domain)
	require.Equal(t, "admin", c.Username)
	require.Equal(t, "123456", c.Password)
	require.Equal(t, true, c.Package)
	require.Equal(t, "charts", c.RepoName)

}

func Test_parseConfig_ErrCase(t *testing.T) {
	content := `
helm-repo:
  - endpoint: 10.10.1.3:80   # harbor访问地址
  domain: harbor.wl.io      # harbor域
  username: admin           # harbor用户
  password: 123456          # harbor密码
  preserveDir: /root/charts # chart包持久化目录
  package: true             # 是否打成tar包
  repo-name: charts         # chart repo harbor内的名称
`
	c, err := parseConfig([]byte(content), mockLogger)
	require.NotNil(t, err)
	require.Equal(t, "", c.Endpoint)
	require.Equal(t, "", c.Domain)
	require.Equal(t, "", c.Username)
	require.Equal(t, "", c.Password)
	require.Equal(t, false, c.Package)
	require.Equal(t, "", c.RepoName)

}

func Test_Save_Mock(t *testing.T) {

	defer os.RemoveAll("charts")

	h := &mocks.HandlerInterface{}

	url := ApiListChartsUrl(mockSchema, mockEndpoint)

	r := request.HTTPRequestItem{
		Url:      url,
		Method:   http.MethodGet,
		Body:     nil,
		Headers:  nil,
		User:     mockUser,
		Timeout:  mockRequestTimout,
		Password: mockPasswd,
		Mock:     false,
	}

	m := &Manager{
		Config: Config{
			Schema:      mockSchema,
			Endpoint:    mockEndpoint,
			Username:    mockUser,
			Password:    mockPasswd,
			PreserveDir: mockPreserveDir,
			Package:     mockPackage,
			RepoName:    mockRepoName,
		},
		Logger:             mockLogger,
		Handler:            h,
		HttpRequestTimeout: mockRequestTimout,
	}

	h.On("GetChartList", r).Return(chartListBytes, nil)
	//r.Url = url1

	h.On("GetChartByte", request.HTTPRequestItem{
		Url:      ApiChartBytesUrl(mockSchema, mockEndpoint, "adminer-1.0.0.tgz"),
		Method:   http.MethodGet,
		Body:     nil,
		Headers:  nil,
		User:     mockUser,
		Timeout:  mockRequestTimout,
		Password: mockPasswd,
		Mock:     false,
	}).Return(adminerBytes, nil)

	h.On("GetChartByte", request.HTTPRequestItem{
		Url:      ApiChartBytesUrl(mockSchema, mockEndpoint, "httpd-1.0.0.tgz"),
		Method:   http.MethodGet,
		Body:     nil,
		Headers:  nil,
		User:     mockUser,
		Timeout:  mockRequestTimout,
		Password: mockPasswd,
		Mock:     false,
	}).Return(httpdBytes, nil)

	h.On("GetChartByte", request.HTTPRequestItem{
		Url:      ApiChartBytesUrl(mockSchema, mockEndpoint, "nginx-1.0.0.tgz"),
		Method:   http.MethodGet,
		Body:     nil,
		Headers:  nil,
		User:     mockUser,
		Timeout:  mockRequestTimout,
		Password: mockPasswd,
		Mock:     false,
	}).Return(nginxBytes, nil)

	sort.Sort(mockChartList)
	err := m.Save()
	require.Equal(t, nil, err)
}

func Test_Save_GetChartByteErr_Mock(t *testing.T) {

	defer os.RemoveAll("charts")

	h := &mocks.HandlerInterface{}

	url := ApiListChartsUrl(mockSchema, mockEndpoint)

	r := request.HTTPRequestItem{
		Url:      url,
		Method:   http.MethodGet,
		Body:     nil,
		Headers:  nil,
		User:     mockUser,
		Timeout:  mockRequestTimout,
		Password: mockPasswd,
		Mock:     false,
	}

	m := &Manager{
		Config: Config{
			Schema:      mockSchema,
			Endpoint:    mockEndpoint,
			Username:    mockUser,
			Password:    mockPasswd,
			PreserveDir: mockPreserveDir,
			Package:     mockPackage,
			RepoName:    mockRepoName,
		},
		Logger:             mockLogger,
		Handler:            h,
		HttpRequestTimeout: mockRequestTimout,
	}

	h.On("GetChartList", r).Return(chartListBytes, nil)
	//r.Url = url1

	h.On("GetChartByte", request.HTTPRequestItem{
		Url:      ApiChartBytesUrl(mockSchema, mockEndpoint, "adminer-1.0.0.tgz"),
		Method:   http.MethodGet,
		Body:     nil,
		Headers:  nil,
		User:     mockUser,
		Timeout:  mockRequestTimout,
		Password: mockPasswd,
		Mock:     false,
	}).Return(nil, mockErr)

	h.On("GetChartByte", request.HTTPRequestItem{
		Url:      ApiChartBytesUrl(mockSchema, mockEndpoint, "httpd-1.0.0.tgz"),
		Method:   http.MethodGet,
		Body:     nil,
		Headers:  nil,
		User:     mockUser,
		Timeout:  mockRequestTimout,
		Password: mockPasswd,
		Mock:     false,
	}).Return(nil, mockErr)

	h.On("GetChartByte", request.HTTPRequestItem{
		Url:      ApiChartBytesUrl(mockSchema, mockEndpoint, "nginx-1.0.0.tgz"),
		Method:   http.MethodGet,
		Body:     nil,
		Headers:  nil,
		User:     mockUser,
		Timeout:  mockRequestTimout,
		Password: mockPasswd,
		Mock:     false,
	}).Return(nil, mockErr)

	sort.Sort(mockChartList)
	err := m.Save()
	require.NotNil(t, err)
}

func Test_Save_ListErr_Mock(t *testing.T) {

	defer os.RemoveAll("charts")

	h := &mocks.HandlerInterface{}

	url := ApiListChartsUrl(mockSchema, mockEndpoint)

	r := request.HTTPRequestItem{
		Url:      url,
		Method:   http.MethodGet,
		Body:     nil,
		Headers:  nil,
		User:     mockUser,
		Timeout:  mockRequestTimout,
		Password: mockPasswd,
		Mock:     false,
	}

	m := &Manager{
		Config: Config{
			Schema:      mockSchema,
			Endpoint:    mockEndpoint,
			Username:    mockUser,
			Password:    mockPasswd,
			PreserveDir: mockPreserveDir,
			Package:     mockPackage,
			RepoName:    mockRepoName,
		},
		Logger:             mockLogger,
		Handler:            h,
		HttpRequestTimeout: mockRequestTimout,
	}

	h.On("GetChartList", r).Return(chartListBytes, mockErr)
	//r.Url = url1

	sort.Sort(mockChartList)
	err := m.Save()
	require.NotNil(t, err)
}

func Test_Chart(t *testing.T) {

	defer func() {
		r := recover()
		if r != nil {
			fmt.Println(r)
		}
	}()

	Run(command.OperationItem{
		B:          nil,
		Logger:     nil,
		OptionFunc: nil,
		Interface:  nil,
		UnitTest:   false,
		Mock:       false,
		LocalRun:   false,
	})
}

func Test_Chart_ParseErr(t *testing.T) {

	defer func() {
		r := recover()
		if r != nil {
			fmt.Println(r)
		}
	}()

	content := `
helm-repo:
  - endpoint: 10.10.1.3:80   # harbor访问地址
  domain: harbor.wl.io      # harbor域
  username: admin           # harbor用户
  password: 123456          # harbor密码
  preserveDir: /root/charts # chart包持久化目录
  package: true             # 是否打成tar包
  repo-name: charts         # chart repo harbor内的名称
`

	Run(command.OperationItem{
		B:          []byte(content),
		Logger:     nil,
		OptionFunc: nil,
		Interface:  nil,
		UnitTest:   false,
		Mock:       false,
		LocalRun:   false,
	})
}

func Test_getHandlerInterface(t *testing.T) {

	var h HandlerInterface
	r := getHandlerInterface(h)
	require.Equal(t, new(Handler), r)

	h2 := &mocks.HandlerInterface{}
	r2 := getHandlerInterface(h2)
	require.Equal(t, h2, r2)
}

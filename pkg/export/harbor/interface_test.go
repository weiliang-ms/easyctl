package harbor

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/weiliang-ms/easyctl/pkg/util/request"
	"testing"
	"time"
)

func TestInterface(t *testing.T) {
	var r Requester
	_, _ = r.ProjectsByPage(mockPage, mockPageSize, mockLogger, mockSchema, mockAddress, mockUser, mockPassword, time.Millisecond)
	_, _ = r.ProjectCount(mockUser, mockPassword, mockSchema, mockAddress, time.Millisecond)
	_, _ = r.RepoCount(mockUser, mockPassword, mockSchema, mockAddress, mockProjectId, time.Millisecond)
	_, _ = r.TagsWithinRepoByPage(mockPage, mockPageSize, mockSchema, mockAddress, mockProjectName, mockRepoName, mockUser, mockPassword, time.Millisecond)
	_, _ = r.TagsNumWithRepo(mockUser, mockPassword, mockSchema, mockAddress, mockProjectName, mockRepoName, time.Millisecond)
	_, _ = r.ListRepoByPage(mockPage, mockPageSize, mockUser, mockPassword, mockSchema, mockAddress, mockProjectName, time.Millisecond)

	// test DoRequest
	b, err := r.DoRequest(request.HTTPRequestItem{Method: "！！！"})
	require.Nil(t, b)
	require.Equal(t, "net/http: invalid method \"！！！\"", err.Error())

	// test DoRequest
	defer func() {
		r := recover()
		require.Equal(t, "runtime error: invalid memory address or nil pointer dereference", fmt.Sprintf("%s", r))
	}()
	r.DoRequest(request.HTTPRequestItem{Mock: true})

}

package harbor

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestApis(t *testing.T) {
	require.Equal(t, mockApiGetStatisticsUrl, ApiGetStatistics(mockSchema, mockAddress))
	require.Equal(t, mockApiGetProjectsByPageUrl, ApiGetProjectsByPage(mockSchema, mockAddress, mockPage, mockPageSize))
	require.Equal(t, mockApiGetProjectsByNameUrl, ApiGetProjectsByName(mockSchema, mockAddress, mockProjectName))
	require.Equal(t, mockApiGetListRepoByPageUrl, ApiGetListRepoByPage(mockSchema, mockAddress, mockProjectName, mockPage, mockPageSize))
	require.Equal(t, mockApiGetTagsWithinRepoByPageUrl, ApiGetTagsWithinRepoByPage(mockSchema, mockAddress, mockProjectName, mockRepoName, mockPage, mockPageSize))
	require.Equal(t, mockApiGetTagsNumWithRepoUrl, ApiGetTagsNumWithRepo(mockSchema, mockAddress, mockProjectName, mockRepoName))
	require.Equal(t, mockApiGetProjectMetaInfoUrl, ApiGetProjectMetaInfo(mockSchema, mockAddress, mockProjectId))
}


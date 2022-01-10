package harbor

import (
	"fmt"
	"strings"
)

func ApiGetStatistics(schema, address string) string {
	return fmt.Sprintf("%s://%s/api/v2.0/statistics", schema, address)
}

func ApiGetProjectsByPage(schema, address string, page, pageSize int) string {
	return fmt.Sprintf("%s://%s/api/v2.0/projects?page=%d&page_size=%d",
		schema, address, page, pageSize)
}

func ApiGetProjectsByName(schema, address, projectName string) string {
	return fmt.Sprintf("%s://%s/api/v2.0/search?q=%s",
		schema, address, projectName)
}

func ApiGetListRepoByPage(schema, address, projectName string, page, pageSize int) string {
	return fmt.Sprintf("%s://%s/api/v2.0/projects/%s/repositories?page=%d&page_size=%d",
		schema,
		address,
		projectName,
		page,
		pageSize,
	)
}

func ApiGetTagsWithinRepoByPage(schema, address, projectName, repoName string, page, pageSize int) string {
	args := fmt.Sprintf("page=%d"+
		"&page_size=%d&with_tag=true"+
		"&with_label=false&with_scan_overview=false"+
		"&with_signature=false"+
		"&with_immutable_status=false", page, pageSize)
	return fmt.Sprintf("%s://%s/api/v2.0/projects/%s/repositories/%s/artifacts?%s",
		schema,
		address,
		projectName,
		strings.ReplaceAll(repoName, "/", "%252F"),
		args)
}

func ApiGetTagsNumWithRepo(schema, address, projectName, repoName string) string {
	return fmt.Sprintf("%s://%s/api/v2.0/projects/%s/repositories/%s",
		schema,
		address,
		projectName,
		strings.ReplaceAll(repoName, "/", "%252F"),
	)
}

func ApiGetProjectMetaInfo(schema, address string, projectId int) string {
	return fmt.Sprintf("%s://%s/api/v2.0/projects/%d", schema, address, projectId)
}

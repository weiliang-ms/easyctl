package chart

import (
	"github.com/sirupsen/logrus"
	"github.com/weiliang-ms/easyctl/pkg/export"
	"os"
	"testing"
)

var (
	testHarborExecutor *export.HarborExecutor
	testLogger         *logrus.Logger
)

//func TestParseHarborRepoConfig(t *testing.T) {
//	b, readErr := os.ReadFile("../../../asset/config.yaml")
//	if readErr != nil {
//		panic(readErr)
//	}
//
//	logger := logrus.New()
//	logger.SetLevel(logrus.DebugLevel)
//
//	c, err := export.ParseHarborConfig(b, logger)
//	if err != nil {
//		panic(err)
//	}
//	fmt.Printf("%v\n", c.HarborConfig)
//	testHarborExecutor = c
//	testLogger = logger
//}
//
//func TestProjectsByPage(t *testing.T) {
//
//	list , err := testHarborExecutor.ProjectsByPage(3)
//	if err != nil {
//		panic(err)
//	}
//
//	fmt.Printf("数量为: %d", len(list))
//	for _, v := range list{
//		fmt.Println(v)
//	}
//}
//
//func TestAllProjects(t *testing.T) {
//	b, readErr := os.ReadFile("../../../asset/config.yaml")
//	if readErr != nil {
//		panic(readErr)
//	}
//
//	logger := logrus.New()
//	logger.SetLevel(logrus.InfoLevel)
//
//	executor, err := export.ParseHarborConfig(b, logger)
//	if err != nil {
//		panic(err)
//	}
//	execErr := executor.AllProjects()
//	if execErr != nil {
//		 panic(execErr)
//	}
//
//	fmt.Printf("数量为：%d\n", len(executor.ProjectSlice))
//}
//
//func TestProjectList(t *testing.T) {
//	b, readErr := os.ReadFile("../../../asset/config.yaml")
//	if readErr != nil {
//		panic(readErr)
//	}
//
//	logger := logrus.New()
//	logger.SetLevel(logrus.InfoLevel)
//
//	executor, err := export.ParseHarborConfig(b, logger)
//	if err != nil {
//		panic(err)
//	}
//	execErr := executor.ProjectList()
//	if execErr != nil {
//		panic(execErr)
//	}
//
//	fmt.Printf("数量为：%d\n", len(executor.ProjectSlice))
//}
//
//func TestFilterProjects(t *testing.T) {
//	b, readErr := os.ReadFile("../../../asset/config.yaml")
//	if readErr != nil {
//		panic(readErr)
//	}
//
//	logger := logrus.New()
//	logger.SetLevel(logrus.InfoLevel)
//
//	executor, err := export.ParseHarborConfig(b, logger)
//	if err != nil {
//		panic(err)
//	}
//	execErr := executor.ProjectList()
//	if execErr != nil {
//		panic(execErr)
//	}
//
//	executor.FilterProjects()
//	fmt.Printf("过滤后数量为：%d\n -> %v", len(executor.ProjectSlice), executor.ProjectSlice)
//}
//
//func TestListRepoByPage(t *testing.T) {
//
//	_ , err := testHarborExecutor.ListRepoByPage(1, "apache")
//	if err != nil {
//		panic(err)
//	}
//}
//
//func TestReposWithinProject(t *testing.T) {
//	testHarborExecutor.ReposInProject = make(map[string][]string)
//	err := testHarborExecutor.ReposWithinProject("apache")
//	if err != nil {
//		panic(err)
//	}
//
//	//fmt.Printf("%v", testHarborExecutor.Repos)
//}
//
//func TestReposWithinProjects(t *testing.T) {
//
//	err := testHarborExecutor.ReposWithinProjects()
//	if err != nil {
//		panic(err)
//	}
//
//	//fmt.Printf("%v", testHarborExecutor.Repos)
//}
//
//func TestTagsWithinRepoByPage(t *testing.T) {
//
//	tags, err := testHarborExecutor.TagsWithinRepoByPage(1, "paas","champ-portal")
//	if err != nil {
//		panic(err)
//	}
//
//	fmt.Printf("%v", tags)
//}
//
//func TestTagsWithinRepo(t *testing.T) {
//
//	testHarborExecutor.TagsInProject = make(map[string][]string)
//	err := testHarborExecutor.TagsWithinRepo("paas","champ-portal")
//	if err != nil {
//		panic(err)
//	}
//
//	fmt.Printf("%v", testHarborExecutor.TagsInProject)
//}
//
//func TestTagsWithinProjects(t *testing.T) {
//
//	testHarborExecutor.TagsInProject = make(map[string][]string)
//
//	err := testHarborExecutor.TagsWithinProjects()
//	if err != nil {
//		panic(err)
//	}
//}

func TestHarborImageList(t *testing.T) {
	b, readErr := os.ReadFile("../../../asset/config.yaml")
	if readErr != nil {
		panic(readErr)
	}

	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	export.HarborImageList(b, logger)
}

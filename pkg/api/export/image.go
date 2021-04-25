package export

import (
	"encoding/json"
	"fmt"
	"github.com/weiliang-ms/easyctl/pkg/util"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type harbor struct {
	Harbor harborMeta `yaml:"harbor"`
}

type harborMeta struct {
	User      string   `yaml:"user"`
	Password  string   `yaml:"password"`
	Address   string   `yaml:"address"`
	Namespace []string `yaml:"namespace"`
	Domain    string   `yaml:"domain"`
}

type project struct {
	Name string `json:"name"`
	Repo []string
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

// 单机本地离线
func HarborImage(path string) {
	h := parseHarbor(path).Harbor
	d := h.Domain
	p := projects(h)

	// 组装project repo 映射关系
	var projectRepoMap = make(map[string][]string)

	for _, v := range p {
		repoNames := repos(h, v.Name)
		if v.Name != "" {
			projectRepoMap[v.Name] = repoNames
		}

	}

	_ = os.Mkdir("images", 0755)

	// 初始化image-list.txt
	imageList, _ := os.OpenFile("images/image-list.txt", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
	defer imageList.Close()

	// 组装镜像库镜像Tag
	for k, v := range projectRepoMap {
		for _, r := range v {
			for _, t := range imageTags(h, k, r, d) {
				_, _ = imageList.WriteString(t)
				_, _ = imageList.WriteString("\n")
			}
		}
	}

}

func projects(harbor harborMeta) []project {
	var slice []project
	i := 1
	for {
		re := listProject(i, harbor.User, harbor.Password, harbor.Address)
		if len(re) == 0 {
			break
		}
		for _, v := range re {
			var project project
			project.Name = v.Name
			slice = append(slice, project)
		}
		i++
	}

	return slice
}

func imageTags(harbor harborMeta, projectName string, repoName string, domain string) []string {
	var tags []string
	i := 1
	for {
		re := listRepoTags(i, harbor.User, harbor.Password, harbor.Address, projectName, repoName)
		if len(re) == 0 {
			break
		}
		for _, v := range re {
			for _, t := range v.Tags {
				tag := fmt.Sprintf("%s/%s/%s:%s", domain, projectName, repoName, t.Name)
				tags = append(tags, tag)
			}
		}
		i++
	}
	return tags
}

func repos(harbor harborMeta, projectName string) []string {
	var repos []string
	i := 1
	for {
		re := listRepo(i, harbor.User, harbor.Password, harbor.Address, projectName)
		if len(re) == 0 {
			break
		}
		for _, v := range re {
			s := strings.Split(v.Name, "/")
			if len(s) > 1 {
				repos = append(repos, s[len(s)-1])
			}
		}
		i++
	}

	//fmt.Printf("repos: %v",repos)
	return repos
}

func listRepoTags(page int, username string, password string, address string, projectName string, repoName string) []tags {

	var tags []tags

	args := fmt.Sprintf("page=%d"+
		"&page_size=10&with_tag=true"+
		"&with_label=false&with_scan_overview=false"+
		"&with_signature=false"+
		"&with_immutable_status=false", page)

	url := fmt.Sprintf("http://%s/api/v2.0/projects/%s/repositories/%s/artifacts?%s",
		address, projectName, repoName, args)

	//var r io.Reader
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}
	req.SetBasicAuth(username, password)
	req.Header.Add("accept", "application/json")
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

func listProject(page int, username string, password string, address string) []project {

	var project []project

	url := fmt.Sprintf("http://%s/api/v2.0/projects?page=%d&page_size=10", address, page)
	//var r io.Reader
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}
	req.SetBasicAuth(username, password)
	req.Header.Add("accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &project); err != nil {
		fmt.Println(err.Error())
	}
	return project
}

func listRepo(page int, username string, password string, address string, projectName string) []repo {
	var repo []repo
	url := fmt.Sprintf("http://%s/api/v2.0/projects/%s/repositories?page=%d&page_size=10", address, projectName, page)
	body := util.Get(url, username, password)
	if err := json.Unmarshal(body, &repo); err != nil {
		fmt.Println(err.Error())
	}
	return repo
}

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

package harbor

import (
	"github.com/sirupsen/logrus"
	"time"
)

// 包含关系
// harbor -> 包含多个project
// project -> 包含多个repo
// Repository -> 包含多个tag

// Executor 与harbor交互的执行器
type Executor struct {
	Config
	Logger         *logrus.Logger
	ProjectSlice   []ProjectInternal
	ReposInProject map[string][]string // 项目下镜像repo集合
	TagsInProject  map[string][]string // 项目下镜像tag集合
	HandlerInterface
}

// ConfigExternal 用于反序列化harbor-repo对象配置
type ConfigExternal struct {
	HarborRepo struct {
		Schema        string   `yaml:"schema"`
		Address       string   `yaml:"address"`
		Domain        string   `yaml:"domain"`
		User          string   `yaml:"user"`
		Password      string   `yaml:"password"`
		PreserveDir   string   `yaml:"preserve-dir"`
		TagWithDomain bool     `yaml:"withDomain"`
		Projects      []string `yaml:"projects"`
		Excludes      []string `yaml:"excludes"`
	} `yaml:"harbor-repo"`
}

// Config harbor-repo对象配置，用于内部使用
type Config struct {
	HarborSchema             string
	HarborAddress            string
	HarborDomain             string
	HarborUser               string
	HarborPassword           string
	PreserveDir              string
	TagWithDomain            bool
	ProjectsToSearch         []string
	ProjectsToSearchExcludes []string
}

// ProjectInternal harbor project内部对象，包含部分必需属性
type ProjectInternal struct {
	Name      string
	ProjectId int
}

// ProjectExternal harbor project外部对象，用于反序列化
type ProjectExternal struct {
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
	Project    []ProjectExternal `json:"project"`
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
	Tags              []TagExternal `json:"tags"`
	Type              string        `json:"type"`
}

// TagExternal 用于反序列化repo tag等
type TagExternal struct {
	ArtifactId   int       `json:"artifact_id"`
	Id           int       `json:"id"`
	Immutable    bool      `json:"immutable"`
	Name         string    `json:"name"`
	PullTime     time.Time `json:"pull_time"`
	PushTime     time.Time `json:"push_time"`
	RepositoryId int       `json:"repository_id"`
	Signed       bool      `json:"signed"`
}

// Statistics project info
type Statistics struct {
	PrivateProjectCount int `json:"private_project_count"`
	PrivateRepoCount    int `json:"private_repo_count"`
	PublicProjectCount  int `json:"public_project_count"`
	PublicRepoCount     int `json:"public_repo_count"`
	TotalProjectCount   int `json:"total_project_count"`
	TotalRepoCount      int `json:"total_repo_count"`
}

// ProjectMeta project metadata
type ProjectMeta struct {
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
		AutoScan             string `json:"auto_scan"`
		EnableContentTrust   string `json:"enable_content_trust"`
		PreventVul           string `json:"prevent_vul"`
		Public               string `json:"public"`
		RetentionId          string `json:"retention_id"`
		ReuseSysCveAllowlist string `json:"reuse_sys_cve_allowlist"`
		Severity             string `json:"severity"`
	} `json:"metadata"`
	Name       string    `json:"name"`
	OwnerId    int       `json:"owner_id"`
	OwnerName  string    `json:"owner_name"`
	ProjectId  int       `json:"project_id"`
	RepoCount  int       `json:"repo_count"`
	UpdateTime time.Time `json:"update_time"`
}

type RepoArtifactInfo struct {
	ArtifactCount int       `json:"artifact_count"`
	CreationTime  time.Time `json:"creation_time"`
	Id            int       `json:"id"`
	Name          string    `json:"name"`
	ProjectId     int       `json:"project_id"`
	PullCount     int       `json:"pull_count"`
	UpdateTime    time.Time `json:"update_time"`
}

type projectMapTags struct {
	ProjectName string
	Tags        []string
}

package install

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/lithammer/dedent"
	"github.com/weiliang-ms/easyctl/asset"
	"github.com/weiliang-ms/easyctl/pkg/runner"
	"github.com/weiliang-ms/easyctl/pkg/util"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"text/template"
	"time"
)

type project struct {
	Name   string `json:"project_name"`
	Public bool   `json:"public"`
}

type pingObject struct {
	Ak       string `json:"access_key"`
	As       string `json:"access_secret"`
	Insecure bool   `json:"insecure"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Url      string `json:"url"`
}

type registry struct {
	ID              int        `json:"id"`
	Credential      credential `json:"credential"`
	Insecure        bool       `json:"insecure"`
	Name            string     `json:"name"`
	Type            string     `json:"type"`
	Url             string     `json:"url"`
	TokenServiceUrl string     `json:"token_service_url"`
}

type credential struct {
	Ak   string `json:"access_key"`
	As   string `json:"access_secret"`
	Type string `json:"type"`
}

type replica struct {
	ID            int      `json:"id"`
	Name          string   `json:"name"`
	Enabled       bool     `json:"enabled"`
	Deletion      bool     `json:"deletion"`
	Override      bool     `json:"override"`
	Description   string   `json:"description"`
	SrcRegistry   []string `json:"src_registry"`
	DestNamespace []string `json:"dest_namespace"`
	Trigger       trigger  `json:"trigger"`
	DestRegistry  registry `json:"dest_registry"`
	Filters       []string `json:"filters"`
}

type trigger struct {
	Type            string          `json:"type"`
	TriggerSettings triggerSettings `json:"trigger_settings"`
}

type triggerSettings struct {
	Cron string `json:"cron"`
}

var (
	registryData registry
	cred         credential
	repl         replica
)

func init() {

	cred = credential{
		Ak:   "admin",
		Type: "harbor",
	}

	registryData = registry{
		Credential: cred,
		Insecure:   false,
		Name:       "replication",
		Type:       "harbor",
	}

	repl = replica{
		Name:         "replication",
		Enabled:      true,
		Deletion:     false,
		Override:     true,
		Trigger:      trigger{Type: "event_based"},
		DestRegistry: registryData,
	}
}

var (
	DockerConfigTempl = template.Must(template.New("DockerConfig").Parse(
		dedent.Dedent(`{
  "log-opts": {
    "max-size": "5m",
    "max-file":"3"
  },
  {{- if .Mirrors }}
  "registry-mirrors": [{{ .Mirrors }}],
  {{- end}}
  {{- if .InsecureRegistries }}
  "insecure-registries": ["{{ .InsecureRegistries }}"],
  {{- end}}
  "exec-opts": ["native.cgroupdriver=systemd"]
}
    `)))
)

func Harbor(i runner.Installer) {
	re := runner.ParseHarborServerList(i.ServerListPath)
	list := re.Harbor.Server

	if i.Offline && len(list) == 2 {
		//addDockerRegistry(re.Harbor)
		//offlineInstall(i,re.Harbor)
		initHarbor(re.Harbor)
	}
}

func offlineInstall(i runner.Installer, h runner.Harbor) {
	var wg sync.WaitGroup
	name := "harbor-offline-installer-v2.1.4.tgz"
	script, _ := asset.Asset("static/script/install_harbor.sh")

	for _, v := range h.Server {
		wg.Add(1)
		runner.ScpFile(name, fmt.Sprintf("/tmp/%s", name), v, 0755)
		i.Cmd = fmt.Sprintf("password=%s domain=%s http_port=%s data_dir=%s resolv_ip=%s %s",
			h.Password, h.Domain, h.HttpPort, h.DataDir, another(h, v.Host), string(script))
		go func(server runner.Server) {
			defer wg.Done()
			server.RemoteShell(i.Cmd)
		}(v)
	}
	wg.Wait()
}

func another(h runner.Harbor, ip string) string {
	for _, v := range h.Server {
		if v.Host != ip {
			return v.Host
		}
	}
	return "127.0.0.1"
}

func addRegistry(h runner.Harbor) {

	pingData := pingObject{
		Ak: "admin", As: h.Password, Insecure: false,
		Name: "replication", Type: "harbor",
		Url: fmt.Sprintf("http://%s", h.Domain),
	}

	registryData.Credential.As = h.Password
	registryData.Url = fmt.Sprintf("http://%s", h.Domain)

	// 序列化
	pingBody, _ := json.Marshal(&pingData)

	registryBody, _ := json.Marshal(&registryData)

	// 检测harbor状态
	for _, v := range h.Server {
		url := fmt.Sprintf("http://%s:%s/api/v2.0/registries/ping", h.Domain, h.HttpPort)
		log.Printf("Registry -> 检测%s可达性...", v.Host)
		rep, _ := post(url, bytes.NewBuffer(pingBody), "admin", h.Password)
		if !successful(rep.StatusCode) {
			log.Printf("Registry -> %s不可达...", v.Host)
			b, _ := ioutil.ReadAll(rep.Body)
			log.Println(string(b))
			log.Fatal(rep.StatusCode)
		}
		log.Println("Registry -> 可达...")
	}

	// 添加registry
	for _, v := range h.Server {
		url := fmt.Sprintf("http://%s:%s/api/v2.0/registries", v.Host, h.HttpPort)
		log.Printf("Registry -> %s添加...", v.Host)
		rep, _ := post(url, bytes.NewBuffer(registryBody), "admin", h.Password)
		if !successful(rep.StatusCode) {
			log.Printf("Registry -> %s添加失败...", v.Host)
			b, _ := ioutil.ReadAll(rep.Body)
			log.Println(string(b))
			log.Fatal(rep.StatusCode)
		}
		log.Printf("Registry -> %s添加成功...", v.Host)
	}
}

func setReplication(h runner.Harbor) {
	//var url string
	repl.DestRegistry = registryData
	repl.DestRegistry.Credential = cred
	repl.DestRegistry.Credential.As = h.Password
	trigger := trigger{
		Type:            "event_based",
		TriggerSettings: triggerSettings{Cron: ""},
	}
	repl.Trigger = trigger
	repl.Filters = []string{}

	registryData.Url = fmt.Sprintf("http://%s:%s", h.Domain, h.HttpPort)

	for _, v := range h.Server {
		// 设置replica规则
		host := fmt.Sprintf("http://%s:%s", v.Host, h.HttpPort)
		setReplica(repl, host, h.Password)
	}
}

func setReplica(repl replica, host string, password string) {

	log.Printf("Registry ID -> 获取 %s registry id...", host)
	searchRegistryUrl := host + "/api/v2.0/registries?q%3Dname=replication"
	repl.DestRegistry.ID = registryID(searchRegistryUrl, "admin", password)

	// 序列化
	replicaBody, _ := json.Marshal(&repl)
	log.Printf("Replicas -> %s添加复制规则...", host)
	url := fmt.Sprintf("%s/api/v2.0/replication/policies", host)
	rep, _ := post(url, bytes.NewBuffer(replicaBody), "admin", password)

	if successful(rep.StatusCode) {
		log.Printf("Replicas -> %s添加复制规则成功...", host)
	} else {
		log.Println(rep.StatusCode)
		b, _ := ioutil.ReadAll(rep.Body)
		log.Println(string(b))
		log.Fatalf("Replicas -> %s添加复制规则失败...", host)
	}

}

// 根据registry name获取registry id
func registryID(url string, username string, password string) int {
	var registry []registry
	b := util.Get(url, username, password)
	json.Unmarshal(b, &registry)
	log.Println(string(b))
	if len(registry) != 0 {
		return registry[0].ID
	}
	return 0
}

func initHarbor(re runner.Harbor) {

	url := fmt.Sprintf("http://%s/api/v2.0/projects", re.Domain)
	health := fmt.Sprintf("http://%s/api/v2.0/health", re.Domain)
	private := re.Project.Private
	public := re.Project.Public

	// 组装
	var projects []project
	for _, v := range private {
		var p project
		p.Public = false
		p.Name = v
		projects = append(projects, p)
	}

	for _, v := range public {
		var p project
		p.Public = true
		p.Name = v
		projects = append(projects, p)
	}

	if !isHealthy(health, 5, time.Second*5) {
		log.Fatal("harbor访问异常...")
	}

	addRegistry(re)
	setReplication(re)
	newProject(projects, url, re.Password)
}

func addDockerRegistry(re runner.Harbor) {

	var wg sync.WaitGroup
	content, _ := util.Render(DockerConfigTempl, util.Data{
		"InsecureRegistries": re.Domain,
	})
	for _, v := range re.Server {
		wg.Add(1)
		v.WriteRemoteFile([]byte(content), "/etc/docker/daemon.json", 0755)

		go func(server runner.Server) {
			defer wg.Done()
			server.RemoteShell("systemctl daemon-reload && systemctl restart docker")
		}(v)
	}
	wg.Wait()
}

func isHealthy(url string, retry int, interval time.Duration) bool {
	for {
		if retry == 0 {
			break
		}
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		reps, _ := http.DefaultClient.Do(req)
		if reps.StatusCode == http.StatusOK {
			return true
		}
		time.Sleep(interval)
	}
	return false
}

func newProject(projects []project, url string, password string) {
	for _, v := range projects {
		b, _ := json.Marshal(v)
		resp, _ := post(url, bytes.NewBuffer(b), "admin", password)
		log.Printf("创建项目 -> %s", v.Name)
		if successful(resp.StatusCode) {
			log.Printf("项目%s -> 创建成功...", v.Name)
		} else {
			log.Fatal(resp.StatusCode)
		}
	}
}

func post(url string, body io.Reader, user string, password string) (*http.Response, error) {

	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		panic(err)
	}

	req.SetBasicAuth(user, password)
	req.Header.Add("accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	return http.DefaultClient.Do(req)
}

func successful(code int) bool {
	switch code {
	case http.StatusOK:
		return true
	case http.StatusCreated:
		return true
	case http.StatusConflict:
		return true
	default:
		return false
	}
}

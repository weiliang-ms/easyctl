package harbor

//
//import (
//	"bytes"
//	"encoding/json"
//	"fmt"
//	"github.com/lithammer/dedent"
//	"github.com/weiliang-ms/easyctl/pkg/runner"
//	"github.com/weiliang-ms/easyctl/pkg/util"
//	"io/ioutil"
//	"log"
//	"net/http"
//	"os"
//	"sync"
//	"text/template"
//	"time"
//)
//
//type project struct {
//	Name   string `json:"project_name"`
//	Public bool   `json:"public"`
//}
//
//type pingObject struct {
//	Ak       string `json:"access_key"`
//	As       string `json:"access_secret"`
//	Insecure bool   `json:"insecure"`
//	Name     string `json:"name"`
//	Type     string `json:"type"`
//	Url      string `json:"url"`
//}
//
//type registry struct {
//	ID              int        `json:"id"`
//	Credential      credential `json:"credential"`
//	Insecure        bool       `json:"insecure"`
//	Name            string     `json:"name"`
//	Type            string     `json:"type"`
//	Url             string     `json:"url"`
//	TokenServiceUrl string     `json:"token_service_url"`
//}
//
//type credential struct {
//	Ak   string `json:"access_key"`
//	As   string `json:"access_secret"`
//	Type string `json:"type"`
//}
//
//type replica struct {
//	ID            int      `json:"id"`
//	Name          string   `json:"name"`
//	Enabled       bool     `json:"enabled"`
//	Deletion      bool     `json:"deletion"`
//	Override      bool     `json:"override"`
//	Description   string   `json:"description"`
//	SrcRegistry   []string `json:"src_registry"`
//	DestNamespace []string `json:"dest_namespace"`
//	Trigger       trigger  `json:"trigger"`
//	DestRegistry  registry `json:"dest_registry"`
//	Filters       []string `json:"filters"`
//}
//
//type trigger struct {
//	Type            string          `json:"type"`
//	TriggerSettings triggerSettings `json:"trigger_settings"`
//}
//
//type triggerSettings struct {
//	Cron string `json:"cron"`
//}
//
//var (
//	registryData registry
//	cred         credential
//	repl         replica
//)
//
//const (
//	reloadDockerDaemonCmd  = "systemctl daemon-reload && systemctl restart docker"
//	dockerDaemonConfigPath = "/etc/docker/daemon.json"
//)
//
//func init() {
//
//	cred = credential{
//		Ak:   "admin",
//		Type: "harbor",
//	}
//
//	registryData = registry{
//		Credential: cred,
//		Insecure:   false,
//		Name:       "replication",
//		Type:       "harbor",
//	}
//
//	repl = replica{
//		Name:         "replication",
//		Enabled:      true,
//		Deletion:     false,
//		Override:     true,
//		Trigger:      trigger{Type: "event_based"},
//		DestRegistry: registryData,
//	}
//}
//
//var (
//	DockerConfigTempl = template.Must(template.New("DockerConfig").Parse(
//		dedent.Dedent(`{
//  "log-opts": {
//    "max-size": "5m",
//    "max-file":"3"
//  }
//  {{- if .Mirrors }}
//  ,"registry-mirrors": [{{ .Mirrors }}]
//  {{- end}}
//  {{- if .InsecureRegistries }}
//  ,"insecure-registries": ["{{ .InsecureRegistries }}"]
//  {{- end}}
//}
//    `)))
//
//	HarborConfigTmpl = template.Must(template.New("HarborConfig").Parse(dedent.Dedent(`
//hostname: {{ .domain }}
//
//# http related config
//http:
//  # port for http, default is 80. If https enabled, this port will redirect to https port
//  port: {{ .http_port }}
//
//# https related config
//# https:
//  # https port for harbor, default is 443
//  # port: 443
//  # The path of cert and key files for nginx
//  # certificate: /your/certificate/path
//  # private_key: /your/private/key/path
//harbor_admin_password: {{ .admin_password }}
//database:
//  password: root123
//  max_idle_conns: 50
//  max_open_conns: 1000
//data_volume: {{ .data_dir }}
//clair:
//  updaters_interval: 12
//trivy:
//  ignore_unfixed: false
//  skip_update: false
//jobservice:
//  max_job_workers: 40
//notification:
//  webhook_job_max_retry: 10
//chart:
//  absolute_url: disabled
//log:
//  level: info
//  local:
//    rotate_size: 200M
//    location: /var/log/harbor
//_version: 2.0.0
//proxy:
//  http_proxy:
//  https_proxy:
//  no_proxy:
//  components:
//    - core
//    - jobservice
//    - clair
//    - trivy
//`)))
//	HarborPrepareTmpl = template.Must(template.New("HarborPrepareCmd").Parse(dedent.Dedent(`
//sed -i '/docker-compose/d' /etc/rc.local
//cat <<EOF >>/etc/rc.local
//docker-compose -f /usr/local/harbor/docker-compose.yml down
//docker-compose -f /usr/local/harbor/docker-compose.yml up -d
//EOF
//
//chmod +x /etc/rc.local
//
//sed -i "/{{ .domain }}/d" /etc/hosts
//echo "{{ .address }} {{ .domain }}" >> /etc/hosts
//
//sh -c /usr/local/harbor/prepare
//`)))
//	HarborInstallTmpl = template.Must(template.New("HarborInstallCmd").Parse(dedent.Dedent(`
///usr/local/harbor/install.sh --with-clair --with-trivy --with-chartmuseum
//sed -i "/volumes:/a\      - /etc/hosts:/etc/hosts" /usr/local/harbor/docker-compose.yml
//cd /usr/local/harbor
//docker-compose down
//docker-compose up -d
//
//firewall-cmd --zone=public --add-port={{ .http_port }}/tcp --permanent
//firewall-cmd --reload
//`)))
//)
//
//const (
//	harborYml = "/usr/local/harbor/harbor.yml"
//)
//
//type harborMgr struct {
//	Installer      runner.Installer
//	Harbor         runner.Harbor
//	HarborProjects []project
//	ProjectsMap    map[int]ProjectMetadataSub
//}
//
//type ProjectMetadataSub struct {
//	Name        string
//	RetentionID string
//}
//
//// harbor mgr
//var hm harborMgr
//
//func Harbor(i runner.Installer) {
//
//	if i.ServerListPath != "" {
//		hm.Harbor = runner.ParseServerList(i.ServerListPath, runner.HarborServerList{}).Harbor.Attribute
//		hm.Installer = i
//	}
//
//	if i.Offline {
//		harborOffline()
//	}
//}
//
//func harborOffline() {
//	hm.
//		reloadDockerDaemon().
//		offlineInstall().
//		defaultProjects().
//		newProject().
//		gc().
//		listProjects().
//		setProjectsReserveNumPolicy()
//
//	if len(hm.Harbor.Server) > 1 {
//		hm.
//			addRegistries().
//			setReplications()
//	}
//}
//
//// 安装
//func (hm harborMgr) offlineInstall() *harborMgr {
//	return func() *harborMgr {
//		hm.initHarborEnv()
//		hm.cleanRetainData()
//		hm.depress()
//		hm.config()
//		hm.loadHarborImages()
//		hm.prepare()
//		hm.install()
//		return &hm
//	}()
//}
//
//func (hm *harborMgr) initHarborEnv() {
//
//	hm.detectHarborRuntime()
//
//	if len(hm.Harbor.Server) != 0 {
//		name := "harbor-offline-installer-v2.1.4.tgz"
//		log.Println("拷贝安装介质...")
//		for _, v := range hm.Harbor.Server {
//			runner.ScpFile(fmt.Sprintf("src/%s", name), fmt.Sprintf("/tmp/%s", name), v, 0755)
//		}
//	} else {
//
//	}
//}
//
//// 检测运行时
//func (hm *harborMgr) detectHarborRuntime() {
//	if len(hm.Harbor.Server) > 0 {
//		for _, v := range hm.Harbor.Server {
//			log.Printf("%s -> 检测docker...\n", v.Host)
//			if !v.CommandExists("docker") {
//				log.Fatalf("%s -> 检测失败...", v.Host)
//			}
//			log.Printf("%s -> 检测通过...\n", v.Host)
//
//			log.Printf("%s -> 检测docker-compose...\n", v.Host)
//			for _, v := range hm.Harbor.Server {
//				if !v.CommandExists("docker-compose") {
//					log.Fatalf("%s -> 检测失败...", v.Host)
//				}
//			}
//			log.Printf("%s -> 检测通过...\n", v.Host)
//
//			log.Printf("%s -> 初始化目录...\n", v.Host)
//			for _, v := range hm.Harbor.Server {
//				v.RemoteShell(fmt.Sprintf("mkdir -p %s", hm.Harbor.DataDir))
//			}
//		}
//	} else {
//
//		log.Println("检测docker...")
//		if !runner.LocalCommandExists("docker") {
//			log.Fatalf("docker检测失败...")
//		}
//
//		log.Println("检测docker-compose...")
//		if !runner.LocalCommandExists("docker-compose") {
//			log.Fatalf("docker-compose检测失败...")
//		}
//
//		log.Print("初始化harbor目录...")
//		os.Mkdir(hm.Harbor.DataDir, 0644)
//
//	}
//}
//
//func (hm *harborMgr) cleanRetainData() {
//	dirItems := []string{"ca_download", "database", "job_logs", "redis", "registry", "secret"}
//	if len(hm.Harbor.Server) > 0 {
//		for _, s := range hm.Harbor.Server {
//			for _, v := range dirItems {
//				_ = s.DelDirectory(fmt.Sprintf("%s/%s", hm.Harbor.DataDir, v))
//			}
//		}
//	} else {
//		for _, v := range dirItems {
//			_ = os.Remove(fmt.Sprintf("%s/%s", hm.Harbor.DataDir, v))
//		}
//	}
//}
//
//// 调整配置
//func (hm *harborMgr) config() {
//	content, _ := tmplutil.Render(HarborConfigTmpl, util.Data{
//		"data_dir":       hm.Harbor.DataDir,
//		"admin_password": hm.Harbor.Password,
//		"http_port":      hm.Harbor.HttpPort,
//		"domain":         hm.Harbor.Domain,
//	})
//	if len(hm.Harbor.Server) > 0 {
//		log.Println("配置harbor...")
//		for _, v := range hm.Harbor.Server {
//			log.Printf("生成配置文件 -> %s: %s ...", v.Host, harborYml)
//			v.WriteRemoteFile([]byte(content), harborYml, 0755)
//		}
//	} else {
//		log.Printf("生成配置文件 -> %s ...", harborYml)
//		err := ioutil.WriteFile(harborYml, []byte(content), 0644)
//		if err != nil {
//			panic(err)
//		}
//	}
//}
//
//func (hm *harborMgr) depress() {
//	if len(hm.Harbor.Server) > 0 {
//		cmd := "tar zxvf /tmp/harbor-offline-installer-v2.1.4.tgz -C /usr/local"
//		log.Println("解压安装包...")
//		remoteShell(hm.Harbor, cmd)
//	} else {
//		log.Println("解压安装包...")
//		re := runner.Shell(fmt.Sprintf("tar zxvf %s -C /usr/local", hm.Installer.OfflineFilePath))
//		if re.ExitCode != 0 {
//			log.Fatalf("解压安装包失败: %s", re.StdErr)
//		}
//	}
//}
//
//func (hm *harborMgr) loadHarborImages() {
//	// todo: 正则
//	img := "/usr/local/harbor/harbor.v2.1.4.tar.gz"
//
//	if len(hm.Harbor.Server) > 0 {
//		var wg sync.WaitGroup
//		wg.Add(len(hm.Harbor.Server))
//
//		for _, v := range hm.Harbor.Server {
//			log.Printf("%s -> 导入harbor镜像 ...", v.Host)
//			go func(path string) {
//				v.RemoteShell(fmt.Sprintf("docker load -i %s", path))
//				defer wg.Done()
//			}(img)
//		}
//		wg.Wait()
//	} else {
//		log.Printf("导入harbor镜像 ...")
//		runner.Shell(fmt.Sprintf("docker load -i %s", img))
//	}
//}
//
//// harbor预安装前准备
//func (hm *harborMgr) prepare() {
//	log.Println("预安装harbor...")
//	if len(hm.Harbor.Server) > 0 {
//		for _, v := range hm.Harbor.Server {
//			content, _ := tmplutil.Render(HarborPrepareTmpl, util.Data{
//				"domain":  hm.Harbor.Domain,
//				"address": hm.anotherAddress(v.Host),
//			})
//			v.RemoteShell(content)
//		}
//	} else {
//		content, _ := tmplutil.Render(HarborPrepareTmpl, util.Data{
//			"domain":  hm.Harbor.Domain,
//			"address": hm.Harbor.ResolvAddress,
//		})
//		runner.Shell(content)
//	}
//}
//
//func (hm *harborMgr) anotherAddress(ip string) string {
//	if len(hm.Harbor.Server) != 2 {
//		return ip
//	}
//	if hm.Harbor.Server[0].Host == ip {
//		return hm.Harbor.Server[1].Host
//	} else {
//		return hm.Harbor.Server[0].Host
//	}
//}
//
//// harbor安装
//func (hm *harborMgr) install() {
//	content, _ := tmplutil.Render(HarborInstallTmpl, util.Data{
//		"http_port": hm.Harbor.HttpPort,
//	})
//	log.Println("安装harbor...")
//	if len(hm.Harbor.Server) > 0 {
//		remoteShell(hm.Harbor, content)
//	} else {
//		runner.Shell(content)
//	}
//}
//
//func remoteShell(h runner.Harbor, cmd string) {
//	var wg sync.WaitGroup
//	for _, v := range h.Server {
//		wg.Add(1)
//		go func(server runner.Server, shell string) {
//			defer wg.Done()
//			server.RemoteShell(shell)
//		}(v, cmd)
//	}
//	wg.Wait()
//}
//
//// harbor状态
//func reachable(h runner.Harbor, host string) bool {
//
//	pingData := pingObject{
//		Ak: "admin", As: h.Password, Insecure: false,
//		Name: "replication", Type: "harbor",
//		Url: fmt.Sprintf("http://%s", h.Domain),
//	}
//	// 序列化
//	pingBody, _ := json.Marshal(&pingData)
//
//	url := fmt.Sprintf("http://%s:%s/api/v2.0/registries/ping", h.Domain, h.HttpPort)
//	log.Printf("Registry -> 检测%s可达性...", host)
//
//	retry := 5
//	for {
//		rep, _ := post(url, bytes.NewBuffer(pingBody), "admin", h.Password)
//		if successful(rep.StatusCode) {
//			break
//		}
//
//		log.Printf("Registry -> %s不可达...", host)
//		b, _ := ioutil.ReadAll(rep.Body)
//		log.Println(string(b))
//		time.Sleep(15 * time.Second)
//		retry--
//	}
//	log.Println("Registry -> 可达...")
//	return true
//}
//
//// 多节点添加
//func (hm harborMgr) addRegistries() *harborMgr {
//
//	// 检测harbor状态
//	for _, v := range hm.Harbor.Server {
//		reachable(hm.Harbor, v.Host)
//	}
//
//	// 添加registry
//	for _, v := range hm.Harbor.Server {
//		hm.addRegistry(v)
//	}
//
//	return &hm
//}
//
//// 单节点添加
//func (hm *harborMgr) addRegistry(server runner.Server) bool {
//
//	h := hm.Harbor
//	registryData.Credential.As = h.Password
//	registryData.Url = fmt.Sprintf("http://%s", h.Domain)
//	registryBody, _ := json.Marshal(&registryData)
//	url := fmt.Sprintf("http://%s:%s/api/v2.0/registries", server.Host, h.HttpPort)
//
//	// 添加registry
//	log.Printf("添加Registry -> %s...", server.Host)
//	rep, _ := post(url, bytes.NewBuffer(registryBody), "admin", h.Password)
//	if !successful(rep.StatusCode) {
//		return false
//	}
//	log.Printf("Registry -> %s添加成功...", server.Host)
//	return true
//}
//
//// 配置复制规则
//func (hm harborMgr) setReplications() *harborMgr {
//	//var url string
//	h := hm.Harbor
//	repl.DestRegistry = registryData
//	repl.DestRegistry.Credential = cred
//	repl.DestRegistry.Credential.As = h.Password
//	trigger := trigger{
//		Type:            "event_based",
//		TriggerSettings: triggerSettings{Cron: ""},
//	}
//	repl.Trigger = trigger
//	repl.Filters = []string{}
//
//	registryData.Url = fmt.Sprintf("http://%s:%s", h.Domain, h.HttpPort)
//
//	for _, v := range h.Server {
//		// 设置replica规则
//		host := fmt.Sprintf("http://%s:%s", v.Host, h.HttpPort)
//		setReplication(repl, host, h.Password)
//	}
//
//	return &hm
//}
//
//func setReplication(repl replica, host string, password string) {
//
//	log.Printf("Registry ID -> 获取 %s registry id...", host)
//	searchRegistryUrl := host + "/api/v2.0/registries?q%3Dname=replication"
//	repl.DestRegistry.ID = registryID(searchRegistryUrl, "admin", password)
//
//	// 序列化
//	replicaBody, _ := json.Marshal(&repl)
//	log.Printf("Replicas -> %s添加复制规则...", host)
//	url := fmt.Sprintf("%s/api/v2.0/replication/policies", host)
//	rep, _ := post(url, bytes.NewBuffer(replicaBody), "admin", password)
//
//	if successful(rep.StatusCode) {
//		log.Printf("Replicas -> %s添加复制规则成功...", host)
//	} else {
//		log.Println(rep.StatusCode)
//		b, _ := ioutil.ReadAll(rep.Body)
//		log.Println(string(b))
//		log.Fatalf("Replicas -> %s添加复制规则失败...", host)
//	}
//
//}
//
//// 根据registry name获取registry id
//func registryID(url string, username string, password string) int {
//	var registry []registry
//	b := util.Get(url, username, password)
//	err := json.Unmarshal(b, &registry)
//	if err != nil {
//		log.Fatalf("反序列化失败")
//	}
//	if len(registry) != 0 {
//		return registry[0].ID
//	}
//	return 0
//}
//
//// 初始化harbor
//func (hm *harborMgr) defaultProjects() harborMgr {
//	health := fmt.Sprintf("http://%s/api/v2.0/health", hm.Harbor.Domain)
//	private := hm.Harbor.Project.Private
//	public := hm.Harbor.Project.Public
//
//	// 组装
//	var projects []project
//	for _, v := range private {
//		var p project
//		p.Public = false
//		p.Name = v
//		projects = append(projects, p)
//	}
//
//	for _, v := range public {
//		var p project
//		p.Public = true
//		p.Name = v
//		projects = append(projects, p)
//	}
//
//	if !isHealthy(health, 5, time.Second*20) {
//		log.Fatal("harbor访问异常...")
//	}
//
//	hm.HarborProjects = projects
//	return *hm
//}
//
//func (hm harborMgr) reloadDockerDaemon() *harborMgr {
//	content, _ := tmplutil.Render(DockerConfigTempl, util.Data{
//		"InsecureRegistries": hm.Harbor.Domain,
//	})
//
//	if len(hm.Harbor.Server) != 0 {
//		hm.reloadRemoteDockerDaemon(content)
//	} else {
//		log.Printf("生成文件：%s\n", dockerDaemonConfigPath)
//		err := ioutil.WriteFile(dockerDaemonConfigPath, []byte(content), 0644)
//		if err != nil {
//			log.Fatalf("写文件：%s失败，失败原因为：%s\n", dockerDaemonConfigPath, err.Error())
//		}
//		log.Println("重载docker配置...")
//		runner.Shell(reloadDockerDaemonCmd)
//	}
//
//	return &hm
//}
//
//func (hm harborMgr) reloadRemoteDockerDaemon(content string) *harborMgr {
//
//	var wg sync.WaitGroup
//
//	for _, v := range hm.Harbor.Server {
//		wg.Add(1)
//		v.WriteRemoteFile([]byte(content), dockerDaemonConfigPath, 0755)
//
//		go func(server runner.Server) {
//			defer wg.Done()
//			server.RemoteShell(reloadDockerDaemonCmd)
//		}(v)
//	}
//	wg.Wait()
//
//	return &hm
//}
//
//func isHealthy(url string, retry int, interval time.Duration) bool {
//	for {
//		if retry == 0 {
//			break
//		}
//		log.Println("检测harbor是否可访问...")
//		req, err := http.NewRequest(http.MethodGet, url, nil)
//		if err != nil {
//			panic(err)
//		}
//		log.Println("向harbor发起探测请求...")
//		reps, err := http.DefaultClient.Do(req)
//		if err != nil {
//			panic(err)
//		}
//		log.Printf("解析状态码:%d...", reps.StatusCode)
//		if reps.StatusCode == http.StatusOK {
//			log.Println("harbor可提供api服务...")
//			return true
//		}
//		log.Printf("harbor暂时不可访问，%f秒后将再再次检测...", interval.Seconds())
//		time.Sleep(interval)
//	}
//	return false
//}
//
//// 初始化项目
//func (hm harborMgr) newProject() *harborMgr {
//
//	log.Println("初始化project...")
//	if len(hm.Harbor.Server) > 0 {
//		for _, s := range hm.Harbor.Server {
//			hm.createProject(s.Host)
//		}
//	} else {
//		hm.createProject(hm.Harbor.ResolvAddress)
//	}
//
//	return &hm
//}
//
//func (hm harborMgr) createProject(host string) {
//	url := fmt.Sprintf("http://%s:%s/api/v2.0/projects", host, hm.Harbor.HttpPort)
//	for _, v := range hm.HarborProjects {
//		b, err := json.Marshal(v)
//		if err != nil {
//			panic(err)
//		}
//		resp, err := post(url, bytes.NewBuffer(b), "admin", hm.Harbor.Password)
//		if err != nil {
//			panic(err)
//		}
//		log.Printf("创建项目 -> %s", v.Name)
//		if successful(resp.StatusCode) {
//			log.Printf("项目%s -> 创建成功...", v.Name)
//		} else {
//			log.Fatal(resp.StatusCode)
//		}
//	}
//}
//
//// GCSettingObject 配置Gc策略对象
//type GCSettingObject struct {
//	Schedule struct {
//		Type string `json:"type"`
//		Cron string `json:"cron"`
//	} `json:"schedule"`
//
//	Parameters struct {
//		DeleteUntagged bool `json:"delete_untagged"`
//	} `json:"parameters"`
//}
//
//func (hm harborMgr) gc() *harborMgr {
//
//	if !hm.Harbor.OpenGC {
//		return &hm
//	}
//	log.Println("配置gc策略...")
//	url := fmt.Sprintf("http://%s:%s/api/v2.0/system/gc/schedule", hm.Harbor.ResolvAddress, hm.Harbor.HttpPort)
//	gc := &GCSettingObject{}
//	gc.Parameters.DeleteUntagged = true
//	gc.Schedule.Type = "Hourly"
//	gc.Schedule.Cron = "0 0 * * * *"
//	b, err := json.Marshal(&gc)
//	if err != nil {
//		panic(err)
//	}
//	re, err := post(url, bytes.NewBuffer(b), "admin", hm.Harbor.Password)
//	if err != nil {
//		panic(err)
//	}
//
//	if re.StatusCode >= http.StatusCreated {
//		log.Println(re.StatusCode)
//		log.Println("配置gc失败...")
//	}
//
//	return &hm
//}
//
//func successful(code int) bool {
//	switch code {
//	case http.StatusOK:
//		return true
//	case http.StatusCreated:
//		return true
//	case http.StatusConflict:
//		return true
//	default:
//		return false
//	}
//}

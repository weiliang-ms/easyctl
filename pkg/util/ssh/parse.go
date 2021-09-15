package ssh

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"os"
	"reflect"
)

type ServerList struct {
	Common    CommonServerList
	Harbor    HarborServerList
	HA        HaProxyServerList
	Keepalive KeepaliveServerList
	Docker    DockerServerList
	Compose   DockerComposeServerList
}

type ExecResult struct {
	ExitCode int
	StdErr   string
	StdOut   string
}

type CommonServerList struct {
	Server []Server `yaml:"server,flow"`
}

type KeepaliveServerList struct {
	Attribute Keepalive `yaml:"keepalive"`
}

type Keepalive struct {
	Vip       string   `yaml:"vip"`
	Interface string   `yaml:"interface"`
	Server    []Server `yaml:"server,flow"`
}

type HaProxyServerList struct {
	Attribute HaProxy `yaml:"haproxy"`
}

type DockerServerList struct {
	Attribute Docker `yaml:"docker"`
}

type Docker struct {
	Servers []Server `yaml:"server,flow"`
}

type DockerComposeServerList struct {
	Attribute DockerCompose `yaml:"docker-compose"`
}

type DockerCompose struct {
	Server []Server `yaml:"server,flow"`
}

type HaProxy struct {
	Server      []Server  `yaml:"server,flow"`
	BalanceList []Balance `yaml:"balance,flow"`
}

type Balance struct {
	Name    string   `yaml:"name"`
	Port    int      `yaml:"port"`
	Address []string `yaml:"address"`
}

type HarborServerList struct {
	Attribute Harbor `yaml:"harbor"`
}

type Harbor struct {
	Server   []Server      `yaml:"server,flow"`
	Project  HarborProject `yaml:"project"`
	Password string        `yaml:"password"`
	DataDir  string        `yaml:"dataDir"`
	Domain   string        `yaml:"domain"`
	HttpPort string        `yaml:"http-port"`
	Vip      string        `yaml:"vip"`
}

type HarborProject struct {
	Private []string `yaml:"private"`
	Public  []string `yaml:"public"`
}

func ParseServerList(yamlPath string, v interface{}) (err error, list ServerList) {

	var decodeErr, marshalErr error

	f, err := os.Open(yamlPath)
	if err != nil {
		return err, list
	}

	// todo:优化反射方式
	switch reflect.ValueOf(v).Type().String() {
	case "runner.DockerServerList":
		decodeErr = yaml.NewDecoder(f).Decode(&list.Docker)
		_, marshalErr = json.Marshal(&list.Docker)
	case "runner.HarborServerList":
		decodeErr = yaml.NewDecoder(f).Decode(&list.Harbor)
		_, marshalErr = json.Marshal(&list.Docker)
	case "runner.KeepaliveServerList":
		decodeErr = yaml.NewDecoder(f).Decode(&list.Keepalive)
		_, marshalErr = json.Marshal(&list.Docker)
	default:
		decodeErr = yaml.NewDecoder(f).Decode(&list.Common)
		_, marshalErr = json.Marshal(&list.Docker)
	}

	if decodeErr != nil {
		return err, list
	}

	if marshalErr != nil {
		return err, list
	}

	return nil, list
}

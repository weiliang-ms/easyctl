package sys

import (
	"easyctl/constant"
	"easyctl/util"
	"errors"
)

type Service struct {
	Name    string
	Content string
}

const (
	redis  = "redis"
	docker = "docker"
)

var unSupportSystemErr = errors.New("暂不支持当前系统...")
var unrecognizedServiceErr = errors.New("暂不识别的服务名称")

func ConfigService(serviceName string) {
	switch serviceName {
	case redis:
		configRedisService(constant.Redhat7RedisServiceFilePath, constant.Redhat7RedisServiceContent)
	case docker:
		configDockerService(constant.Redhat7DockerServiceFilePath, constant.Redhat7DockerServiceContent)
	default:
		panic(unrecognizedServiceErr)
	}
}

func configRedisService(path string, content string) {
	systemType := SystemInfoObject.OSVersion.ReleaseType
	mainNumber := SystemInfoObject.OSVersion.MainVersionNumber
	//fmt.Println("----" + systemType + mainNumber)
	if systemType != RedhatReleaseType {
		panic(unSupportSystemErr)
	}
	if mainNumber == "7" {
		util.CreateFile(path, content)
		util.ExecuteCmd("systemctl daemon-reload")
	} else {
		panic(unSupportSystemErr)
	}
}

func configDockerService(path string, content string) {
	systemType := SystemInfoObject.OSVersion.ReleaseType
	mainNumber := SystemInfoObject.OSVersion.MainVersionNumber
	//fmt.Println("----" + systemType + mainNumber)
	if systemType != RedhatReleaseType {
		panic(unSupportSystemErr)
	}
	if mainNumber == "7" {
		util.CreateFile(path, content)
		util.ExecuteCmd("systemctl daemon-reload")
	} else {
		panic(unSupportSystemErr)
	}
}

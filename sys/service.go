package sys

import (
	"easyctl/resources"
	"easyctl/util"
	"errors"
)

type Service struct {
	Name    string
	Content string
}

const (
	redis                       = "redis"
	Redhat7RedisServiceFilePath = "/usr/lib/systemd/system/redis.service"
)

var unSupportSystemErr = errors.New("暂不支持当前系统...")
var unrecognizedServiceErr = errors.New("暂不识别的服务名称")

func ConfigService(serviceName string) {
	switch serviceName {
	case redis:
		configRedisService()
	default:
		panic(unrecognizedServiceErr)
	}
}

func configRedisService() {
	systemType := SystemInfoObject.OSVersion.ReleaseType
	mainNumber := SystemInfoObject.OSVersion.MainVersionNumber
	//fmt.Println("----" + systemType + mainNumber)
	if systemType != RedhatReleaseType {
		panic(unSupportSystemErr)
	}
	if mainNumber == "7" {
		util.CreateFile(Redhat7RedisServiceFilePath, resources.Redhat7RedisServiceContent)
		util.ExecuteCmd("systemctl daemon-reload")
	} else {
		panic(unSupportSystemErr)
	}
}

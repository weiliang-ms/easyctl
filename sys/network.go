package sys

import (
	"fmt"
	"github.com/weiliang-ms/easyctl/util"
)

const defaultDNSAddress = "114.114.114.114"

// 检测主机是否连接到Internet
func AccessAliMirrors() bool {
	// 检测是否可以访问阿里镜像源
	resolvRe, _ := util.ExecuteCmdAcceptResult("ping mirrors.aliyun.com -c 1 -W 2;echo $?")
	accessDNSRe, _ := util.ExecuteCmdAcceptResult("ping 114.114.114.114 -c 1 -W 2;echo $?")
	if resolvRe == "0" {
		return true
	} else {
		if accessDNSRe == "0" {
			SetDefaultDNS()
			return true
		} else {
			fmt.Println("无法连接互联网...")
			return false
		}
	}
}

func SetDefaultDNS() {
	fmt.Println("配置默认dns: 114.114.114.114")
	SetDNS(defaultDNSAddress)
}

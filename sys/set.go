package sys

import (
	"easycfg/util"
	"fmt"
)

func SetDNS(dnsAddress string) (err error, result string) {

	cmd := "sed -i \"/nameserver " + dnsAddress + "/d\" /etc/resolv.conf;" +
		"echo \"nameserver " + dnsAddress + "\" >> /etc/resolv.conf\n"

	fmt.Printf("[check] 检测dns地址：%s合法性...\n", dnsAddress)
	err, result = util.CheckIP(dnsAddress)
	if err != nil {
		return err, result
	}
	shellErr, shellResult := util.ExecuteCmd(cmd)
	return shellErr, shellResult
}

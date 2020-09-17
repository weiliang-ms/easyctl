package sys

import "easycfg/util"

// 查询端口监听
func SearchPortStatus(port string) error {
	cmd := "ss -alnpt|grep " + port
	return util.ExecuteCmdAcceptResult(cmd)
}

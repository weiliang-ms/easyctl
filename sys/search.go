package sys

import "github.com/weiliang-ms/easyctl/util"

// 查询端口监听
func SearchPortStatus(port string) (result string, err error) {
	cmd := "ss -alnpt|grep " + port
	return util.ExecuteCmdAcceptResult(cmd), nil
}

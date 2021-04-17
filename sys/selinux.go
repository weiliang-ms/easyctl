package sys

import (
	"github.com/weiliang-ms/easyctl/constant"
	"github.com/weiliang-ms/easyctl/util"
)

func SelinuxStatus() string {
	return util.ExecuteCmdAcceptResult(constant.SelinuxStatusCmd)
}

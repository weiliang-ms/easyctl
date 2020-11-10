package sys

import (
	"easyctl/constant"
	"easyctl/util"
)

func SelinuxStatus() string {
	return util.ExecuteCmdAcceptResult(constant.SelinuxStatusCmd)
}

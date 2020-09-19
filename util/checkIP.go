package util

import (
	"easyctl/cutomErr"
	"regexp"
)

// ip合法性检测
func CheckIP(address string) (err error, result string) {
	matched, err := regexp.MatchString("^(\\d{1}|[1-9]{1}\\d{1}|1\\d\\d|2[0-4]\\d|25[0-5])\\.(\\d{1}|[1-9]{1}\\d{1}|1\\d\\d|2[0-4]\\d|25[0-5])\\.(\\d{1}|[1-9]{1}\\d{1}|1\\d\\d|2[0-4]\\d|25[0-5])\\.(\\d{1}|[1-9]{1}\\d{1}|1\\d\\d|2[0-4]\\d|25[0-5])$", address)
	if matched == false {
		err = cutomErr.InvalidateIPAddress
		result = ""
	} else if matched == true {
		err = nil
		result = ""
	}
	return err, result
}

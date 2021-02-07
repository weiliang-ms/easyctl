package sys

import (
	"easyctl/constant"
	"easyctl/util"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func KernelVersion() (version string) {
	file, err := os.OpenFile("/proc/version", os.O_RDONLY, 0664)
	if err != nil {
		log.Fatal(err.Error())
	}
	b, readErr := ioutil.ReadAll(file)
	if readErr != nil {
		log.Fatal(readErr.Error())
	}
	for _, v := range strings.Split(string(b), " ") {
		matched, _ := regexp.MatchString("^[0-9]\\.[0-9]*\\.[0-9]*-[0-9]*.*.el[0-9].*x86_64$", v)
		if matched {
			version = v
		}
	}
	util.PrintDirectBanner([]string{constant.Kernel}, fmt.Sprintf("内核版本: %s", version))
	return version
}

func KernelMainVersion() int {
	var mVersion string

	version := KernelVersion()
	if version != "" {
		mVersion = strings.Split(version, ".")[0]
	} else {
		log.Fatal("获取内核版本失败...")
	}
	int, err := strconv.Atoi(mVersion)
	if err != nil {
		log.Fatal("内核主版本获取失败...")
	}
	return int
}

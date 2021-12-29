package validity

import (
	"fmt"
	"strings"
)

// DataPath 目录合法性检测
/*
	- 必须为绝对路径
	- 不应为/proc、/sys、/boot等
*/

var (
	invalidDataDir = map[string]struct{}{
		"/bin":   {},
		"/boot":  {},
		"/lib":   {},
		"/lib64": {},
		"/proc":  {},
		"/run":   {},
		"/sbin":  {},
		"/sys":   {},
	}
)

type relativePathErr struct {
	Msg string
}

type sensitivePathErr struct {
	Msg string
}

func (e relativePathErr) Error() string {
	return fmt.Sprintf("%s 路径非法，不是绝对路径", e.Msg)
}

func (e sensitivePathErr) Error() string {
	return fmt.Sprintf("路径非法，不允许在%s目录下", e.Msg)
}

func DataPath(path string) error {
	// 1.判断是否绝对路径
	if !strings.HasPrefix(path, "/") {
		return relativePathErr{Msg: path}
	}

	// 2.判断是否为敏感目录
	dirSlice := strings.Split(path, "/")
	ancestorDir := fmt.Sprintf("/%s", dirSlice[1])
	if _, ok := invalidDataDir[ancestorDir]; ok {
		return sensitivePathErr{Msg: ancestorDir}
	}

	return nil
}
